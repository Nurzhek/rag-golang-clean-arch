package entity

import "time"

// JobStatus is the lifecycle state of an asynchronous ingestion job.
type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// Job tracks the progress of a background ingestion request so clients can poll
// for completion instead of blocking on the HTTP call.
type Job struct {
	ID              string
	Status          JobStatus
	TotalChunks     int
	ProcessedChunks int
	// DocumentID is the source ID of the ingested document, set once known.
	DocumentID string
	// Error holds the failure reason when Status is JobStatusFailed.
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
