package logging

import (
	"fmt"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (l *logger) MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         200,
		}

		start := time.Now()
		next.ServeHTTP(recorder, r)
		end := time.Now()

		line := fmt.Sprintf("type:%s path:%s, el_time:%s", r.Method, r.RequestURI, end.Sub(start))
		switch {
		case recorder.status < 400:
			l.InfoF(line)
		case recorder.status >= 400 && recorder.status < 500:
			l.WarningF(line)
		case recorder.status >= 500:
			l.ErrorF(line)
		}
	})
}
