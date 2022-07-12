package http

import "net/http"

func (s *Server) getNumAgreedTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (s *Server) getNumRejectedTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (s *Server) getTotalTime(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
