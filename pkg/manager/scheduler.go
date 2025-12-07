package manager

import (
	"fmt"
	"net/rpc"
	"time"

	"titan/pkg/logger"
	"titan/pkg/models"
	pb "titan/pkg/proto"
)

// Scheduler is responsible for assigning jobs to workers
type Scheduler struct {
	store        *Store
	stopChan     chan struct{}
	workerClients map[string]*rpc.Client
}

// NewScheduler creates a new scheduler
func NewScheduler(store *Store) *Scheduler {
	return &Scheduler{
		store:        store,
		stopChan:     make(chan struct{}),
		workerClients: make(map[string]*rpc.Client),
	}
}

// Run starts the scheduling loop
func (s *Scheduler) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	logger.Info("Scheduler started")
	
	for {
		select {
		case <-ticker.C:
			s.schedule()
		case <-s.stopChan:
			logger.Info("Scheduler stopped")
			return
		}
	}
}

// Stop halts the scheduler
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// schedule attempts to assign pending jobs to available workers
func (s *Scheduler) schedule() {
	pendingJobs := s.store.GetPendingJobs()
	if len(pendingJobs) == 0 {
		return
	}
	
	healthyWorkers := s.store.GetHealthyWorkers()
	if len(healthyWorkers) == 0 {
		logger.Warn("No healthy workers available for scheduling")
		return
	}
	
	// Round-robin scheduling
	workerIndex := 0
	for _, job := range pendingJobs {
		worker := healthyWorkers[workerIndex%len(healthyWorkers)]
		
		// Assign job to worker
		if err := s.assignJobToWorker(job, worker); err != nil {
			logger.Error("Failed to assign job to worker", 
				"job_id", job.ID, 
				"worker_id", worker.ID, 
				"error", err)
			continue
		}
		
		// Update job status
		job.Status = models.JobStatusScheduled
		job.WorkerID = worker.ID
		s.store.UpdateJob(job)
		
		logger.Info("Job scheduled", 
			"job_id", job.ID, 
			"worker_id", worker.ID)
		
		workerIndex++
	}
}

// assignJobToWorker sends a StartTask RPC to the worker
func (s *Scheduler) assignJobToWorker(job *models.Job, worker *models.Worker) error {
	client, err := s.getWorkerClient(worker)
	if err != nil {
		return fmt.Errorf("failed to get worker client: %w", err)
	}
	
	taskID := job.ID // Use job ID as task ID for simplicity
	req := pb.TaskRequest{
		TaskId:  taskID,
		JobId:   job.ID,
		Command: job.Command,
		Env:     job.Env,
	}
	
	var resp pb.TaskResponse
	err = client.Call("WorkerService.StartTask", req, &resp)
	if err != nil {
		// If call fails, remove client to force reconnect next time
		delete(s.workerClients, worker.ID)
		return fmt.Errorf("failed to start task on worker: %w", err)
	}
	
	if !resp.Accepted {
		return fmt.Errorf("worker rejected task: %s", resp.Message)
	}
	
	return nil
}

// getWorkerClient returns or creates a net/rpc client for the worker
func (s *Scheduler) getWorkerClient(worker *models.Worker) (*rpc.Client, error) {
	if client, ok := s.workerClients[worker.ID]; ok {
		return client, nil
	}
	
	client, err := rpc.Dial("tcp", worker.Address)
	if err != nil {
		return nil, err
	}
	
	s.workerClients[worker.ID] = client
	
	return client, nil
}
