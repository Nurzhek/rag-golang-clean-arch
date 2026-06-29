package usecase

import (
	"context"
	"fmt"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// documentInteractor implements DocumentUseCase. It presents stored chunks at
// the document (source) level by grouping on SourceID.
type documentInteractor struct {
	store port.VectorStore
}

// NewDocumentUseCase wires the document management dependencies.
func NewDocumentUseCase(store port.VectorStore) DocumentUseCase {
	return &documentInteractor{store: store}
}

// List groups stored chunks by source document, preserving first-seen order.
func (uc *documentInteractor) List(ctx context.Context) ([]entity.DocumentSummary, error) {
	chunks, err := uc.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}

	index := make(map[string]int)
	summaries := make([]entity.DocumentSummary, 0)
	for _, c := range chunks {
		if i, ok := index[c.SourceID]; ok {
			summaries[i].ChunkCount++
			continue
		}
		index[c.SourceID] = len(summaries)
		summaries = append(summaries, entity.DocumentSummary{
			ID:         c.SourceID,
			ChunkCount: 1,
			Metadata:   c.Metadata,
		})
	}
	return summaries, nil
}

// Delete removes a source document and all of its chunks. It returns
// domain.ErrDocumentNotFound when the ID matches nothing.
func (uc *documentInteractor) Delete(ctx context.Context, documentID string) (int, error) {
	deleted, err := uc.store.Delete(ctx, documentID)
	if err != nil {
		return 0, fmt.Errorf("delete document: %w", err)
	}
	if deleted == 0 {
		return 0, domain.ErrDocumentNotFound
	}
	return deleted, nil
}
