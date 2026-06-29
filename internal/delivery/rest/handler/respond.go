package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/dto"
)

// writeJSON encodes payload as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("encode response", "error", err)
	}
}

// writeError writes a standard JSON error envelope.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, dto.ErrorResponse{Error: msg})
}
