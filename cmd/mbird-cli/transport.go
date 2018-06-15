package main

import (
	"net/http"
)

type transport struct {
	http.RoundTripper
	responseReceiver chan *http.Response
	loggingDone      chan bool
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		panic(err)
	}
	t.responseReceiver <- resp
	<-t.loggingDone
	return resp, nil
}
