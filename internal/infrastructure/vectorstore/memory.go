package vectorstore

import (
	"context"
	"errors"
	"math"
	"sort"
	"sync"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// ErrDimensionMismatch is returned when the documents and vectors slices passed
// to Add are not the same length.
var ErrDimensionMismatch = errors.New("documents and vectors length mismatch")

type record struct {
	doc    entity.Document
	vector []float32
}

// Memory is a dependency-free, in-process vector store backed by brute-force
// cosine similarity. It is ideal for development, tests, and small datasets.
// For production scale, provide another port.VectorStore implementation
// (pgvector, Qdrant, Weaviate, ...) — nothing else in the system changes.
type Memory struct {
	mu      sync.RWMutex
	records []record
}

// compile-time check that Memory satisfies the port.
var _ port.VectorStore = (*Memory)(nil)

// NewMemory creates an empty in-memory vector store.
func NewMemory() *Memory {
	return &Memory{}
}

// Add stores documents alongside their pre-computed vectors.
func (m *Memory) Add(_ context.Context, docs []entity.Document, vectors [][]float32) error {
	if len(docs) != len(vectors) {
		return ErrDimensionMismatch
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range docs {
		m.records = append(m.records, record{doc: docs[i], vector: vectors[i]})
	}
	return nil
}

// Search returns the topK most similar documents, highest score first.
func (m *Memory) Search(_ context.Context, queryVector []float32, topK int) ([]entity.ScoredDocument, error) {
	if topK <= 0 {
		return nil, nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	scored := make([]entity.ScoredDocument, 0, len(m.records))
	for _, r := range m.records {
		scored = append(scored, entity.ScoredDocument{
			Document: r.doc,
			Score:    cosineSimilarity(queryVector, r.vector),
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if len(scored) > topK {
		scored = scored[:topK]
	}
	return scored, nil
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		av, bv := float64(a[i]), float64(b[i])
		dot += av * bv
		normA += av * av
		normB += bv * bv
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return float32(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}
