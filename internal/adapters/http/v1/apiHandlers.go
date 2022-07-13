package v1

import (
	"net/http"
)

func getNumAgreedTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func getNumRejectedTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func getTotalTime(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
