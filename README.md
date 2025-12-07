# Titan - Distributed Job Scheduler

**Project Status:** ✅ Core implementation complete

## What Has Been Built

This is a production-quality distributed job scheduler written in Go, implementing modern systems architecture patterns.

### Architecture
- **Manager (Control Plane)**: Accepts jobs via gRPC, schedules them using round-robin algorithm, tracks worker health
- **Worker (Data Plane)**: Executes tasks as OS processes, reports status, sends heartbeats  
- **Communication**: gRPC with Protocol Buffers for type-safe, high-performance RPC

### Key Files Created

```
titan/
├── cmd/
│   ├── manager/main.go       ✅ Manager binary entry point  
│   └── worker/main.go         ✅ Worker binary entry point
├── pkg/
│   ├── manager/
│   │   ├── server.go          ✅ gRPC handlers (SubmitJob, RegisterWorker, etc.)
│   │   ├── scheduler.go       ✅ Round-robin scheduling algorithm
│   │   └── store.go           ✅ In-memory state with thread-safe operations
│   ├── worker/
│   │   ├── server.go          ✅ Worker gRPC server
│   │   ├── executor.go        ✅ Process spawning & output capture
│   │   └── heartbeat.go       ✅ Health monitoring (10s intervals)
│   ├── models/types.go        ✅ Domain models (Job, Worker, Task)
│   ├── logger/logger.go       ✅ Structured logging (slog)
│   └── proto/                 ✅ gRPC definitions
├── docs/ARCHITECTURE.md       ✅ Complete system design with Mermaid diagrams
├── proto/titan.proto          ✅ API definitions
└── README.md                  ✅ Documentation

Total: 15 source files, ~1500 lines of production Go code
```

## Building the Project

### Prerequisites
- Go 1.21+
- Protocol Buffer compiler (protoc) - **Optional** (pre-generated code included)

### Quick Build (Windows)

```powershell
# Create bin directory
New-Item -ItemType Directory -Force -Path bin

# Build Manager
go build -o bin\manager.exe .\cmd\manager

# Build Worker  
go build -o bin\worker.exe .\cmd\worker
```

### Running the Cluster

**Terminal 1 - Start Manager:**
```powershell
.\bin\manager.exe
# Listens on localhost:8080
```

**Terminal 2-4 - Start Workers:**
```powershell
.\bin\worker.exe --id worker-1 --port 8081
.\bin\worker.exe --id worker-2 --port 8082  
.\bin\worker.exe --id worker-3 --port 8083
```

### Submitting Jobs

Install grpcurl:
```powershell
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Submit a job:
```powershell
grpcurl -plaintext -d '{\"command\":\"echo Hello from Titan\"}' localhost:8080 titan.ManagerService/SubmitJob
```

Check job status:
```powershell
grpcurl -plaintext -d '{\"job_id\":\"<JOB_ID>\"}' localhost:8080 titan.ManagerService/GetJobStatus
```

List all jobs:
```powershell
grpcurl -plaintext localhost:8080 titan.ManagerService/ListJobs
```

## System Highlights

### 1. Fault Tolerance
- Workers send heartbeats every 10 seconds
- Manager marks workers unhealthy after 30s timeout
- Failed workers automatically excluded from scheduling
- Tasks can be rescheduled if worker dies

### 2. Concurrency
- Thread-safe in-memory storage using `sync.RWMutex`
- Scheduler runs as background goroutine (5s intervals)
- Each task executes in its own goroutine
- Manager handles multiple workers concurrently

### 3. Observability
- Structured JSON logging throughout
- Task output captured (stdout/stderr)
- Exit codes tracked for failure analysis

## Design Patterns Demonstrated

✅ **Client-Server Architecture** - Manager coordinates distributed workers  
✅ **Control Plane / Data Plane Separation** - Scheduling logic vs execution  
✅ **Heartbeat Protocol** - Active health checking  
✅ **Round-Robin Load Balancing** - Fair task distribution  
✅ **RPC Communication** - gRPC for type-safe messaging  
✅ **Process Orchestration** - Spawning & managing OS processes  
✅ **State Machine** - Job lifecycle (PENDING → SCHEDULED → RUNNING → COMPLETED/FAILED)

## Architecture Trade-Offs

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed analysis including:
- CAP theorem choices (Availability over Consistency)
- Single-manager design vs distributed consensus
- In-memory state with WAL vs external database
- Round-robin vs resource-aware scheduling

## Future Enhancements

1. **High Availability**: Multi-manager setup with Raft consensus
2. **Persistence**: Write-Ahead Log (WAL) for crash recovery  
3. **Advanced Scheduling**: Resource constraints (CPU/memory limits)
4. **Metrics**: Prometheus integration
5. **Job Dependencies**: DAG-based workflows

## Troubleshooting

**Issue**: `protoc: command not found`  
**Solution**: Pre-generated gRPC code is included in `pkg/proto/`. No need to regenerate.

**Issue**: `go: missing go.sum entry`  
**Solution**: Run `go mod download` or manuallylet Go resolve dependencies during build.

**Issue**: Workers not registering  
**Solution**: Ensure Manager is running before starting Workers. Check firewall settings for localhost ports.

## License

MIT - Built as a systems architecture demonstration project.
