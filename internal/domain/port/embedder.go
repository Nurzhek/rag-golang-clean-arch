package port

import "context"

// Embedder is the port for turning text into dense vector representations.
type Embedder interface {
	// EmbedDocuments embeds a batch of texts, preserving order.
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
	// EmbedQuery embeds a single query string.
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
}
