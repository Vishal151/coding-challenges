package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from server %s\n", s.addr)
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s)
}
