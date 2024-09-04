package main

import (
	"fmt"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Output(1, "Served by: "+r.Host)
	fmt.Fprint(w, "Hello from: "+r.Host)
}

func pingHander(w http.ResponseWriter, r *http.Request) {
	log.Output(1, "health check for "+r.Host)
	fmt.Fprint(w, "pong")
}

func runServer(server *http.Server) {
	log.Output(1, "Listening at: "+server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func main() {
	ch := make(chan interface{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/ping/", pingHander)

	server_one := &http.Server{
		Addr:    "127.0.0.1:3333",
		Handler: mux,
	}
	server_two := &http.Server{
		Addr:    "127.0.0.1:3334",
		Handler: mux,
	}
	server_three := &http.Server{
		Addr:    "127.0.0.1:3335",
		Handler: mux,
	}

	for _, server := range []*http.Server{server_one, server_two, server_three} {
		go func(s *http.Server) {
			runServer(s)
			ch <- true
		}(server)
	}

	<-ch
}
