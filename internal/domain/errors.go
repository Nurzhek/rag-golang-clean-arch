package domain

import "errors"

// Domain-level sentinel errors. Outer layers match on these with errors.Is to
// translate failures into transport responses without importing infrastructure.
var (
	// ErrNoRelevantDocuments is returned when a query retrieves no context.
	ErrNoRelevantDocuments = errors.New("no relevant documents found")
	// ErrEmptyContent is returned when ingest is called with blank content.
	ErrEmptyContent = errors.New("document content must not be empty")
	// ErrEmptyQuestion is returned when a query is called with a blank question.
	ErrEmptyQuestion = errors.New("question must not be empty")
)
