package handler

import "net/http"

// Health handles GET /health and reports service liveness.
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
