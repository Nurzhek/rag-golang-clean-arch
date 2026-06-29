package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/dto"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
)

const (
	// jobsBasePath is used to build the poll URL returned on ingestion.
	jobsBasePath = "/api/v1/jobs/"
	// maxFileBytes caps the PUT file-upload body size.
	maxFileBytes = 32 << 20 // 32 MiB
)

// DocumentHandler exposes document ingestion and management over HTTP. It
// depends on use case interfaces, not concrete implementations.
type DocumentHandler struct {
	ingest usecase.IngestUseCase
	docs   usecase.DocumentUseCase
}

// NewDocumentHandler constructs a DocumentHandler.
func NewDocumentHandler(ingest usecase.IngestUseCase, docs usecase.DocumentUseCase) *DocumentHandler {
	return &DocumentHandler{ingest: ingest, docs: docs}
}

// Ingest handles POST /api/v1/documents. It accepts inline JSON content, queues
// asynchronous ingestion, and returns 200 immediately with a job ID to poll.
//
// @Summary     Queue inline document ingestion
// @Description Accepts inline text content, queues asynchronous embedding + storage, and returns a job ID to poll.
// @Tags        documents
// @Accept      json
// @Produce     json
// @Param       request body dto.IngestRequest true "Document content and optional metadata"
// @Success     200 {object} dto.JobAcceptedResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /api/v1/documents [post]
func (h *DocumentHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var req dto.IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	jobID, err := h.ingest.Submit(r.Context(), usecase.IngestInput{
		Content:  req.Content,
		Metadata: req.Metadata,
	})
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, accepted(jobID))
}

// IngestFile handles PUT /api/v1/documents. It accepts a raw file body and queues
// chunked asynchronous ingestion in the background, returning 202 with a job ID.
// Query parameters become document metadata, e.g. ?source=manual.txt&title=Manual.
//
// @Summary     Upload a file for chunked async ingestion
// @Description Accepts a raw request body (the file content) and queues chunked background ingestion. Query parameters are stored as document metadata.
// @Tags        documents
// @Accept      plain
// @Produce     json
// @Param       source query string false "Optional metadata: source name"
// @Param       title  query string false "Optional metadata: title"
// @Param       file   body  string true  "Raw file content"
// @Success     202 {object} dto.JobAcceptedResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     413 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /api/v1/documents [put]
func (h *DocumentHandler) IngestFile(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxFileBytes))
	if err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "file exceeds maximum allowed size")
		return
	}

	metadata := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			metadata[k] = v[0]
		}
	}

	jobID, err := h.ingest.Submit(r.Context(), usecase.IngestInput{
		Content:  string(body),
		Metadata: metadata,
	})
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusAccepted, accepted(jobID))
}

// List handles GET /api/v1/documents.
//
// @Summary     List ingested documents
// @Description Returns all ingested documents with their chunk counts and metadata.
// @Tags        documents
// @Produce     json
// @Success     200 {object} dto.ListDocumentsResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /api/v1/documents [get]
func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	docs, err := h.docs.List(r.Context())
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	resp := dto.ListDocumentsResponse{
		Count:     len(docs),
		Documents: make([]dto.DocumentSummary, 0, len(docs)),
	}
	for _, d := range docs {
		resp.Documents = append(resp.Documents, dto.DocumentSummary{
			ID:         d.ID,
			ChunkCount: d.ChunkCount,
			Metadata:   d.Metadata,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /api/v1/documents/{id}.
//
// @Summary     Delete a document
// @Description Deletes a document and all of its chunks by source document ID.
// @Tags        documents
// @Produce     json
// @Param       id path string true "Document ID"
// @Success     200 {object} dto.DeleteResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /api/v1/documents/{id} [delete]
func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	deleted, err := h.docs.Delete(r.Context(), id)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dto.DeleteResponse{DocumentID: id, DeletedChunks: deleted})
}

func accepted(jobID string) dto.JobAcceptedResponse {
	return dto.JobAcceptedResponse{
		JobID:   jobID,
		Status:  "queued",
		PollURL: jobsBasePath + jobID,
	}
}
