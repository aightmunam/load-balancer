package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at /")
	fmt.Fprint(w, "Backend: Success")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at about/")
	fmt.Fprint(w, "About Me: I am a very bare bones backend server")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/", indexHandler)

	server := &http.Server{
		Addr:    "127.0.0.1:3333",
		Handler: mux,
	}
	fmt.Println("Backend is listening requests at 127.0.0.1:3333")
	server.ListenAndServe()
}
