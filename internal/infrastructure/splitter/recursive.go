package splitter

import (
	"github.com/tmc/langchaingo/textsplitter"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// Config controls chunking behaviour.
type Config struct {
	ChunkSize    int
	ChunkOverlap int
}

// recursiveSplitter adapts langchaingo's recursive character splitter to the
// domain port.TextSplitter.
type recursiveSplitter struct {
	splitter textsplitter.RecursiveCharacter
}

// compile-time check that recursiveSplitter satisfies the port.
var _ port.TextSplitter = (*recursiveSplitter)(nil)

// NewRecursive builds a port.TextSplitter, applying sensible defaults when
// values are unset.
func NewRecursive(cfg Config) port.TextSplitter {
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = 1000
	}
	if cfg.ChunkOverlap < 0 {
		cfg.ChunkOverlap = 0
	}

	s := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(cfg.ChunkSize),
		textsplitter.WithChunkOverlap(cfg.ChunkOverlap),
	)
	return &recursiveSplitter{splitter: s}
}

func (r *recursiveSplitter) Split(text string) ([]string, error) {
	return r.splitter.SplitText(text)
}
