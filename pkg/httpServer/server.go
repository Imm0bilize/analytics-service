package httpServer

import (
	"analytic-service/internal/config"
	"context"
	"errors"
	"net"
	"net/http"
)

type Server struct {
	notify chan error
	server *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	s := &Server{
		notify: make(chan error, 1),
		server: &http.Server{
			Handler:      handler,
			Addr:         net.JoinHostPort("", cfg.Http.Port),
			ReadTimeout:  cfg.Http.ReadTimeout,
			WriteTimeout: cfg.Http.WriteTimeout,
		},
	}
	return s
}

func (s *Server) Run() {
	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
			close(s.notify)
		}
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
