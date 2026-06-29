package vectorstore

import (
	"context"
	"testing"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
)

func TestMemorySearchRanksByCosineSimilarity(t *testing.T) {
	store := NewMemory()
	docs := []entity.Document{
		{ID: "a", Content: "apple"},
		{ID: "b", Content: "banana"},
		{ID: "c", Content: "cherry"},
	}
	vectors := [][]float32{
		{1, 0, 0},     // identical to the query direction
		{0, 1, 0},     // orthogonal
		{0.9, 0.1, 0}, // close to the query
	}
	if err := store.Add(context.Background(), docs, vectors); err != nil {
		t.Fatalf("Add: %v", err)
	}

	got, err := store.Search(context.Background(), []float32{1, 0, 0}, 2)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].ID != "a" {
		t.Errorf("expected top result 'a', got %q", got[0].ID)
	}
	if got[1].ID != "c" {
		t.Errorf("expected second result 'c', got %q", got[1].ID)
	}
}

func TestMemoryAddRejectsLengthMismatch(t *testing.T) {
	store := NewMemory()
	err := store.Add(context.Background(), []entity.Document{{ID: "a"}}, [][]float32{})
	if err == nil {
		t.Fatal("expected error on length mismatch")
	}
}
