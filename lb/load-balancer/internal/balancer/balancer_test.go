package balancer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadBalancer(t *testing.T) {
	servers := []string{"http://server1", "http://server2"}
	lb := NewLoadBalancer("8080", servers)

	// Test round-robin behavior
	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		lb.ServeHTTP(rr, req)

		expectedServer := servers[i%len(servers)]
		if lb.getNextServerURL().String() != expectedServer {
			t.Errorf("Request %d: expected %s, got %s", i, expectedServer, lb.getNextServerURL().String())
		}
	}
}
