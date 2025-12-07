package worker

import (
	"fmt"
	"net/rpc"

	"titan/pkg/logger"
	pb "titan/pkg/proto"
)

// Server implements the Worker RPC service
type Server struct {
	executor     *Executor
	heartbeater  *Heartbeater
	workerID     string
	address      string
	managerAddr  string
}

// NewServer creates a new Worker server
func NewServer(workerID, address, managerAddr string) (*Server, error) {
	executor, err := NewExecutor(managerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}
	
	return &Server{
		executor:    executor,
		workerID:    workerID,
		address:     address,
		managerAddr: managerAddr,
	}, nil
}

// Start initializes the worker and begins background tasks
func (s *Server) Start() error {
	// Register with manager
	if err := s.register(); err != nil {
		return fmt.Errorf("failed to register with manager: %w", err)
	}
	
	// Start heartbeat
	s.heartbeater = NewHeartbeater(s.workerID, s.executor.managerClient, 10*1e9) // 10 seconds
	go s.heartbeater.Start()
	
	logger.Info("Worker server started", "worker_id", s.workerID, "address", s.address)
	return nil
}

// register registers the worker with the manager
func (s *Server) register() error {
	req := pb.WorkerInfo{
		WorkerId: s.workerID,
		Address:  s.address,
		Capacity: pb.ResourceCapacity{
			TotalCpuMillicores: 4000, // 4 cores
			TotalMemoryMb:      8192, // 8 GB
		},
	}
	
	var resp pb.RegistrationResponse
	err := s.executor.managerClient.Call("ManagerService.RegisterWorker", req, &resp)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}
	
	if !resp.Accepted {
		return fmt.Errorf("registration rejected: %s", resp.Message)
	}
	
	logger.Info("Worker registered successfully", "worker_id", s.workerID)
	return nil
}

// StartTask handles task execution requests from the manager
func (s *Server) StartTask(req pb.TaskRequest, resp *pb.TaskResponse) error {
	logger.Info("Received task", "task_id", req.TaskId, "job_id", req.JobId)
	
	err := s.executor.StartTask(req.TaskId, req.JobId, req.Command, req.Env)
	if err != nil {
		logger.Error("Failed to start task", "task_id", req.TaskId, "error", err)
		*resp = pb.TaskResponse{
			Accepted: false,
			Message:  err.Error(),
		}
		return nil
	}
	
	*resp = pb.TaskResponse{
		Accepted: true,
		Message:  "Task accepted",
	}
	return nil
}

// StopTask handles task termination requests
func (s *Server) StopTask(req pb.StopTaskRequest, resp *pb.StopTaskResponse) error {
	logger.Info("Stopping task", "task_id", req.TaskId)
	
	err := s.executor.StopTask(req.TaskId)
	if err != nil {
		logger.Error("Failed to stop task", "task_id", req.TaskId, "error", err)
		*resp = pb.StopTaskResponse{
			Stopped: false,
		}
		return nil
	}
	
	*resp = pb.StopTaskResponse{
		Stopped: true,
	}
	return nil
}

// RegisterRPC registers the server with the net/rpc handler
func (s *Server) RegisterRPC(server *rpc.Server) {
	server.RegisterName("WorkerService", s)
}

// GetRPCServer creates and returns a net/rpc server
func (s *Server) GetRPCServer() *rpc.Server {
	server := rpc.NewServer()
	s.RegisterRPC(server)
	return server
}
