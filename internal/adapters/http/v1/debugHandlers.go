package v1

import (
	"net/http"
)

// HealthCheck godoc
// @Summary
// @Description  check service health
// @ID           health-check
// @Tags         debug
// @Success      200  {string}  string  "ok"
// @Router       /debug/healthz [get]
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		// add logging
		return
	}
}
