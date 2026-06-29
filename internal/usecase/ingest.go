package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

const defaultEmbedBatchSize = 64

// IDGenerator produces unique IDs for jobs, source documents, and chunks. It is
// injected so the use case stays deterministic and easy to test.
type IDGenerator func() string

// asyncIngestInteractor implements IngestUseCase. Submit validates the input,
// records a queued job, and runs the ingestion pipeline (split -> embed -> store)
// in the background, embedding in batches so progress is observable via polling.
type asyncIngestInteractor struct {
	splitter  port.TextSplitter
	embedder  port.Embedder
	store     port.VectorStore
	jobs      port.JobRepository
	newID     IDGenerator
	batchSize int
	log       *slog.Logger
}

// NewIngestUseCase wires the asynchronous ingestion dependencies.
func NewIngestUseCase(
	splitter port.TextSplitter,
	embedder port.Embedder,
	store port.VectorStore,
	jobs port.JobRepository,
	newID IDGenerator,
	batchSize int,
	log *slog.Logger,
) IngestUseCase {
	if batchSize <= 0 {
		batchSize = defaultEmbedBatchSize
	}
	return &asyncIngestInteractor{
		splitter:  splitter,
		embedder:  embedder,
		store:     store,
		jobs:      jobs,
		newID:     newID,
		batchSize: batchSize,
		log:       log,
	}
}

// Submit records a queued job and starts background processing, returning the
// job ID immediately so the HTTP handler can respond without blocking.
func (uc *asyncIngestInteractor) Submit(ctx context.Context, in IngestInput) (string, error) {
	in.Content = strings.TrimSpace(in.Content)
	if in.Content == "" {
		return "", domain.ErrEmptyContent
	}

	now := time.Now().UTC()
	job := entity.Job{
		ID:        uc.newID(),
		Status:    entity.JobStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.jobs.Save(ctx, job); err != nil {
		return "", fmt.Errorf("create job: %w", err)
	}

	// Run in the background with a fresh context: the request context is
	// cancelled the moment the HTTP handler returns.
	go uc.process(job.ID, in)

	return job.ID, nil
}

func (uc *asyncIngestInteractor) process(jobID string, in IngestInput) {
	ctx := context.Background()

	job, err := uc.jobs.Get(ctx, jobID)
	if err != nil {
		uc.log.Error("ingest worker: job lookup failed", "job_id", jobID, "error", err)
		return
	}

	fail := func(cause error) {
		job.Status = entity.JobStatusFailed
		job.Error = cause.Error()
		job.UpdatedAt = time.Now().UTC()
		_ = uc.jobs.Save(ctx, job)
		uc.log.Error("ingest job failed", "job_id", jobID, "error", cause)
	}

	job.Status = entity.JobStatusRunning
	job.UpdatedAt = time.Now().UTC()
	_ = uc.jobs.Save(ctx, job)

	chunks, err := uc.splitter.Split(in.Content)
	if err != nil {
		fail(fmt.Errorf("split content: %w", err))
		return
	}
	if len(chunks) == 0 {
		fail(domain.ErrEmptyContent)
		return
	}

	sourceID := uc.newID()
	metadata := cloneMetadata(in.Metadata)

	job.DocumentID = sourceID
	job.TotalChunks = len(chunks)
	job.UpdatedAt = time.Now().UTC()
	_ = uc.jobs.Save(ctx, job)

	// Embed and store in batches so a large file reports incremental progress.
	for start := 0; start < len(chunks); start += uc.batchSize {
		end := min(start+uc.batchSize, len(chunks))
		batch := chunks[start:end]

		vectors, err := uc.embedder.EmbedDocuments(ctx, batch)
		if err != nil {
			fail(fmt.Errorf("embed chunks: %w", err))
			return
		}

		docs := make([]entity.Document, len(batch))
		for i, chunk := range batch {
			docs[i] = entity.Document{
				ID:       uc.newID(),
				SourceID: sourceID,
				Content:  chunk,
				Metadata: metadata,
			}
		}
		if err := uc.store.Add(ctx, docs, vectors); err != nil {
			fail(fmt.Errorf("store documents: %w", err))
			return
		}

		job.ProcessedChunks = end
		job.UpdatedAt = time.Now().UTC()
		_ = uc.jobs.Save(ctx, job)
	}

	job.Status = entity.JobStatusCompleted
	job.UpdatedAt = time.Now().UTC()
	_ = uc.jobs.Save(ctx, job)
	uc.log.Info("ingest job completed", "job_id", jobID, "document_id", sourceID, "chunks", job.TotalChunks)
}

func cloneMetadata(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
