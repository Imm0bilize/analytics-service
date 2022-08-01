package prometheus

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
)

type metrics struct {
	notify       chan error
	port         string
	opsProcessed *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
}

func New(port string) *metrics {
	opsProcessed := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "team29",
		Subsystem: "analytics",
		Name:      "processed_ops_total",
		Help:      "The total number of processed events",
	}, []string{"method", "path", "statuscode"})

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "team29",
		Subsystem: "analytics",
		Name:      "http_response_time_seconds",
		Help:      "Duration of HTTP requests.",
	}, []string{"path"})

	notify := make(chan error, 1)
	return &metrics{
		notify:       notify,
		port:         port,
		opsProcessed: opsProcessed,
		httpDuration: httpDuration,
	}
}

func (m *metrics) Run() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(net.JoinHostPort("", m.port), nil); !errors.Is(err, http.ErrServerClosed) {
			m.notify <- err
			close(m.notify)
		}
	}()
}

func (m *metrics) Notify() <-chan error {
	return m.notify
}
