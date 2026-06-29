package usecase

import "github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"

// IngestInput is the application-level request to add a document to the store.
type IngestInput struct {
	Content  string
	Metadata map[string]string
}

// IngestOutput reports the result of an ingest operation.
type IngestOutput struct {
	ChunksCreated int
}

// QueryInput is the application-level request to answer a question via RAG.
type QueryInput struct {
	Question string
	TopK     int
}

// QueryOutput carries the generated answer and its grounding sources.
type QueryOutput struct {
	Answer  string
	Sources []entity.ScoredDocument
}
