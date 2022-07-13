package v1

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func CreateHandler(auth func(http.Handler) http.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()

	r.Use(middlewares...)

	r.Route("/debug", func(r chi.Router) {
		r.Get("/healthz", healthCheck)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(auth)
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/num-agreed", getNumAgreedTasks)
			r.Get("/num-rejected", getNumRejectedTasks)
			r.Get("/total-time", getTotalTime)
		})
	})

	return r
}
