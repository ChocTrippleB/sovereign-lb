package main

import (
    "fmt"
    "net/http" // handles HTTP requests/responses
	"net/http/httputil" // reverse proxy, forwards traffic to backend servers. In load-balancers
	"net/url" // Used for working with URLs.
	"sync" // For concurrency safety (threads) 'go'
	"time"
)

type Server struct {
	URL string
	Alive bool
}

var servers = []*Server{
	{URL: "http://localhost:8081", Alive: true},
	{URL: "http://localhost:8082", Alive: true},
	{URL: "http://localhost:8083", Alive: true},
}

// Mutex to protect access to the current server index.
// "Only one person can enter this room at once."
var mu sync.Mutex
var current int

func getNextServer() *Server {
	mu.Lock() //only one goroutine can execute this critical section at a time.
	defer mu.Unlock() // Run this later before the function exits.
	
	total := len(servers)
	for i := 0; i < total; i++ {
		idx := (current + i) % total
		if servers[idx].Alive {
			current = (idx + 1 ) % total
			return servers[idx]
		}
	}
	return nil
}

func healthCheck() {
	for { 
		for _, server := range servers { 
			resp, err := http.Get(server.URL)
			mu.Lock()

			if err != nil || resp == nil || resp.StatusCode != 200 {
				if server.Alive {
					fmt.Println("Server Down: ", server.URL)
				}
				server.Alive = false
			} else {
				if !server.Alive {
					fmt.Println("Server Up: ", server.URL)
				}
				server.Alive = true
			}
			if err == nil && resp != nil {
				resp.Body.Close() // Close the response body to prevent resource leaks.
			}
			mu.Unlock()
		}
		time.Sleep(5 * time.Second)
	}
}


func handleRequest(w http.ResponseWriter, r *http.Request) {
	server := getNextServer()

	if server == nil {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	fmt.Println("Forwarding to: ", server.URL)

	url, _ := url.Parse(server.URL) // no try-catch or exceptions in Go.
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}


func main() {
	go healthCheck()

    http.HandleFunc("/", handleRequest)
	fmt.Println("Load blancer running on port 8080")
	http.ListenAndServe(":8080", nil)
}