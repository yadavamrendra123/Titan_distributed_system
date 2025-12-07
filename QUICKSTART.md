# 🚀 How to Run Titan - Quick Start Guide

## ✅ Prerequisites Check

Your system is ready! Both binaries have been built successfully:
- ✅ `bin\manager.exe` - Control Plane
- ✅ `bin\worker.exe` - Data Plane

---

## 📝 Step-by-Step Instructions

### Step 1: Start the Manager

Open a **new PowerShell terminal** and run:

```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\bin\manager.exe
```

**Expected Output:**
```json
{"level":"info","msg":"Starting Titan Manager","address":"0.0.0.0:8080"}
{"level":"info","msg":"Manager server started"}
{"level":"info","msg":"Scheduler started"}
{"level":"info","msg":"Manager listening","address":"0.0.0.0:8080"}
```

✅ **Keep this terminal open!** Manager is now running on port 8080.

---

### Step 2: Start Worker 1

Open a **second PowerShell terminal** and run:

```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\bin\worker.exe --id worker-1 --port 8081
```

**Expected Output:**
```json
{"level":"info","msg":"Starting Titan Worker","worker_id":"worker-1","address":"localhost:8081","manager":"localhost:8080"}
{"level":"info","msg":"Worker registered successfully","worker_id":"worker-1"}
{"level":"info","msg":"Heartbeater started","worker_id":"worker-1","interval":"10s"}
{"level":"info","msg":"Worker listening","worker_id":"worker-1","address":"localhost:8081"}
```

✅ Worker 1 is now registered and sending heartbeats!

---

### Step 3: Start Worker 2 (Optional)

Open a **third PowerShell terminal**:

```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\bin\worker.exe --id worker-2 --port 8082
```

---

### Step 4: Start Worker 3 (Optional)

Open a **fourth PowerShell terminal**:

```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\bin\worker.exe --id worker-3 --port 8083
```

---

## 🎯 Submit Your First Job

### Option A: Using the CLI Client (Recommended)

In a **new terminal**, run:

```powershell
.\bin\client.exe --command "echo Hello from Titan!"
```

**Expected Response:**
```
Job submitted successfully!
Job ID: 550e8400-e29b-41d4-a716-446655440000
Status: PENDING
```

### Option B: Check Job Status

```powershell
.\bin\client.exe --status <JOB_ID>
```

### Option C: List All Jobs

```powershell
.\bin\client.exe --list
```

---

## 📊 Observe the System

### Watch Manager Logs

In the Manager terminal, you'll see:

```json
{"level":"info","msg":"Job submitted","job_id":"550e8400...","command":"echo Hello from Titan!"}
{"level":"info","msg":"Job scheduled","job_id":"550e8400...","worker_id":"worker-1"}
```

### Watch Worker Logs

In Worker 1 terminal, you'll see:

```json
{"level":"info","msg":"Received task","task_id":"...","job_id":"550e8400..."}
{"level":"info","msg":"Starting task","task_id":"...","command":"echo Hello from Titan!"}
{"level":"info","msg":"Task completed","task_id":"..."}
```

---

## 🧪 Test Scenarios

### 1. Submit Multiple Jobs

```powershell
grpcurl -plaintext -d '{\"command\":\"echo Job 1\"}' localhost:8080 titan.ManagerService/SubmitJob
grpcurl -plaintext -d '{\"command\":\"echo Job 2\"}' localhost:8080 titan.ManagerService/SubmitJob
grpcurl -plaintext -d '{\"command\":\"echo Job 3\"}' localhost:8080 titan.ManagerService/SubmitJob
```

Watch the scheduler distribute them round-robin across workers!

### 2. Long-Running Job

```powershell
grpcurl -plaintext -d '{\"command\":\"ping -n 10 localhost\"}' localhost:8080 titan.ManagerService/SubmitJob
```

(Windows equivalent of `sleep 10`)

### 3. Failing Job

```powershell
grpcurl -plaintext -d '{\"command\":\"exit 1\"}' localhost:8080 titan.ManagerService/SubmitJob
```

Check the status - it should show `"status":"FAILED"` with `"exitCode":1`.

### 4. List All Jobs

```powershell
grpcurl -plaintext localhost:8080 titan.ManagerService/ListJobs
```

### 5. Test Fault Tolerance

1. Submit a long-running job
2. Press `Ctrl+C` in one Worker terminal to kill it
3. Submit another job
4. Observe that the Manager only schedules to healthy workers!

---

## 🛑 Stopping the System

**Graceful Shutdown:**
- Press `Ctrl+C` in each terminal (Manager and Workers)
- The system will cleanly shut down

**Output:**
```json
{"level":"info","msg":"Shutting down Manager..."}
{"level":"info","msg":"Heartbeater stopped","worker_id":"worker-1"}
```

---

## 🎨 Advanced Usage

### Custom Commands

```powershell
# Run PowerShell commands
grpcurl -plaintext -d '{\"command\":\"Get-Date\"}' localhost:8080 titan.ManagerService/SubmitJob

# File operations (use sh -c on Git Bash)
grpcurl -plaintext -d '{\"command\":\"dir\"}' localhost:8080 titan.ManagerService/SubmitJob
```

### Environment Variables

```powershell
grpcurl -plaintext -d '{\"command\":\"echo $MY_VAR\",\"env\":{\"MY_VAR\":\"Hello\"}}' localhost:8080 titan.ManagerService/SubmitJob
```

---

## 🐛 Troubleshooting

### Issue: "Failed to connect to Manager"
**Solution:** Ensure Manager is running on port 8080
```powershell
netstat -an | findstr "8080"
```

### Issue: "Worker not registering"
**Solution:** Check firewall settings for localhost ports

### Issue: Commands not executing
**Solution:** On Windows, use `cmd /c` or PowerShell syntax:
```powershell
grpcurl -plaintext -d '{\"command\":\"cmd /c echo Hello\"}' localhost:8080 titan.ManagerService/SubmitJob
```

---

## 📈 What's Happening Under the Hood

1. **Manager** starts and begins the scheduler loop (every 5 seconds)
2. **Workers** register and send heartbeats (every 10 seconds)
3. You **submit a job** → Manager marks it `PENDING`
4. **Scheduler** wakes up → finds pending jobs + healthy workers
5. **Round-robin** assignment → Manager sends `StartTask` RPC to Worker
6. **Worker** spawns OS process → captures output
7. **Worker** reports status back to Manager
8. You **query status** → see `COMPLETED` with output!

---

## 🎓 Learning Exercises

1. **Monitor heartbeats:** Watch the logs and count heartbeats (should be every 10s)
2. **Test scheduling:** Submit 5 jobs with 2 workers - observe round-robin
3. **Break and recover:** Kill a worker mid-task, see how Manager handles it
4. **Resource limits:** Try to overload workers with many concurrent jobs

---

## 🚀 Next Steps

- **Add persistence:** Implement Write-Ahead Log (WAL)
- **Add metrics:** Integrate Prometheus
- **Add a Web UI:** Build a dashboard to visualize jobs
- **Deploy to cloud:** Run Manager and Workers on separate VMs

---

**Congratulations!** 🎉 You now have a fully functional distributed job scheduler running on your machine!
