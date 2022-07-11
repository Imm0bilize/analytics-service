package http

import "net/http"

func (s *Server) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		// add logging
		return
	}
}
