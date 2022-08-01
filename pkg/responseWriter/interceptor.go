package responseWriter

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type Interceptor struct {
	http.ResponseWriter
	StatusCode int
}

func (i *Interceptor) WriteHeader(statusCode int) {
	i.StatusCode = statusCode
	i.ResponseWriter.WriteHeader(statusCode)
}

func (i *Interceptor) Write(p []byte) (int, error) {
	return i.ResponseWriter.Write(p)
}

func (i *Interceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := i.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("type assertion failed http.ResponseWriter not a http.Hijacker")
	}
	return h.Hijack()
}

func (i *Interceptor) Flush() {
	f, ok := i.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	f.Flush()
}
