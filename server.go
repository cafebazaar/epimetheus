package epimetheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(listenPort int) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: promhttp.Handler(),
	}

	return &Server{
		httpServer: server,
	}
}

func (s *Server) Serve() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return errors.Wrapf(err, "Failed to start Prometheus listener on %v", s.httpServer.Addr)
	}
	return nil
}

func (s *Server) Stop(server *http.Server) error {
	if err := server.Shutdown(context.Background()); err != nil {
		return errors.Wrap(err, "Failed to shutdown Prometheus listener")
	}
	return nil
}
