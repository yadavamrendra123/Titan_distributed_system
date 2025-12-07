package proto

// Plain Go structs for net/rpc communication
// No protobuf dependency needed

type JobRequest struct {
	Command   string
	Env       map[string]string
	Resources ResourceRequirements
}

type ResourceRequirements struct {
	CpuMillicores int32
	MemoryMb      int64
}

type JobResponse struct {
	JobId  string
	Status string
}

type JobStatusRequest struct {
	JobId string
}

type JobStatusResponse struct {
	JobId    string
	Status   string
	WorkerId string
	Output   string
	ExitCode int32
}

type ListJobsRequest struct {
	// Empty
}

type ListJobsResponse struct {
	Jobs []JobStatusResponse
}

type WorkerInfo struct {
	WorkerId string
	Address  string
	Capacity ResourceCapacity
}

type ResourceCapacity struct {
	TotalCpuMillicores int32
	TotalMemoryMb      int64
}

type RegistrationResponse struct {
	Accepted bool
	Message  string
}

type HeartbeatRequest struct {
	WorkerId     string
	Timestamp    int64
	CurrentUsage ResourceUsage
}

type ResourceUsage struct {
	UsedCpuMillicores int32
	UsedMemoryMb      int64
}

type HeartbeatResponse struct {
	Acknowledged bool
}

type TaskRequest struct {
	TaskId  string
	JobId   string
	Command string
	Env     map[string]string
}

type TaskResponse struct {
	Accepted bool
	Message  string
}

type StopTaskRequest struct {
	TaskId string
}

type StopTaskResponse struct{
	Stopped bool
}

type TaskStatusUpdate struct {
	TaskId   string
	Status   string
	Output   string
	ExitCode int32
}

type Ack struct {
	Ok bool
}
