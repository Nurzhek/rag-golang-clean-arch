package embedding

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// Config holds the settings needed to construct the OpenAI embedder adapter.
type Config struct {
	APIKey  string
	Model   string
	BaseURL string
}

// openAIEmbedder adapts a langchaingo embedder to the domain port.Embedder.
type openAIEmbedder struct {
	embedder embeddings.Embedder
}

// compile-time check that openAIEmbedder satisfies the port.
var _ port.Embedder = (*openAIEmbedder)(nil)

// NewOpenAI builds a port.Embedder backed by langchaingo's OpenAI embeddings.
func NewOpenAI(cfg Config) (port.Embedder, error) {
	opts := []openai.Option{
		openai.WithToken(cfg.APIKey),
		openai.WithEmbeddingModel(cfg.Model),
	}
	if cfg.BaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
	}

	client, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("init openai embedding client: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, fmt.Errorf("init embedder: %w", err)
	}

	return &openAIEmbedder{embedder: embedder}, nil
}

func (e *openAIEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	vecs, err := e.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("embed documents: %w", err)
	}
	return vecs, nil
}

func (e *openAIEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	vec, err := e.embedder.EmbedQuery(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}
	return vec, nil
}
