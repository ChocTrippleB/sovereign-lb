package main

import (
    "fmt"
    "net/http" // handles HTTP requests/responses
	"net/http/httputil" // reverse proxy, forwards traffic to backend servers. In load-balancers
	"net/url" // Used for working with URLs.
	"sync" // For concurrency safety (threads) 'go'
)

var servers = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
}

// Mutex to protect access to the current server index.
// "Only one person can enter this room at once."
var mu sync.Mutex
var current int

func getNextServer() string {
	mu.Lock() //only one goroutine can execute this critical section at a time.
	defer mu.Unlock() // Run this later before the function exits.
	server := servers[current] //   := variable declaration and assignment in one step. in C# += 
	current = (current + 1) % len (servers)  // Move to the next server index, get remainder to loop back to the start of the list.
	return server
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	target := getNextServer()
	fmt.Println("Forwarding to: ", target)

	url, _ := url.Parse(target) // no try-catch or exceptions in Go.
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}


func main() {
    
}