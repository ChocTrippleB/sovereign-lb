package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello from Server 2")
    })

    fmt.Println("Server 1 running on port 8082")
    http.ListenAndServe(":8082", nil)
}