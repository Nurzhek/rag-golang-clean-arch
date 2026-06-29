package dto

// --- Ingestion ---

// IngestRequest is the JSON body for POST /api/v1/documents.
type IngestRequest struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// JobAcceptedResponse is returned when ingestion is queued. Poll PollURL for status.
type JobAcceptedResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	PollURL string `json:"poll_url"`
}

// --- Jobs ---

// JobResponse is the polling payload for GET /api/v1/jobs/{id}.
type JobResponse struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	TotalChunks     int    `json:"total_chunks"`
	ProcessedChunks int    `json:"processed_chunks"`
	DocumentID      string `json:"document_id,omitempty"`
	Error           string `json:"error,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// --- Documents ---

// DocumentSummary describes one ingested document in a listing.
type DocumentSummary struct {
	ID         string            `json:"id"`
	ChunkCount int               `json:"chunk_count"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// ListDocumentsResponse is returned by GET /api/v1/documents.
type ListDocumentsResponse struct {
	Count     int               `json:"count"`
	Documents []DocumentSummary `json:"documents"`
}

// DeleteResponse is returned by DELETE /api/v1/documents/{id}.
type DeleteResponse struct {
	DocumentID    string `json:"document_id"`
	DeletedChunks int    `json:"deleted_chunks"`
}

// --- Query ---

// QueryRequest is the JSON body for POST /api/v1/query.
type QueryRequest struct {
	Question string `json:"question"`
	TopK     int    `json:"top_k,omitempty"`
}

// Source describes a retrieved chunk that grounded the answer.
type Source struct {
	ID         string            `json:"id"`
	DocumentID string            `json:"document_id,omitempty"`
	Content    string            `json:"content"`
	Score      float32           `json:"score"`
	Metadata   map[string]string `json:"metadata,omitempty"`
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
