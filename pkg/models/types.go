package models

import "time"

// JobStatus represents the lifecycle state of a job
type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusScheduled JobStatus = "SCHEDULED"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

// Job represents a unit of work to be executed
type Job struct {
	ID        string
	Command   string
	Env       map[string]string
	Status    JobStatus
	WorkerID  string // Assigned worker
	Output    string
	ExitCode  int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorkerStatus represents the health state of a worker
type WorkerStatus string

const (
	WorkerStatusHealthy   WorkerStatus = "HEALTHY"
	WorkerStatusUnhealthy WorkerStatus = "UNHEALTHY"
)

// Worker represents a compute node in the cluster
type Worker struct {
	ID               string
	Address          string
	TotalCPU         int32
	TotalMemory      int64
	UsedCPU          int32
	UsedMemory       int64
	Status           WorkerStatus
	LastHeartbeat    time.Time
	RegisteredAt     time.Time
}

// Task represents a running instance of a job on a worker
type Task struct {
	ID       string
	JobID    string
	WorkerID string
	Command  string
	Env      map[string]string
	Status   JobStatus
	Output   string
	ExitCode int32
}
