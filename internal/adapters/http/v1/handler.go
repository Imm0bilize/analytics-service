package v1

import (
	"analytic-service/internal/ports"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Handler struct {
	domain ports.ClientDomain
}

func CreateHandler(domain ports.ClientDomain) *Handler {
	return &Handler{domain: domain}
}

func (h *Handler) GetHttpHandler(auth func(http.Handler) http.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()

	r.Use(middlewares...)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/debug", func(r chi.Router) {
		r.Get("/healthz", healthCheck)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(auth)
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/num-accepted", h.getNumAcceptedTasks)
			r.Get("/num-rejected", h.getNumRejectedTasks)
			r.Get("/total-time", h.getTotalTime)
		})
	})

	return r
}
