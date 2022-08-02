package prometheus

import (
	"analytic-service/pkg/responseWriter"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

func (m *metrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		interceptor := &responseWriter.Interceptor{
			StatusCode:     http.StatusOK,
			ResponseWriter: w,
		}
		timer := prometheus.NewTimer(m.httpDuration.WithLabelValues(r.RequestURI))

		next.ServeHTTP(interceptor, r)

		timer.ObserveDuration()
		m.opsProcessed.With(
			prometheus.Labels{
				"method":     r.Method,
				"path":       r.RequestURI,
				"statuscode": strconv.Itoa(interceptor.StatusCode),
			},
		).Inc()
	})
}
