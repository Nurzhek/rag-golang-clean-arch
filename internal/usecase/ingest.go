package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// IDGenerator produces unique IDs for stored chunks. It is injected so the use
// case stays deterministic and easy to test.
type IDGenerator func() string

// ingestInteractor implements IngestUseCase. It orchestrates the ingestion
// pipeline (split -> embed -> store) and depends only on domain ports, so any
// splitter/embedder/store implementation can be plugged in.
type ingestInteractor struct {
	splitter port.TextSplitter
	embedder port.Embedder
	store    port.VectorStore
	newID    IDGenerator
}

// NewIngestUseCase wires the ingestion dependencies.
func NewIngestUseCase(splitter port.TextSplitter, embedder port.Embedder, store port.VectorStore, newID IDGenerator) IngestUseCase {
	return &ingestInteractor{splitter: splitter, embedder: embedder, store: store, newID: newID}
}

func (uc *ingestInteractor) Execute(ctx context.Context, in IngestInput) (IngestOutput, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return IngestOutput{}, domain.ErrEmptyContent
	}

	chunks, err := uc.splitter.Split(content)
	if err != nil {
		return IngestOutput{}, fmt.Errorf("split content: %w", err)
	}
	if len(chunks) == 0 {
		return IngestOutput{}, domain.ErrEmptyContent
	}

	vectors, err := uc.embedder.EmbedDocuments(ctx, chunks)
	if err != nil {
		return IngestOutput{}, fmt.Errorf("embed chunks: %w", err)
	}

	docs := make([]entity.Document, len(chunks))
	for i, chunk := range chunks {
		docs[i] = entity.Document{
			ID:       uc.newID(),
			Content:  chunk,
			Metadata: cloneMetadata(in.Metadata),
		}
	}

	if err := uc.store.Add(ctx, docs, vectors); err != nil {
		return IngestOutput{}, fmt.Errorf("store documents: %w", err)
	}

	return IngestOutput{ChunksCreated: len(docs)}, nil
}

func cloneMetadata(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
