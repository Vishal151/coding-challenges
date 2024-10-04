package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	backends []*url.URL
	mutex    sync.Mutex
	current  int
}

func NewLoadBalancer(backends []string) *LoadBalancer {
	var urls []*url.URL
	for _, backend := range backends {
		url, _ := url.Parse(backend)
		urls = append(urls, url)
	}
	return &LoadBalancer{backends: urls}
}

func (lb *LoadBalancer) NextBackend() *url.URL {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	backend := lb.backends[lb.current]
	lb.current = (lb.current + 1) % len(lb.backends)
	return backend
}

func main() {
	lb := NewLoadBalancer([]string{
		"http://localhost:8081",
		"http://localhost:8082",
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, false)
		fmt.Printf("Received request from %s\n%s\n", r.RemoteAddr, string(dump))

		backend := lb.NextBackend()
		proxy := httputil.NewSingleHostReverseProxy(backend)

		proxy.ServeHTTP(w, r)

		fmt.Printf("Request forwarded to backend server: %s\n\n", backend)
	})

	log.Println("Starting load balancer on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
