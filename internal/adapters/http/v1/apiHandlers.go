package v1

import (
	"encoding/json"
	"net/http"
)

// @Summary      getting the number of accepted tasks
// @Description  the handler allows you to get the total number of accepted tasks stored in the database
// @ID           get-num-accepted-tasks
// @Tags         tasks
// @Security     ApiKeyAuth
// @Accept       json
// @Success      200  {object}  dto.NumAgreedTasksDTO
// @Failure      500  {string}  string  "error when querying the database"
// @Failure      500  {string}  string  "error during encoding to json"
// @Security     ApiKeyAuth
// @Router       /api/tasks/num-accepted [get]
func (h *Handler) getNumAcceptedTasks(w http.ResponseWriter, r *http.Request) {
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

// @Summary      getting the number of rejected tasks
// @Description  the handler allows you to get the total number of rejected tasks stored in the database
// @ID           get-num-rejected-tasks
// @Tags         tasks
// @Accept       json
// @Success      200  {object}  dto.NumRejectedTaskDTO
// @Failure      500  {string}  string  "error when querying the database"
// @Failure      500  {string}  string  "error during encoding to json"
// @Security     ApiKeyAuth
// @Router       /api/tasks/num-rejected [get]
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

// @Summary      getting the total time for all tasks
// @Description  the handler allows you to get the total amount of time spent on confirmed or rejected tasks
// @ID           get-total-time
// @Tags         tasks
// @Accept       json
// @Success      200  {object}  dto.TotalTimeDTO
// @Failure      500  {string}  string  "error when querying the database"
// @Failure      500  {string}  string  "error during encoding to json"
// @Security     ApiKeyAuth
// @Router       /api/tasks/total-time [get]
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
