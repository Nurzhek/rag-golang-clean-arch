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
//
// @Summary     Poll ingestion job status
// @Description Returns the current status and progress of an asynchronous ingestion job.
// @Tags        jobs
// @Produce     json
// @Param       id path string true "Job ID"
// @Success     200 {object} dto.JobResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /api/v1/jobs/{id} [get]
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
