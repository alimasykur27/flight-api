package httpserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server: Represent an HTTP server
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP Server
func NewServer(router *chi.Mux, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
