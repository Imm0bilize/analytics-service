package http

import "net/http"

// DEBUG

func (s *Server) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		// add logging
		return
	}
}

// API
