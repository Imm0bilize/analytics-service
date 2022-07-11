package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"time"
)

type Server struct {
	notify chan error
	server *http.Server
}

func (s *Server) createHandler() *chi.Mux {
	r := chi.NewMux()

	r.Route("/debug", func(r chi.Router) {
		r.Get("/healthz", s.healthCheckHandler)
	})

	return r
}

func New(port string, readTimeout, writeTimeout time.Duration) *Server {
	s := &Server{
		notify: make(chan error, 1),
	}

	server := &http.Server{
		Handler:      s.createHandler(),
		Addr:         net.JoinHostPort("", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	s.server = server

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
