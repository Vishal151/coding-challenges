package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vishal151/load-balancer/internal/balancer"
	"golang.org/x/time/rate"
)

type Backend struct {
	*balancer.Backend
	limiter             *rate.Limiter
	consecutiveFailures int
	lastFailure         time.Time
}

type LoadBalancer struct {
	backends  []*Backend
	algorithm balancer.Algorithm
	mu        sync.Mutex
}

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"backend"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"backend"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
}

func NewLoadBalancer(backends []string, algo string) *LoadBalancer {
	var be []*Backend
	for _, backend := range backends {
		url, _ := url.Parse(backend)
		be = append(be, &Backend{
			Backend: &balancer.Backend{URL: url, Healthy: true},
			limiter: rate.NewLimiter(rate.Limit(100), 200), // 100 requests per second, burst of 200
		})
	}

	var algorithm balancer.Algorithm
	switch algo {
	case "round-robin":
		algorithm = balancer.NewRoundRobin()
	case "least-connections":
		algorithm = &balancer.LeastConnections{}
	case "ip-hash":
		algorithm = &balancer.IPHash{}
	default:
		log.Printf("Unknown algorithm '%s', falling back to round-robin", algo)
		algorithm = balancer.NewRoundRobin()
	}

	return &LoadBalancer{backends: be, algorithm: algorithm}
}

func (lb *LoadBalancer) NextBackend(r *http.Request) *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	backend := lb.algorithm.NextBackend(lb.getHealthyBackends(), r)
	if backend == nil {
		return nil
	}
	return lb.findBackendByURL(backend.URL)
}

func (lb *LoadBalancer) getHealthyBackends() []*balancer.Backend {
	var healthy []*balancer.Backend
	for _, b := range lb.backends {
		if b.IsHealthy() {
			healthy = append(healthy, b.Backend)
		}
	}
	return healthy
}

func (lb *LoadBalancer) findBackendByURL(url *url.URL) *Backend {
	for _, b := range lb.backends {
		if b.URL.String() == url.String() {
			return b
		}
	}
	return nil
}

func (lb *LoadBalancer) HealthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		for _, backend := range lb.backends {
			status := "up"
			alive := isBackendAlive(backend.URL)
			backend.SetHealth(alive)
			if !alive {
				status = "down"
				backend.consecutiveFailures++
				backend.lastFailure = time.Now()
			} else {
				backend.consecutiveFailures = 0
			}
			log.Printf("Backend %s health check: %s\n", backend.URL, status)

			// Circuit breaking logic
			if backend.consecutiveFailures >= 3 && time.Since(backend.lastFailure) < 1*time.Minute {
				backend.SetHealth(false)
				log.Printf("Circuit breaker tripped for backend %s\n", backend.URL)
			}
		}
	}
}

func isBackendAlive(u *url.URL) bool {
	resp, err := http.Get(u.String() + "/health")
	if err != nil {
		log.Printf("Health check error for %s: %v\n", u, err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func main() {
	algorithm := flag.String("algorithm", "round-robin", "Load balancing algorithm (round-robin, least-connections, ip-hash)")
	flag.Parse()

	backendURLs := strings.Split(os.Getenv("BACKEND_URLS"), ",")
	if len(backendURLs) == 0 {
		backendURLs = []string{
			"http://localhost:8081",
			"http://localhost:8082",
		}
	}

	// Ensure all backend URLs have a scheme and host
	var validBackendURLs []string
	for _, backendURL := range backendURLs {
		if backendURL == "" {
			continue
		}
		if !strings.HasPrefix(backendURL, "http://") && !strings.HasPrefix(backendURL, "https://") {
			backendURL = "http://" + backendURL
		}
		parsedURL, err := url.Parse(backendURL)
		if err != nil || parsedURL.Host == "" {
			log.Printf("Invalid backend URL: %s, skipping", backendURL)
			continue
		}
		validBackendURLs = append(validBackendURLs, backendURL)
	}

	if len(validBackendURLs) == 0 {
		log.Fatal("No valid backend URLs provided")
	}

	lb := NewLoadBalancer(validBackendURLs, *algorithm)

	go lb.HealthCheck()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, false)
		fmt.Printf("Received request from %s\n%s\n", r.RemoteAddr, string(dump))

		backend := lb.NextBackend(r)
		if backend == nil {
			http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
			return
		}

		// Rate limiting
		if !backend.limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		proxy := httputil.NewSingleHostReverseProxy(backend.URL)
		proxy.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()

		requestsTotal.WithLabelValues(backend.URL.String()).Inc()
		requestDuration.WithLabelValues(backend.URL.String()).Observe(duration)

		fmt.Printf("Request forwarded to backend server: %s\n\n", backend.URL)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("cmd/lb/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Add Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server
	log.Println("Starting HTTPS server on :8443")
	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatalf("HTTPS server error: %v", err)
	}
}
