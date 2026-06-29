package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/dto"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
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

// writeUseCaseError maps domain errors onto HTTP status codes. Centralised so
// every handler reports failures consistently (DRY).
func writeUseCaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyContent), errors.Is(err, domain.ErrEmptyQuestion):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrDocumentNotFound),
		errors.Is(err, domain.ErrJobNotFound),
		errors.Is(err, domain.ErrNoRelevantDocuments):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
