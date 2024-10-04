package balancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []string
	mu              sync.Mutex
}

func NewLoadBalancer(port string, servers []string) *LoadBalancer {
	return &LoadBalancer{
		port:    port,
		servers: servers,
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextServerURL()
	proxy := httputil.NewSingleHostReverseProxy(targetServer)
	proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) getNextServerURL() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	lb.roundRobinCount++
	serverURL, _ := url.Parse(server)
	return serverURL
}

func (lb *LoadBalancer) Start() error {
	return http.ListenAndServe(":"+lb.port, lb)
}
