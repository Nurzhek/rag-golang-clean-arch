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

func TestMemoryListAndDeleteBySource(t *testing.T) {
	store := NewMemory()
	docs := []entity.Document{
		{ID: "c1", SourceID: "a"},
		{ID: "c2", SourceID: "a"},
		{ID: "c3", SourceID: "b"},
	}
	vectors := [][]float32{{1, 0}, {0, 1}, {1, 1}}
	if err := store.Add(context.Background(), docs, vectors); err != nil {
		t.Fatalf("Add: %v", err)
	}

	list, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(list))
	}

	deleted, err := store.Delete(context.Background(), "a")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if deleted != 2 {
		t.Errorf("expected 2 deleted, got %d", deleted)
	}

	remaining, _ := store.List(context.Background())
	if len(remaining) != 1 {
		t.Errorf("expected 1 remaining chunk, got %d", len(remaining))
	}

	if n, _ := store.Delete(context.Background(), "missing"); n != 0 {
		t.Errorf("expected 0 deleted for missing source, got %d", n)
	}
}
