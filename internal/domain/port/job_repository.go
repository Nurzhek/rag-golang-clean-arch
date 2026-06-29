package port

import (
	"context"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
)

// JobRepository is the port for persisting and retrieving asynchronous
// ingestion job state.
type JobRepository interface {
	// Save upserts a job by its ID.
	Save(ctx context.Context, job entity.Job) error
	// Get returns a job by ID, or domain.ErrJobNotFound if it does not exist.
	Get(ctx context.Context, id string) (entity.Job, error)
}
