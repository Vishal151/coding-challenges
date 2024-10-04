package balancer

import (
	"net/url"
	"sync"
)

// Backend represents a backend server
type Backend struct {
	URL               *url.URL
	Healthy           bool
	ActiveConnections int
	mux               sync.RWMutex
}

func (b *Backend) IsHealthy() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Healthy
}

func (b *Backend) SetHealth(health bool) {
	b.mux.Lock()
	b.Healthy = health
	b.mux.Unlock()
}

type Algorithm interface {
	NextBackend(backends []*Backend) *Backend
}

type RoundRobin struct {
	current int
	mu      sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{current: -1}
}

func (rr *RoundRobin) NextBackend(backends []*Backend) *Backend {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if len(backends) == 0 {
		return nil
	}
	rr.current = (rr.current + 1) % len(backends)
	return backends[rr.current]
}

func (rr *RoundRobin) ResetCurrent() {
	rr.mu.Lock()
	rr.current = -1
	rr.mu.Unlock()
}

type LeastConnections struct{}

func (lc *LeastConnections) NextBackend(backends []*Backend) *Backend {
	var leastConn *Backend
	minConn := int(^uint(0) >> 1) // Max int

	for _, b := range backends {
		if b.ActiveConnections < minConn && b.IsHealthy() {
			minConn = b.ActiveConnections
			leastConn = b
		}
	}

	return leastConn
}

type IPHash struct{}

func (ih *IPHash) NextBackend(backends []*Backend) *Backend {
	// Note: This implementation doesn't use the request, so it's not a true IP hash.
	// For demonstration purposes, we'll just return the first healthy backend.
	for _, b := range backends {
		if b.IsHealthy() {
			return b
		}
	}
	return nil
}
