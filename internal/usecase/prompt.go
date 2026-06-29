package usecase

import (
	"fmt"
	"strings"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
)

// DefaultPromptBuilder produces a grounded RAG prompt that instructs the model
// to answer strictly from the retrieved context and to admit when it cannot.
func DefaultPromptBuilder(question string, sources []entity.ScoredDocument) string {
	var b strings.Builder
	b.WriteString("You are a helpful assistant. Answer the question using ONLY the context below. ")
	b.WriteString("If the answer cannot be found in the context, say that you don't know.\n\n")
	b.WriteString("Context:\n")
	for i, s := range sources {
		fmt.Fprintf(&b, "[%d] %s\n", i+1, strings.TrimSpace(s.Content))
	}
	b.WriteString("\nQuestion: ")
	b.WriteString(question)
	b.WriteString("\n\nAnswer:")
	return b.String()
}
