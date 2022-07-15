package v1

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) getNumAgreedTasks(w http.ResponseWriter, r *http.Request) {
	num, err := h.domain.GetNumAgreedTasks(r.Context())
	if err != nil {
		http.Error(w, "error when querying the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(num); err != nil {
		http.Error(w, "error during encoding to json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getNumRejectedTasks(w http.ResponseWriter, r *http.Request) {
	num, err := h.domain.GetNumRejectedTasks(r.Context())
	if err != nil {
		http.Error(w, "error when querying the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(num); err != nil {
		http.Error(w, "error during encoding to json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getTotalTime(w http.ResponseWriter, r *http.Request) {
	num, err := h.domain.GetTotalTime(r.Context())
	if err != nil {
		http.Error(w, "error when querying the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(num); err != nil {
		http.Error(w, "error during encoding to json", http.StatusInternalServerError)
		return
	}
}
