package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sankalpjonn/mockingbird/bird"
	"gopkg.in/abiosoft/ishell.v2"
	"gopkg.in/fatih/color.v1"
)

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

func (self *client) generateEggFromRequest(r *http.Request) string {
	return "abc"
}

func (self *client) serveEgg(egg string, w http.ResponseWriter, r *http.Request) {

}

func (self *client) proxyFunc(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(self.recDomain)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if !self.isRecording {
		// egg := self.generateEggFromRequest(r)
		//
		// urlStr := fmt.Sprintf("http://%s/egg/%s", self.server, egg)
		// u, _ := url.Parse(urlStr)
		// fmt.Println("GOT URL", u)
		// // proxy := httputil.NewSingleHostReverseProxy(u)
		// proxy := reverseproxy.NewReverseProxy(u)
		// // req, _ := http.NewRequest("GET", urlStr, nil)
		// // fmt.Println("REQ", req, w)
		// proxy.ServeHTTP(w, req)
		// fmt.Println("SERVED", u)
	} else {
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(w, r)
	}
}

func (self *client) startProxyServer() {
	http.HandleFunc("/", self.proxyFunc)
	http.ListenAndServe(":"+self.proxyServerPort, nil)
}

func (self *client) newRecordingShell() error {
	shell := ishell.New()

	// display welcome info.
	shell.Println("WELCOME TO THE MOCKING BIRD RECORDER!")

	shell.AddCmd(&ishell.Cmd{
		Name: "domain",
		Help: "Set domain which will be proxied",
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
			c.Println(blue("started recording ..."))
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

func (self *client) callServer() (error, string) {
	if self.isCreate {
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
	} else if self.isGet {
		err, egg := self.getEgg(self.eggid)
		if err != nil {
			return err, ""
		}
		b, err := json.Marshal(egg)
		if err != nil {
			return err, ""
		}
		return nil, string(b)
	} else if self.isRecord {
		err := self.newRecordingShell()
		if err != nil {
			return err, ""
		}
		return nil, ""
	} else {
		return errors.New("Please provide either -create or -get or -record in the flags"), ""
	}
}
