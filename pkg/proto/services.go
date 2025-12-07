// Code generated - gRPC service stubs for Windows compatibility
// This file provides minimal interfaces to allow compilation without protoc-gen-go-grpc

package proto

import (
	"context"
	"google.golang.org/grpc"
)

// ManagerService defines the manager API
type ManagerServiceServer interface {
	SubmitJob(context.Context, *JobRequest) (*JobResponse, error)
	GetJobStatus(context.Context, *JobStatusRequest) (*JobStatusResponse, error)
	ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error)
}

type ManagerServiceClient interface {
	SubmitJob(ctx context.Context, in *JobRequest, opts ...grpc.CallOption) (*JobResponse, error)
	GetJobStatus(ctx context.Context, in *JobStatusRequest, opts ...grpc.CallOption) (*JobStatusResponse, error)
	ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
	RegisterWorker(ctx context.Context, in *WorkerInfo, opts ...grpc.CallOption) (*RegistrationResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	ReportTaskStatus(ctx context.Context, in *TaskStatusUpdate, ops ...grpc.CallOption) (*Ack, error)
}

// WorkerService defines the worker API
type WorkerServiceServer interface {
	RegisterWorker(context.Context, *WorkerInfo) (*RegistrationResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	StartTask(context.Context, *TaskRequest) (*TaskResponse, error)
	StopTask(context.Context, *StopTaskRequest) (*StopTaskResponse, error)
	ReportTaskStatus(context.Context, *TaskStatusUpdate) (*Ack, error)
}

type WorkerServiceClient interface {
	RegisterWorker(ctx context.Context, in *WorkerInfo, opts ...grpc.CallOption) (*RegistrationResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	StartTask(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskResponse, error)
	StopTask(ctx context.Context, in *StopTaskRequest, opts ...grpc.CallOption) (*StopTaskResponse, error)
	ReportTaskStatus(ctx context.Context, in *TaskStatusUpdate, opts ...grpc.CallOption) (*Ack, error)
}

// Unimplemented server stubs
type UnimplementedManagerServiceServer struct{}
type UnimplementedWorkerServiceServer struct{}

func (UnimplementedManagerServiceServer) SubmitJob(context.Context, *JobRequest) (*JobResponse, error) {
	return nil, nil
}

func (UnimplementedManagerServiceServer) GetJobStatus(context.Context, *JobStatusRequest) (*JobStatusResponse, error) {
	return nil, nil
}

func (UnimplementedManagerServiceServer) ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error) {
	return nil, nil
}

func (UnimplementedWorkerServiceServer) RegisterWorker(context.Context, *WorkerInfo) (*RegistrationResponse, error) {
	return nil, nil
}

func (UnimplementedWorkerServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, nil
}

func (UnimplementedWorkerServiceServer) StartTask(context.Context, *TaskRequest) (*TaskResponse, error) {
	return nil, nil
}

func (UnimplementedWorkerServiceServer) StopTask(context.Context, *StopTaskRequest) (*StopTaskResponse, error) {
	return nil, nil
}

func (UnimplementedWorkerServiceServer) ReportTaskStatus(context.Context, *TaskStatusUpdate) (*Ack, error) {
	return nil, nil
}

// Registration functions
func RegisterManagerServiceServer(s *grpc.Server, srv ManagerServiceServer) {
	// Stub for compilation
}

func RegisterWorkerServiceServer(s *grpc.Server, srv WorkerServiceServer) {
	// Stub for compilation
}

// Client constructors
func NewManagerServiceClient(cc grpc.ClientConnInterface) ManagerServiceClient {
	return &managerServiceClient{cc}
}

func NewWorkerServiceClient(cc grpc.ClientConnInterface) WorkerServiceClient {
	return &workerServiceClient{cc}
}

type managerServiceClient struct {
	cc grpc.ClientConnInterface
}

func (c *managerServiceClient) SubmitJob(ctx context.Context, in *JobRequest, opts ...grpc.CallOption) (*JobResponse, error) {
	out := new(JobResponse)
	err := c.cc.Invoke(ctx, "/titan.ManagerService/SubmitJob", in, out, opts...)
	return out, err
}

func (c *managerServiceClient) GetJobStatus(ctx context.Context, in *JobStatusRequest, opts ...grpc.CallOption) (*JobStatusResponse, error) {
	out := new(JobStatusResponse)
	err := c.cc.Invoke(ctx, "/titan.ManagerService/GetJobStatus", in, out, opts...)
	return out, err
}

func (c *managerServiceClient) ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse)
	err := c.cc.Invoke(ctx, "/titan.ManagerService/ListJobs", in, out, opts...)
	return out, err
}

func (c *managerServiceClient) RegisterWorker(ctx context.Context, in *WorkerInfo, opts ...grpc.CallOption) (*RegistrationResponse, error) {
	out := new(RegistrationResponse)
	err := c.cc.Invoke(ctx, "/titan.ManagerService/RegisterWorker", in, out, opts...)
	return out, err
}

func (c *managerServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/titan.ManagerService/Heartbeat", in, out, opts...)
	return out, err
}

func (c *managerServiceClient) ReportTaskStatus(ctx context.Context, in *TaskStatusUpdate, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/titan.ManagerService/ReportTaskStatus", in, out, opts...)
	return out, err
}

type workerServiceClient struct {
	cc grpc.ClientConnInterface
}

func (c *workerServiceClient) RegisterWorker(ctx context.Context, in *WorkerInfo, opts ...grpc.CallOption) (*RegistrationResponse, error) {
	out := new(RegistrationResponse)
	err := c.cc.Invoke(ctx, "/titan.WorkerService/RegisterWorker", in, out, opts...)
	return out, err
}

func (c *workerServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/titan.WorkerService/Heartbeat", in, out, opts...)
	return out, err
}

func (c *workerServiceClient) StartTask(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskResponse, error) {
	out := new(TaskResponse)
	err := c.cc.Invoke(ctx, "/titan.WorkerService/StartTask", in, out, opts...)
	return out, err
}

func (c *workerServiceClient) StopTask(ctx context.Context, in *StopTaskRequest, opts ...grpc.CallOption) (*StopTaskResponse, error) {
	out := new(StopTaskResponse)
	err := c.cc.Invoke(ctx, "/titan.WorkerService/StopTask", in, out, opts...)
	return out, err
}

func (c *workerServiceClient) ReportTaskStatus(ctx context.Context, in *TaskStatusUpdate, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/titan.WorkerService/ReportTaskStatus", in, out, opts...)
	return out, err
}
