package port

// TextSplitter is the port for chunking source text into retrievable units.
type TextSplitter interface {
	// Split divides text into overlapping chunks suitable for embedding.
	Split(text string) ([]string, error)
}
