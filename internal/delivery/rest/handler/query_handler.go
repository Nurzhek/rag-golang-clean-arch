package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/dto"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
)

// QueryHandler exposes the RAG query use case over HTTP.
type QueryHandler struct {
	query usecase.QueryUseCase
}

// NewQueryHandler constructs a QueryHandler.
func NewQueryHandler(query usecase.QueryUseCase) *QueryHandler {
	return &QueryHandler{query: query}
}

// Query handles POST /api/v1/query.
//
// @Summary     Ask a question (RAG)
// @Description Embeds the question, retrieves the most relevant chunks, and generates a grounded answer with sources.
// @Tags        query
// @Accept      json
// @Produce     json
// @Param       request body dto.QueryRequest true "Question and optional top_k"
// @Success     200 {object} dto.QueryResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /api/v1/query [post]
func (h *QueryHandler) Query(w http.ResponseWriter, r *http.Request) {
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
			ID:         s.ID,
			DocumentID: s.SourceID,
			Content:    s.Content,
			Score:      s.Score,
			Metadata:   s.Metadata,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}
