package logging

import (
	"analytic-service/pkg/responseWriter"
	"fmt"
	"net/http"
	"time"
)

func (l *logger) MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		interceptor := &responseWriter.Interceptor{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(interceptor, r)
		end := time.Now()

		line := fmt.Sprintf("type:%s path:%s, el_time:%s", r.Method, r.RequestURI, end.Sub(start))
		switch {
		case interceptor.StatusCode < 400:
			l.Infof(line)
		case interceptor.StatusCode >= 400 && interceptor.StatusCode < 500:
			l.Warningf(line)
		case interceptor.StatusCode >= 500:
			l.Errorf(line)
		}
	})
}
