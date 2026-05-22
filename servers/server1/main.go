package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello from Server 1")
    })

    fmt.Println("Server 1 running on port 8081")
    http.ListenAndServe(":8081", nil)
}