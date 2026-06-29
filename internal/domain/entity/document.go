package entity

// Document is the core unit of knowledge in the system. It is produced by
// splitting source material into a retrievable chunk and is what the vector
// store persists and returns during retrieval.
type Document struct {
	// ID uniquely identifies the chunk within the store.
	ID string
	// SourceID groups every chunk that originated from the same ingested
	// document, enabling listing and deletion at the document level.
	SourceID string
	// Content is the raw text of the chunk.
	Content string
	// Metadata carries arbitrary key/value annotations (e.g. source, title).
	Metadata map[string]string
}

// ScoredDocument is a Document paired with its similarity score relative to a
// query. Scores are produced by the vector store during retrieval; a higher
// score means greater similarity.
type ScoredDocument struct {
	Document
	Score float32
}

// DocumentSummary is an aggregate view of one ingested document: its source ID,
// how many chunks it produced, and the metadata shared by those chunks.
type DocumentSummary struct {
	ID         string
	ChunkCount int
	Metadata   map[string]string
}
