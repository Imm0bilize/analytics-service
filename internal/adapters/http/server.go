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

func (s *Server) createHandler(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()

	r.Use(middlewares...)

	r.Route("/debug", func(r chi.Router) {
		r.Get("/healthz", s.healthCheck)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/num-agreed", s.getNumAgreedTasks)
			r.Get("/num-rejected", s.getNumRejectedTasks)
			r.Get("/total-time", s.getTotalTime)
		})
	})

	return r
}

func New(port string, readTimeout, writeTimeout time.Duration, middlewares ...func(http.Handler) http.Handler) *Server {
	s := &Server{
		notify: make(chan error, 1),
	}

	server := &http.Server{
		Handler:      s.createHandler(middlewares...),
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
