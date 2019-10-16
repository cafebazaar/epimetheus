package epimetheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server contains http server of Epimetheus instance
type Server struct {
	httpServer *http.Server
}

func newServer(listenPort int) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: promhttp.Handler(),
	}

	return &Server{
		httpServer: server,
	}
}

func (s *Server) serve() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return errors.Wrapf(err, "Failed to start Prometheus listener on %v", s.httpServer.Addr)
	}
	return nil
}

// Stop stops the server and returns it's error if exists
func (s *Server) Stop(server *http.Server) error {
	if err := server.Shutdown(context.Background()); err != nil {
		return errors.Wrap(err, "Failed to shutdown Prometheus listener")
	}
	return nil
}
