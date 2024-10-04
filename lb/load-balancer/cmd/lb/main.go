package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Backend struct {
	URL     *url.URL
	Healthy bool
	mux     sync.RWMutex
}

type LoadBalancer struct {
	backends []*Backend
	mutex    sync.Mutex
	current  int
}

func NewLoadBalancer(backends []string) *LoadBalancer {
	var be []*Backend
	for _, backend := range backends {
		url, _ := url.Parse(backend)
		be = append(be, &Backend{URL: url, Healthy: true})
	}
	return &LoadBalancer{backends: be, current: -1}
}

func (b *Backend) SetHealth(health bool) {
	b.mux.Lock()
	b.Healthy = health
	b.mux.Unlock()
}

func (b *Backend) IsHealthy() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Healthy
}

func (lb *LoadBalancer) NextBackend() *Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i := 0; i < len(lb.backends); i++ {
		lb.current = (lb.current + 1) % len(lb.backends)
		if lb.backends[lb.current].IsHealthy() {
			return lb.backends[lb.current]
		}
	}
	return nil
}

func (lb *LoadBalancer) HealthCheck() {
	for _, backend := range lb.backends {
		go func(b *Backend) {
			for {
				resp, err := http.Get(b.URL.String() + "/health")
				if err != nil || resp.StatusCode != http.StatusOK {
					b.SetHealth(false)
					log.Printf("Backend %v is unhealthy\n", b.URL)
				} else {
					b.SetHealth(true)
					log.Printf("Backend %v is healthy\n", b.URL)
				}
				time.Sleep(5 * time.Second)
			}
		}(backend)
	}
}

func main() {
	backendURLs := strings.Split(os.Getenv("BACKEND_URLS"), ",")
	if len(backendURLs) == 0 {
		backendURLs = []string{
			"http://localhost:8081",
			"http://localhost:8082",
		}
	}

	lb := NewLoadBalancer(backendURLs)

	go lb.HealthCheck()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, false)
		fmt.Printf("Received request from %s\n%s\n", r.RemoteAddr, string(dump))

		backend := lb.NextBackend()
		if backend == nil {
			http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(backend.URL)
		proxy.ServeHTTP(w, r)

		fmt.Printf("Request forwarded to backend server: %s\n\n", backend.URL)
	})

	log.Println("Starting load balancer on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
