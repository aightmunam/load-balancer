package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

var backendServerAddrs = []string{"http://127.0.0.1:3333", "http://127.0.0.1:3334", "http://127.0.0.1:3335"}
var totalServers uint32 = uint32(len(backendServerAddrs))
var currentServer uint32 = 0

// Simple server choice based on Round Robin without any weights
func getEligibleServer() string {
	atomic.AddUint32(&currentServer, 1)
	return backendServerAddrs[(currentServer)%totalServers]
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at load balancer: 127.0.0.1:8080")
	response, err := http.Get(getEligibleServer() + r.URL.Path + "/")
	if err != nil {
		fmt.Fprint(w, "Server is down.")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprint(w, "Something went wrong.")
	}
	defer response.Body.Close()

	fmt.Println("Backend server responded with: ", string(body))
	fmt.Fprint(w, string(body))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", generalHandler)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	fmt.Println("Load balancer is listening requests at 127.0.0.1:8080")
	server.ListenAndServe()
}
