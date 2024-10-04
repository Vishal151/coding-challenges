package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vishal151/load-balancer/internal/balancer"
)

func TestLoadBalancer(t *testing.T) {
	lb := NewLoadBalancer([]string{
		"http://localhost:8081",
		"http://localhost:8082",
	}, "round-robin")

	// Reset the RoundRobin algorithm
	if rr, ok := lb.algorithm.(*balancer.RoundRobin); ok {
		rr.ResetCurrent()
	}

	// Test round-robin
	expectedPorts := []int{8081, 8082, 8081, 8082}
	for i, expectedPort := range expectedPorts {
		backend := lb.NextBackend()
		if backend == nil {
			t.Fatalf("Expected a backend, got nil")
		}
		if backend.URL.Port() != fmt.Sprint(expectedPort) {
			t.Errorf("Request %d: Expected backend port %d, got %s", i, expectedPort, backend.URL.Port())
		}
	}

	// Test health check
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	lb = NewLoadBalancer([]string{server.URL}, "round-robin")
	go lb.HealthCheck()

	// Wait for health check
	time.Sleep(11 * time.Second)

	if !lb.backends[0].IsHealthy() {
		t.Errorf("Expected backend to be healthy")
	}
}
