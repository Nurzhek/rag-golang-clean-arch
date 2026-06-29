package jobstore

import (
	"context"
	"sync"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/entity"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// Memory is a thread-safe, in-process implementation of port.JobRepository.
// Job state is ephemeral and lost on restart — adequate for the in-memory
// vector store it accompanies; swap for Redis/Postgres for durability.
type Memory struct {
	mu   sync.RWMutex
	jobs map[string]entity.Job
}

// compile-time check that Memory satisfies the port.
var _ port.JobRepository = (*Memory)(nil)

// NewMemory creates an empty in-memory job repository.
func NewMemory() *Memory {
	return &Memory{jobs: make(map[string]entity.Job)}
}

// Save upserts a job by its ID.
func (m *Memory) Save(_ context.Context, job entity.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = job
	return nil
}

// Get returns a job by ID, or domain.ErrJobNotFound.
func (m *Memory) Get(_ context.Context, id string) (entity.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[id]
	if !ok {
		return entity.Job{}, domain.ErrJobNotFound
	}
	return job, nil
}
