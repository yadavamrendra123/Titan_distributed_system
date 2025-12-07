package manager

import (
	"sync"
	"time"

	"titan/pkg/models"
)

// Store manages all cluster state in-memory
type Store struct {
	mu      sync.RWMutex
	jobs    map[string]*models.Job
	workers map[string]*models.Worker
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	return &Store{
		jobs:    make(map[string]*models.Job),
		workers: make(map[string]*models.Worker),
	}
}

// AddJob stores a new job
func (s *Store) AddJob(job *models.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

// GetJob retrieves a job by ID
func (s *Store) GetJob(id string) (*models.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	return job, ok
}

// UpdateJob updates an existing job
func (s *Store) UpdateJob(job *models.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job.UpdatedAt = time.Now()
	s.jobs[job.ID] = job
}

// GetAllJobs returns all jobs
func (s *Store) GetAllJobs() []*models.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs := make([]*models.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// GetPendingJobs returns jobs that need to be scheduled
func (s *Store) GetPendingJobs() []*models.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pending := make([]*models.Job, 0)
	for _, job := range s.jobs {
		if job.Status == models.JobStatusPending {
			pending = append(pending, job)
		}
	}
	return pending
}

// RegisterWorker adds or updates a worker
func (s *Store) RegisterWorker(worker *models.Worker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workers[worker.ID] = worker
}

// GetWorker retrieves a worker by ID
func (s *Store) GetWorker(id string) (*models.Worker, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	worker, ok := s.workers[id]
	return worker, ok
}

// UpdateWorkerHeartbeat updates the last heartbeat time for a worker
func (s *Store) UpdateWorkerHeartbeat(workerID string, usage *models.Worker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if worker, ok := s.workers[workerID]; ok {
		worker.LastHeartbeat = time.Now()
		if usage != nil {
			worker.UsedCPU = usage.UsedCPU
			worker.UsedMemory = usage.UsedMemory
		}
		worker.Status = models.WorkerStatusHealthy
	}
}

// GetHealthyWorkers returns all workers that are healthy
func (s *Store) GetHealthyWorkers() []*models.Worker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	healthy := make([]*models.Worker, 0)
	now := time.Now()
	heartbeatTimeout := 30 * time.Second
	
	for _, worker := range s.workers {
		if now.Sub(worker.LastHeartbeat) < heartbeatTimeout {
			healthy = append(healthy, worker)
		}
	}
	return healthy
}

// GetAllWorkers returns all registered workers
func (s *Store) GetAllWorkers() []*models.Worker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	workers := make([]*models.Worker, 0, len(s.workers))
	for _, worker := range s.workers {
		workers = append(workers, worker)
	}
	return workers
}
