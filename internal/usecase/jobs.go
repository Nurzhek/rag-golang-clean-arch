package usecase

import (
	"context"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// jobInteractor implements JobUseCase — the read side of the async pipeline.
type jobInteractor struct {
	jobs port.JobRepository
}

// NewJobUseCase wires the job query dependencies.
func NewJobUseCase(jobs port.JobRepository) JobUseCase {
	return &jobInteractor{jobs: jobs}
}

// Get returns the current state of a job (or domain.ErrJobNotFound).
func (uc *jobInteractor) Get(ctx context.Context, jobID string) (entity.Job, error) {
	return uc.jobs.Get(ctx, jobID)
}
