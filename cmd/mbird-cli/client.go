package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/bogdanovich/dns_resolver"
	"github.com/sankalpjonn/mockingbird/bird"
	"github.com/satori/go.uuid"
	"gopkg.in/abiosoft/ishell.v2"
	"gopkg.in/fatih/color.v1"
)

type Recording struct {
	Headers    map[string]string
	StatusCode int
	Body       string
}

type client struct {
	headers         Headers
	bodyFilePath    string
	body            string
	server          string
	statusCode      int
	ttl             int
	eggid           string
	isCreate        bool
	isGet           bool
	isRecord        bool
	isRecording     bool
	recDomain       string
	proxyServerPort string
}

func New() *client {
	c := new(client)
	flag.Var(&(c.headers), "H", `List of headers ("header=value")`)
	flag.StringVar(&(c.bodyFilePath), "bodyf", "", "Path of the file that contains body (/path/to/body/file)")
	flag.StringVar(&(c.body), "body", "", "Content of the body")
	flag.StringVar(&(c.server), "server", "localhost:8000", "address of the mockingbird server for this client")
	flag.IntVar(&(c.statusCode), "status", 200, "Status code")
	flag.IntVar(&(c.ttl), "ttl", 600, "Time to live for this mock")
	flag.BoolVar(&(c.isCreate), "create", false, "set if this command is for creating a mock")
	flag.BoolVar(&(c.isGet), "get", false, "set if this command is for retrieving a mock")
	flag.BoolVar(&(c.isRecord), "record", false, "set if this command is for recording a mock")
	flag.StringVar(&(c.proxyServerPort), "pport", "8080", "Port for the proxy server to listen on if record option is selected")
	flag.StringVar(&(c.eggid), "egg", "", "id of the mock being retrieved")
	flag.Parse()
	return c
}

func (self *client) getBody() (string, error) {
	body := ""
	if self.bodyFilePath != "" {
		dat, err := ioutil.ReadFile(self.bodyFilePath)
		if err != nil {
			return "", err
		}
		body = string(dat)
	} else if self.body != "" {
		body = self.body
	}

	return body, nil
}

func (self *client) getEgg(eggid string) (error, *bird.Egg) {
	if eggid == "" {
		return errors.New("Please provide -egg"), nil
	}
	url := fmt.Sprintf("http://%s/egg/%s/raw", self.server, self.eggid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, nil
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	egg := &bird.Egg{}
	err = json.Unmarshal(body, egg)
	if err != nil {
		panic(err)
	}
	return nil, egg
}

func (self *client) createEgg(egg *bird.Egg) error {
	url := fmt.Sprintf("http://%s/egg", self.server)
	b, _ := json.Marshal(egg)
	payload := strings.NewReader(string(b))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	bodyMap := map[string]string{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		panic(err)
	}
	egg.Id = bodyMap["egg_id"]
	return nil
}

func (self *client) recordReqRes(db *DB, r *http.Request, receiver chan *http.Response, done chan bool) {
	res := <-receiver
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	res.Body.Close()

	uuid4, _ := uuid.NewV4()
	data := Data{
		Headers: map[string][]string(r.Header),
		Body:    string(b),
		Status:  res.StatusCode,
		Id:      fmt.Sprintf("%s", uuid4),
	}
	collection := fmt.Sprintf("%s/%s/%s", self.recDomain, r.Method, r.URL.Path)

	db.record(collection, data)

	body := ioutil.NopCloser(bytes.NewReader(b))
	res.Body = body
	done <- true
}

func (self *client) setTransporter() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		remote := strings.Split(addr, ":")
		resolver := dns_resolver.New([]string{"114.114.114.114", "114.114.115.115", "119.29.29.29", "223.5.5.5", "8.8.8.8", "208.67.222.222", "208.67.220.220", "1.1.1.1"})
		resolver.RetryTimes = 5
		ip, err := resolver.LookupHost(remote[0])
		if err != nil {
			panic(err)
		}
		addr = ip[0].String() + ":" + remote[1]
		return dialer.DialContext(ctx, network, addr)
	}
}

func (self *client) serveProxyServerResponse(w http.ResponseWriter, r *http.Request, receiver chan *http.Response, done chan bool) {
	u, err := url.Parse(self.recDomain)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &transport{http.DefaultTransport, receiver, done}
	r.Host = u.Host
	proxy.ServeHTTP(w, r)
}

func (self *client) proxyFuncHandler(db *DB) func(http.ResponseWriter, *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !self.isRecording {
			//When recording stops

			collection := fmt.Sprintf("%s/%s/%s", self.recDomain, r.Method, r.URL.Path)
			err, data := db.replayRandom(collection)
			if err != nil {
				panic(err)
			}
			for k, v := range data.Headers {
				if len(v) > 0 {
					w.Header().Add(k, strings.Join(v, ";"))
				}
			}
			w.WriteHeader(data.Status)
			if data.Body != "" {
				w.Write([]byte(data.Body))
			}

		} else {

			//When recording starts
			self.setTransporter()
			responseReceiver := make(chan *http.Response)
			loggingDone := make(chan bool)
			go self.recordReqRes(db, r, responseReceiver, loggingDone)
			self.serveProxyServerResponse(w, r, responseReceiver, loggingDone)
		}
	}
	return fn
}

func (self *client) startProxyServer() {
	db, err := newDB("mockingbird_recordings")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", self.proxyFuncHandler(db))
	http.ListenAndServe(":"+self.proxyServerPort, nil)
}

func (self *client) newRecordingShell() error {
	shell := ishell.New()

	green := color.New(color.FgGreen).SprintFunc()
	shell.Println(green("************************************\nWELCOME TO THE MOCKING BIRD RECORDER\n\nPlease enter `help` for assistance\n************************************"))

	shell.AddCmd(&ishell.Cmd{
		Name: "domain",
		Help: "Set domain which will be proxied (http://<domain>.com or https://<domain>.com)",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				red := color.New(color.FgRed).SprintFunc()
				c.Println(red("please provide the domain in args !"))
				return
			}
			self.recDomain = c.Args[0]
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "Start recording",
		Func: func(c *ishell.Context) {
			if self.recDomain == "" {
				red := color.New(color.FgRed).SprintFunc()
				c.Println(red("Please provide a domain to record using the domain command !"))
				return
			}
			self.isRecording = true
			blue := color.New(color.FgBlue).SprintFunc()
			c.Println(blue("started recording ...\n\nPlease use localhost:" + self.proxyServerPort + " to access " + self.recDomain + ".\n\nuse the `stop` command to stop recording"))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "Stop recording",
		Func: func(c *ishell.Context) {
			self.isRecording = false
			yellow := color.New(color.FgYellow).SprintFunc()
			c.Println(yellow("stopped recording"))
		},
	})

	go self.startProxyServer()
	shell.Run()
	return nil
}

func (self *client) runCmd() (error, string) {

	if self.isCreate {

		//========CREATE COMMAND BEGIN ========//
		body, err := self.getBody()
		if err != nil {
			return err, ""
		}
		egg := &bird.Egg{
			Id:         self.eggid,
			Headers:    self.headers,
			Body:       body,
			StatusCode: self.statusCode,
			TTL:        self.ttl,
		}
		self.createEgg(egg)
		return nil, fmt.Sprintf("-------\n\nhttp://%s/egg/%s\n\n-------", self.server, egg.Id)
		//========CREATE COMMAND END ========//

	} else if self.isGet {

		//========GET COMMAND BEGIN ========//
		err, egg := self.getEgg(self.eggid)
		if err != nil {
			return err, ""
		}
		b, err := json.Marshal(egg)
		if err != nil {
			return err, ""
		}
		return nil, string(b)
		//========GET COMMAND END ========//

	} else if self.isRecord {

		//========RECORD COMMAND BEGIN ========//
		err := self.newRecordingShell()
		if err != nil {
			return err, ""
		}
		return nil, ""
		//========RECORD COMMAND END ========//

	} else {
		return errors.New("Please provide either -create or -get or -record in the flags"), ""
	}
}
