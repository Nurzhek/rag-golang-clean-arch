package usecase

import "context"

// IngestUseCase ingests source documents into the knowledge base.
type IngestUseCase interface {
	Execute(ctx context.Context, in IngestInput) (IngestOutput, error)
}

// QueryUseCase answers a natural-language question using retrieval-augmented
// generation over the knowledge base.
type QueryUseCase interface {
	Execute(ctx context.Context, in QueryInput) (QueryOutput, error)
}
