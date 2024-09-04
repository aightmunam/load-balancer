package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Served by: ", r.Host)
	fmt.Println("Received request at /")
	fmt.Fprint(w, "Hello from: "+r.Host)
}

func main() {
	ch := make(chan interface{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

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

	go func() {
		server_one.ListenAndServe()
		ch <- true
	}()
	go func() {
		server_two.ListenAndServe()
		ch <- true
	}()
	go func() {
		server_three.ListenAndServe()
		ch <- true
	}()

	<-ch
}
