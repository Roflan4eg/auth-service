package http

import (
	"context"
	"net/http"

	"github.com/Roflan4eg/auth-serivce/internal/interfaces/http/handlers"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
)

type Server struct {
	httpServer *http.Server
	port       string
	log        *logger.Logger
	name       string
}

func NewServer(
	handlers *handlers.Container,
	logger *logger.Logger,
	port string,
) *Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", handlers.Metrics)

	// mux.Handle("/health", handlers.HealthHandler())

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return &Server{
		httpServer: server,
		port:       port,
		log:        logger,
		name:       "metrics-server",
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting HTTP server", "port", s.port, "name", s.name)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Stopping HTTP server", "name", s.name)
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Name() string {
	return s.name
}
