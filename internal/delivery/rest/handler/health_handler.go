package handler

import "net/http"

// Health handles GET /health and reports service liveness.
//
// @Summary     Liveness check
// @Description Reports service liveness.
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
