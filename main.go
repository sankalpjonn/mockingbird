package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sankalpjonn/mockingbird/bird"
)

func router() *mux.Router {
	bird := bird.New()
	r := mux.NewRouter()

	// apis
	r.HandleFunc("/egg", bird.CreateHandler).Methods("POST")
	r.HandleFunc("/egg/{eggId}/raw", bird.GetRawHandler).Methods("GET")
	r.HandleFunc("/egg/{eggId}", bird.GetHandler).Methods("POST", "GET", "PUT", "HEAD", "PATCH")

	return r
}

func main() {
	// parse command line args
	host := flag.String("host", "0.0.0.0:8000", "ip:port")
	flag.Parse()

	// get router for apis
	r := router()

	// start http server
	srv := &http.Server{
		Handler: r,
		Addr:    *host,
	}
	log.Fatal(srv.ListenAndServe())
}
