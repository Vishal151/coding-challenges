package balancer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadBalancer(t *testing.T) {
	// Create mock servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from server 1"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from server 2"))
	}))
	defer server2.Close()

	lb := NewLoadBalancer("8080", []string{server1.URL, server2.URL})

	// Test round-robin
	for i := 0; i < 4; i++ {
		backend := lb.getNextServerURL()
		expectedURL := server1.URL
		if i%2 == 1 {
			expectedURL = server2.URL
		}
		if backend.String() != expectedURL {
			t.Errorf("Request %d: expected %s, got %s", i, expectedURL, backend.String())
		}
	}
}
