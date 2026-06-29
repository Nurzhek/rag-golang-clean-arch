package usecase

import (
	"context"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
)

// IngestUseCase queues a document for asynchronous ingestion and returns the ID
// of the job that tracks its progress.
type IngestUseCase interface {
	Submit(ctx context.Context, in IngestInput) (jobID string, err error)
}

// DocumentUseCase manages stored documents at the source-document level.
type DocumentUseCase interface {
	List(ctx context.Context) ([]entity.DocumentSummary, error)
	Delete(ctx context.Context, documentID string) (deletedChunks int, err error)
}

// JobUseCase exposes the status of asynchronous ingestion jobs (the polling API).
type JobUseCase interface {
	Get(ctx context.Context, jobID string) (entity.Job, error)
}

// QueryUseCase answers a natural-language question using retrieval-augmented
// generation over the knowledge base.
type QueryUseCase interface {
	Execute(ctx context.Context, in QueryInput) (QueryOutput, error)
}
