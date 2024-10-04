package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadBalancer(t *testing.T) {
	lb := NewLoadBalancer([]string{
		"http://localhost:8081",
		"http://localhost:8082",
	})

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
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	lb = NewLoadBalancer([]string{server.URL})
	go lb.HealthCheck()

	// Wait for health check
	time.Sleep(6 * time.Second)

	if !lb.backends[0].IsHealthy() {
		t.Errorf("Expected backend to be healthy")
	}
}
