package handler

import (
	"net/http"
	"time"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/dto"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
)

// JobHandler exposes asynchronous job status — the polling endpoint.
type JobHandler struct {
	jobs usecase.JobUseCase
}

// NewJobHandler constructs a JobHandler.
func NewJobHandler(jobs usecase.JobUseCase) *JobHandler {
	return &JobHandler{jobs: jobs}
}

// Get handles GET /api/v1/jobs/{id}.
func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	job, err := h.jobs.Get(r.Context(), id)
	if err != nil {
		writeUseCaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.JobResponse{
		ID:              job.ID,
		Status:          string(job.Status),
		TotalChunks:     job.TotalChunks,
		ProcessedChunks: job.ProcessedChunks,
		DocumentID:      job.DocumentID,
		Error:           job.Error,
		CreatedAt:       job.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       job.UpdatedAt.Format(time.RFC3339),
	})
}
