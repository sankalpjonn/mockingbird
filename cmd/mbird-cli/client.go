package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Egg struct {
	Id string `json:"egg_id"`
}

type client struct {
	headers      Headers
	bodyFilePath string
	body         string
	server       string
	statusCode   int
	ttl          int
	eggid        string
	isCreate     bool
	isGet        bool
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

func (self *client) getPayload() (string, error) {
	body, err := self.getBody()
	if err != nil {
		return "", err
	}
	payload := map[string]interface{}{
		"body":        body,
		"ttl":         self.ttl,
		"status_code": self.statusCode,
	}
	if len(self.headers) > 0 {
		payload["headers"] = self.headers
	}
	if self.eggid != "" {
		payload["id"] = self.eggid
	}
	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return string(b), nil
}

func (self *client) getEgg() (string, error) {
	if self.eggid == "" {
		return "", errors.New("Please provide -egg")
	}
	url := fmt.Sprintf("http://%s/egg/%s/raw", self.server, self.eggid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}

func (self *client) createEgg() (string, error) {
	url := fmt.Sprintf("http://%s/egg", self.server)
	p, err := self.getPayload()
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(p)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	egg := &Egg{}
	err = json.Unmarshal(body, egg)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("------------\n\nhttp://%s/egg/%s\n\n------------", self.server, egg.Id), nil
}

func (self *client) callServer() string {
	var response string
	if self.isCreate {
		res, err := self.createEgg()
		if err != nil {
			response = fmt.Sprintf("%s", err)
		} else {
			response = res
		}
	} else if self.isGet {
		res, err := self.getEgg()
		if err != nil {
			response = fmt.Sprintf("%s", err)
		} else {
			response = res
		}
	} else {
		response = "Please provide either -create or -get in the flags"
	}
	return response
}
