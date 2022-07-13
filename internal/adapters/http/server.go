package http

import (
	"analytic-service/internal/config"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
)

type Server struct {
	notify chan error
	server *http.Server
}

func (s *Server) createHandler(auth func(http.Handler) http.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()

	r.Use(middlewares...)

	r.Route("/debug", func(r chi.Router) {
		r.Get("/healthz", s.healthCheck)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(auth)
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/num-agreed", s.getNumAgreedTasks)
			r.Get("/num-rejected", s.getNumRejectedTasks)
			r.Get("/total-time", s.getTotalTime)
		})
	})

	return r
}

func New(
	cfg *config.Config, auth func(http.Handler) http.Handler, middlewares ...func(http.Handler) http.Handler) *Server {
	s := &Server{
		notify: make(chan error, 1),
	}

	server := &http.Server{
		Handler:      s.createHandler(auth, middlewares...),
		Addr:         net.JoinHostPort("", cfg.Http.Port),
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
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
