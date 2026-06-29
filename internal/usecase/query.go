package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

const fallbackTopK = 4

// PromptBuilder turns a question plus retrieved context into a final prompt. It
// is injected so the prompting strategy can change without touching the use case.
type PromptBuilder func(question string, sources []entity.ScoredDocument) string

// queryInteractor implements QueryUseCase. It embeds the question, retrieves the
// most similar chunks, builds a grounded prompt, and asks the LLM to answer.
type queryInteractor struct {
	embedder    port.Embedder
	store       port.VectorStore
	llm         port.LLM
	buildPrompt PromptBuilder
	defaultTopK int
}

// NewQueryUseCase wires the query dependencies. If topK <= 0 a sane default is used.
func NewQueryUseCase(embedder port.Embedder, store port.VectorStore, llm port.LLM, pb PromptBuilder, topK int) QueryUseCase {
	if topK <= 0 {
		topK = fallbackTopK
	}
	return &queryInteractor{embedder: embedder, store: store, llm: llm, buildPrompt: pb, defaultTopK: topK}
}

func (uc *queryInteractor) Execute(ctx context.Context, in QueryInput) (QueryOutput, error) {
	question := strings.TrimSpace(in.Question)
	if question == "" {
		return QueryOutput{}, domain.ErrEmptyQuestion
	}

	topK := in.TopK
	if topK <= 0 {
		topK = uc.defaultTopK
	}

	queryVec, err := uc.embedder.EmbedQuery(ctx, question)
	if err != nil {
		return QueryOutput{}, fmt.Errorf("embed question: %w", err)
	}

	sources, err := uc.store.Search(ctx, queryVec, topK)
	if err != nil {
		return QueryOutput{}, fmt.Errorf("search store: %w", err)
	}
	if len(sources) == 0 {
		return QueryOutput{}, domain.ErrNoRelevantDocuments
	}

	answer, err := uc.llm.Generate(ctx, uc.buildPrompt(question, sources))
	if err != nil {
		return QueryOutput{}, fmt.Errorf("generate answer: %w", err)
	}

	return QueryOutput{Answer: strings.TrimSpace(answer), Sources: sources}, nil
}
