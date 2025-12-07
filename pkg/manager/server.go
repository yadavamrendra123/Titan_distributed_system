package manager

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/google/uuid"
	"titan/pkg/logger"
	"titan/pkg/models"
	pb "titan/pkg/proto"
)

// Server implements the Manager RPC service
type Server struct {
	store     *Store
	scheduler *Scheduler
}

// NewServer creates a new Manager server
func NewServer() *Server {
	store := NewStore()
	return &Server{
		store:     store,
		scheduler: NewScheduler(store),
	}
}

// Start begins the manager's background tasks
func (s *Server) Start() {
	go s.scheduler.Run()
	logger.Info("Manager server started")
}

// RegisterRPC registers the server with the net/rpc handler
func (s *Server) RegisterRPC(server *rpc.Server) {
	server.RegisterName("ManagerService", s)
	server.RegisterName("WorkerService", s)
}

// SubmitJob handles job submission from clients
// Signature must be: func (t *T) MethodName(argType T1, replyType *T2) error
func (s *Server) SubmitJob(req pb.JobRequest, resp *pb.JobResponse) error {
	jobID := uuid.New().String()
	
	job := &models.Job{
		ID:        jobID,
		Command:   req.Command,
		Env:       req.Env,
		Status:    models.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	s.store.AddJob(job)
	
	logger.Info("Job submitted", "job_id", jobID, "command", req.Command)
	
	*resp = pb.JobResponse{
		JobId:  jobID,
		Status: string(job.Status),
	}
	return nil
}

// GetJobStatus returns the current status of a job
func (s *Server) GetJobStatus(req pb.JobStatusRequest, resp *pb.JobStatusResponse) error {
	job, ok := s.store.GetJob(req.JobId)
	if !ok {
		return fmt.Errorf("job not found: %s", req.JobId)
	}
	
	*resp = pb.JobStatusResponse{
		JobId:    job.ID,
		Status:   string(job.Status),
		WorkerId: job.WorkerID,
		Output:   job.Output,
		ExitCode: job.ExitCode,
	}
	return nil
}

// ListJobs returns all jobs in the cluster
func (s *Server) ListJobs(req pb.ListJobsRequest, resp *pb.ListJobsResponse) error {
	jobs := s.store.GetAllJobs()
	
	resp.Jobs = make([]pb.JobStatusResponse, len(jobs))
	
	for i, job := range jobs {
		resp.Jobs[i] = pb.JobStatusResponse{
			JobId:    job.ID,
			Status:   string(job.Status),
			WorkerId: job.WorkerID,
			Output:   job.Output,
			ExitCode: job.ExitCode,
		}
	}
	return nil
}

// RegisterWorker handles worker registration
func (s *Server) RegisterWorker(req pb.WorkerInfo, resp *pb.RegistrationResponse) error {
	worker := &models.Worker{
		ID:           req.WorkerId,
		Address:      req.Address,
		TotalCPU:     req.Capacity.TotalCpuMillicores,
		TotalMemory:  req.Capacity.TotalMemoryMb,
		Status:       models.WorkerStatusHealthy,
		LastHeartbeat: time.Now(),
		RegisteredAt: time.Now(),
	}
	
	s.store.RegisterWorker(worker)
	
	logger.Info("Worker registered", "worker_id", req.WorkerId, "address", req.Address)
	
	*resp = pb.RegistrationResponse{
		Accepted: true,
		Message:  "Worker registered successfully",
	}
	return nil
}

// Heartbeat handles worker heartbeats
func (s *Server) Heartbeat(req pb.HeartbeatRequest, resp *pb.HeartbeatResponse) error {
	usage := &models.Worker{
		UsedCPU:    req.CurrentUsage.UsedCpuMillicores,
		UsedMemory: req.CurrentUsage.UsedMemoryMb,
	}
	
	s.store.UpdateWorkerHeartbeat(req.WorkerId, usage)
	
	*resp = pb.HeartbeatResponse{
		Acknowledged: true,
	}
	return nil
}

// ReportTaskStatus handles task status updates from workers
func (s *Server) ReportTaskStatus(req pb.TaskStatusUpdate, resp *pb.Ack) error {
	// Find the job associated with this task
	job, ok := s.store.GetJob(req.TaskId)
	if !ok {
		return fmt.Errorf("job not found: %s", req.TaskId)
	}
	
	// Update job status based on task status
	job.Status = models.JobStatus(req.Status)
	job.Output = req.Output
	job.ExitCode = req.ExitCode
	
	s.store.UpdateJob(job)
	
	logger.Info("Task status updated", "task_id", req.TaskId, "status", req.Status)
	
	*resp = pb.Ack{Ok: true}
	return nil
}

// StartTask is called by scheduler to start a task on a worker
// Note: This is usually called client->server, but here it's just a placeholder
// The real StartTask happens in the Worker service
func (s *Server) StartTask(req pb.TaskRequest, resp *pb.TaskResponse) error {
	*resp = pb.TaskResponse{
		Accepted: true,
		Message:  "Task started",
	}
	return nil
}

// StopTask stops a running task
func (s *Server) StopTask(req pb.StopTaskRequest, resp *pb.StopTaskResponse) error {
	*resp = pb.StopTaskResponse{
		Stopped: true,
	}
	return nil
}

// GetRPCServer creates and returns a net/rpc server
func (s *Server) GetRPCServer() *rpc.Server {
	server := rpc.NewServer()
	s.RegisterRPC(server)
	return server
}
