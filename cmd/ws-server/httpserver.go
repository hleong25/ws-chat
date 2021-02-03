package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	done chan bool
}

func New(done chan bool) HTTPServer {
	return HTTPServer{
		done,
	}
}

func (httpServer HTTPServer) Start() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hello world")

		msg := fmt.Sprintf("hello world %s\n", time.Now().Local())

		w.Write([]byte(msg))
	})

	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("quit")
		close(httpServer.done)
	})

	log.Printf("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
