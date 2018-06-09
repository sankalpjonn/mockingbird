package main

import (
  "net/http"
  "log"
  "os"

  "github.com/mockingbird/bird"
  "github.com/gorilla/mux"
)

func main() {
  r := mux.NewRouter()
  bird := bird.New()

	r.HandleFunc("/egg", bird.CreateHandler).Methods("POST")
	r.HandleFunc("/egg/{eggId}", bird.GetHandler).Methods("POST", "GET", "PUT", "HEAD", "PATCH")

	srv := &http.Server{
		Handler: r,
		Addr:    os.Args[1],
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
