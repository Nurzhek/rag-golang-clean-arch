package port

import (
	"context"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
)

// VectorStore is the port for persisting embedded documents and retrieving the
// most similar ones for a query vector. It is intentionally agnostic of how
// embeddings are produced — callers embed first, then store or search.
type VectorStore interface {
	// Add stores documents alongside their pre-computed vectors. The slices
	// must be the same length and aligned by index.
	Add(ctx context.Context, docs []entity.Document, vectors [][]float32) error
	// Search returns the topK documents most similar to queryVector, ordered
	// by descending similarity.
	Search(ctx context.Context, queryVector []float32, topK int) ([]entity.ScoredDocument, error)
}
