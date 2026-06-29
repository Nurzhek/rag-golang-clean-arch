package usecase_test

import (
	"context"
	"testing"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
)

func TestDocumentUseCaseListGroupsBySource(t *testing.T) {
	store := fakeStore{list: []entity.Document{
		{ID: "c1", SourceID: "doc-a", Content: "x", Metadata: map[string]string{"title": "A"}},
		{ID: "c2", SourceID: "doc-a", Content: "y"},
		{ID: "c3", SourceID: "doc-b", Content: "z"},
	}}

	uc := usecase.NewDocumentUseCase(store)

	got, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(got))
	}
	if got[0].ID != "doc-a" || got[0].ChunkCount != 2 {
		t.Errorf("unexpected first document: %+v", got[0])
	}
	if got[1].ID != "doc-b" || got[1].ChunkCount != 1 {
		t.Errorf("unexpected second document: %+v", got[1])
	}
}

func TestDocumentUseCaseDeleteExisting(t *testing.T) {
	store := fakeStore{list: []entity.Document{
		{ID: "c1", SourceID: "doc-a"},
		{ID: "c2", SourceID: "doc-a"},
	}}

	uc := usecase.NewDocumentUseCase(store)

	n, err := uc.Delete(context.Background(), "doc-a")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 deleted chunks, got %d", n)
	}
}

func TestDocumentUseCaseDeleteMissingReturnsNotFound(t *testing.T) {
	uc := usecase.NewDocumentUseCase(fakeStore{})

	_, err := uc.Delete(context.Background(), "nope")
	if err != domain.ErrDocumentNotFound {
		t.Errorf("expected ErrDocumentNotFound, got %v", err)
	}
}
