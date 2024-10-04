package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, false)
		fmt.Printf("Received request from %s\n%s\n", r.RemoteAddr, string(dump))

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello From Backend Server")
		fmt.Println("Replied with a hello message")
	})

	log.Println("Starting backend server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
