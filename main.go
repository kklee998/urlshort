package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("Path requested: %s\n", path)
		fmt.Fprintf(w, "Hello, you've requested: %s\n", path)
	})
	log.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", nil)
}
