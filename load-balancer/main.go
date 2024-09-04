package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"sync/atomic"
	"time"
)

var backendServerAddrs = []string{"http://127.0.0.1:3333", "http://127.0.0.1:3334", "http://127.0.0.1:3335", "http://127.0.0.1:8001"}

var availableServers []string = append([]string{}, backendServerAddrs...)
var currentServer uint32 = 0

// Naive Round Robin implementation: Find the next server for each incoming request
func getEligibleServer() string {
	var totalServers uint32 = uint32(len(availableServers))
	atomic.AddUint32(&currentServer, 1)
	return availableServers[(currentServer)%totalServers]
}

// Send request to all servers. If any of the servers is down, make it unavailable to the load balancer.
// If a server is up and running, make it available to the load balancer.
func healthCheck() {
	for index, server := range backendServerAddrs {
		_, err := http.Get(server + "/ping")
		isCurrentlyActive := slices.Contains(availableServers, server)
		if err != nil {
			if isCurrentlyActive {
				fmt.Println("Removing server: ", server)
				availableServers = slices.Delete(availableServers, index, index+1)
			}
		} else {
			if !isCurrentlyActive {
				fmt.Println("Adding server: ", server)
				availableServers = append(availableServers, server)
			}
		}
	}
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
	log.Output(1, "Received request at load balancer: "+r.Host)

	response, err := http.Get(getEligibleServer() + r.URL.Path + "/")
	if err != nil {
		fmt.Fprint(w, "Server is down.")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprint(w, "Something went wrong.")
	}
	defer response.Body.Close()

	log.Output(1, "Backend server responded with: "+string(body))
	fmt.Fprint(w, string(body))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", generalHandler)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	log.Output(1, "Load balancer is listening at: "+server.Addr)

	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				healthCheck()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	server.ListenAndServe()
}
