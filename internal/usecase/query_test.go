package usecase_test

import (
	"context"
	"testing"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
)

// The fakes below implement the domain ports, letting the use cases be tested in
// full isolation — no OpenAI calls, no real vector store.

type fakeEmbedder struct{}

func (fakeEmbedder) EmbedDocuments(_ context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i := range texts {
		out[i] = []float32{1, 0}
	}
	return out, nil
}

func (fakeEmbedder) EmbedQuery(_ context.Context, _ string) ([]float32, error) {
	return []float32{1, 0}, nil
}

type fakeStore struct {
	search []entity.ScoredDocument
	list   []entity.Document
}

func (f fakeStore) Add(_ context.Context, _ []entity.Document, _ [][]float32) error { return nil }

func (f fakeStore) Search(_ context.Context, _ []float32, topK int) ([]entity.ScoredDocument, error) {
	if topK > len(f.search) {
		topK = len(f.search)
	}
	return f.search[:topK], nil
}

func (f fakeStore) List(_ context.Context) ([]entity.Document, error) { return f.list, nil }

func (f fakeStore) Delete(_ context.Context, sourceID string) (int, error) {
	n := 0
	for _, d := range f.list {
		if d.SourceID == sourceID {
			n++
		}
	}
	return n, nil
}

type fakeLLM struct{ lastPrompt string }

func (f *fakeLLM) Generate(_ context.Context, prompt string) (string, error) {
	f.lastPrompt = prompt
	return "the answer", nil
}

func TestQueryUseCaseReturnsAnswerAndSources(t *testing.T) {
	store := fakeStore{search: []entity.ScoredDocument{
		{Document: entity.Document{ID: "1", Content: "ground truth"}, Score: 0.9},
	}}
	llm := &fakeLLM{}

	uc := usecase.NewQueryUseCase(fakeEmbedder{}, store, llm, usecase.DefaultPromptBuilder, 4)

	out, err := uc.Execute(context.Background(), usecase.QueryInput{Question: "what is the truth?"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.Answer != "the answer" {
		t.Errorf("unexpected answer: %q", out.Answer)
	}
	if len(out.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(out.Sources))
	}
	if llm.lastPrompt == "" {
		t.Error("expected the LLM to receive a built prompt")
	}
}

func TestQueryUseCaseRejectsEmptyQuestion(t *testing.T) {
	uc := usecase.NewQueryUseCase(fakeEmbedder{}, fakeStore{}, &fakeLLM{}, usecase.DefaultPromptBuilder, 4)

	_, err := uc.Execute(context.Background(), usecase.QueryInput{Question: "   "})
	if err == nil {
		t.Fatal("expected error for empty question")
	}
}

func TestQueryUseCaseReturnsErrWhenNoSources(t *testing.T) {
	uc := usecase.NewQueryUseCase(fakeEmbedder{}, fakeStore{}, &fakeLLM{}, usecase.DefaultPromptBuilder, 4)

	_, err := uc.Execute(context.Background(), usecase.QueryInput{Question: "anything"})
	if err != domain.ErrNoRelevantDocuments {
		t.Errorf("expected ErrNoRelevantDocuments, got %v", err)
	}
}
