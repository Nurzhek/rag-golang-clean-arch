package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/dto"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
)

// RAGHandler exposes the ingest and query use cases over HTTP. It depends on the
// use case interfaces, not their concrete implementations (dependency inversion).
type RAGHandler struct {
	ingest usecase.IngestUseCase
	query  usecase.QueryUseCase
}

// NewRAGHandler constructs a RAGHandler.
func NewRAGHandler(ingest usecase.IngestUseCase, query usecase.QueryUseCase) *RAGHandler {
	return &RAGHandler{ingest: ingest, query: query}
}

// Ingest handles POST /api/v1/documents.
func (h *RAGHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var req dto.IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	out, err := h.ingest.Execute(r.Context(), usecase.IngestInput{
		Content:  req.Content,
		Metadata: req.Metadata,
	})
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.IngestResponse{ChunksCreated: out.ChunksCreated})
}

// Query handles POST /api/v1/query.
func (h *RAGHandler) Query(w http.ResponseWriter, r *http.Request) {
	var req dto.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	out, err := h.query.Execute(r.Context(), usecase.QueryInput{
		Question: req.Question,
		TopK:     req.TopK,
	})
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	resp := dto.QueryResponse{
		Answer:  out.Answer,
		Sources: make([]dto.Source, 0, len(out.Sources)),
	}
	for _, s := range out.Sources {
		resp.Sources = append(resp.Sources, dto.Source{
			ID:       s.ID,
			Content:  s.Content,
			Score:    s.Score,
			Metadata: s.Metadata,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

// writeUseCaseError maps domain errors onto HTTP status codes.
func writeUseCaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyContent), errors.Is(err, domain.ErrEmptyQuestion):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNoRelevantDocuments):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
