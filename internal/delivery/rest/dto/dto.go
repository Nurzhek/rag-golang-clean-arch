package dto

// IngestRequest is the JSON body for POST /api/v1/documents.
type IngestRequest struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// IngestResponse is returned after a successful ingest.
type IngestResponse struct {
	ChunksCreated int `json:"chunks_created"`
}

// QueryRequest is the JSON body for POST /api/v1/query.
type QueryRequest struct {
	Question string `json:"question"`
	TopK     int    `json:"top_k,omitempty"`
}

// Source describes a retrieved chunk that grounded the answer.
type Source struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// QueryResponse is returned for a successful query.
type QueryResponse struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
}

// ErrorResponse is the standard error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}
