package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sankalpjonn/mockingbird/bird"
)

func main() {
	r := mux.NewRouter()
	bird := bird.New()
	host := flag.String("host", "0.0.0.0:8000", "ip:port")
	flag.Parse()

	r.HandleFunc("/egg", bird.CreateHandler).Methods("POST")
	r.HandleFunc("/egg/{eggId}", bird.GetHandler).Methods("POST", "GET", "PUT", "HEAD", "PATCH")
	srv := &http.Server{
		Handler: r,
		Addr:    *host,
	}
	/***** Do not use Fatal anywhere else ********/
	// I don't know who you are. I don't know what you want. If you are looking
	// to use "log.Fatal", I can tell you I don't know you. But what I do have
	// are a very particular set of skills, git skills that I have acquired over
	// a very long career. Skills that make me a find out people from the git
	// commit. If you let this server go now, that'll be the end of it. I will
	// not look for you, I will not pursue you. But if you don't, I will look
	// for you, I will find you, and I will revoke commit access for you.
	log.Fatal(srv.ListenAndServe())
}
