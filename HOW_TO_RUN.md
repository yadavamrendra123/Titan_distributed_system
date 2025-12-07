# 🎯 Titan - How to Run

## ✅ Build Complete!

Your binaries are ready:
- `bin\manager.exe` ✅
- `bin\worker.exe` ✅

---

## 🚀 Quick Start (3 Simple Steps)

### **Method 1: Automated Demo (Easiest)**

Open PowerShell and run:

```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\demo.ps1
```

This will:
- Start Manager and 2 Workers automatically
- Run for 30 seconds
- Clean up automatically

---

### **Method 2: Manual (Full Control)**

#### Step 1: Start Manager
**Terminal 1:**
```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\bin\manager.exe
```

#### Step 2: Start Worker(s)
**Terminal 2:**
```powershell
cd "c:\Users\LEGION\OneDrive\Desktop\system architecture"
.\bin\worker.exe --id worker-1 --port 8081
```

**Terminal 3 (optional):**
```powershell
.\bin\worker.exe --id worker-2 --port 8082
```

#### Step 3: Submit a Job

**Using the CLI Client:**
```powershell
.\bin\client.exe --command "echo Hello Titan!"
```

**Check Status:**
```powershell
.\bin\client.exe --list
```

---

## 🎬 What You'll See

**Manager Terminal:**
```json
{"level":"info","msg":"Manager listening","address":"0.0.0.0:8080"}
{"level":"info","msg":"Job submitted","job_id":"abc123..."}
{"level":"info","msg":"Job scheduled","job_id":"abc123...","worker_id":"worker-1"}
```

**Worker Terminal:**
```json
{"level":"info","msg":"Worker registered successfully","worker_id":"worker-1"}
{"level":"info","msg":"Received task","task_id":"xyz789..."}
{"level":"info","msg":"Task completed","task_id":"xyz789..."}
```

---

## 📖 Full Documentation

See **[QUICKSTART.md](file:///c:/Users/LEGION/OneDrive/Desktop/system%20architecture/QUICKSTART.md)** for:
- Complete step-by-step guide
- Test scenarios (fault tolerance, multiple jobs, etc.)
- Troubleshooting
- Advanced usage

---

## 🛑 Stopping

Press `Ctrl+C` in each terminal or close the windows.

---

**That's it!** Your distributed job scheduler is running! 🎉
