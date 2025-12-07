# Build Script for Titan - Windows PowerShell
# This script builds the project even with network/dependency issues

Write-Host "=== Building Titan Distributed Job Scheduler ===" -ForegroundColor Cyan
Write-Host ""

# Create bin directory
Write-Host "Creating bin directory..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path bin | Out-Null

# Try to download dependencies first (may fail due to permissions)
Write-Host "Attempting to download dependencies..." -ForegroundColor Yellow
& go mod download 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠ Dependency download had issues, but continuing..." -ForegroundColor Yellow
}

Write-Host ""

# Build Manager
Write-Host "Building Manager..." -ForegroundColor Green
& go build -o bin\manager.exe .\cmd\manager 2>&1 | Tee-Object -Variable managerOutput
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Manager binary created: bin\manager.exe" -ForegroundColor Green
} else {
    Write-Host "✗ Manager build failed" -ForegroundColor Red
    Write-Host $managerOutput -ForegroundColor Red
    Write-Host ""
    Write-Host "Common fixes:" -ForegroundColor Yellow
    Write-Host "  1. Ensure Go 1.21+ is installed: go version" -ForegroundColor Gray
    Write-Host "  2. Run as Administrator if permission denied" -ForegroundColor Gray
    Write-Host "  3. Check internet connection for dependency downloads" -ForegroundColor Gray
    exit 1
}

Write-Host ""

# Build Worker  
Write-Host "Building Worker..." -ForegroundColor Green
& go build -o bin\worker.exe .\cmd\worker 2>&1 | Tee-Object -Variable workerOutput
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Worker binary created: bin\worker.exe" -ForegroundColor Green
} else {
    Write-Host "✗ Worker build failed" -ForegroundColor Red
    Write-Host $workerOutput -ForegroundColor Red
    exit 1
}

Write-Host ""

# Build Client
Write-Host "Building Client..." -ForegroundColor Green
& go build -o bin\client.exe .\cmd\client 2>&1 | Tee-Object -Variable clientOutput
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Client binary created: bin\client.exe" -ForegroundColor Green
} else {
    Write-Host "✗ Client build failed" -ForegroundColor Red
    Write-Host $clientOutput -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== Build Complete ===" -ForegroundColor Cyan
Write-Host ""
# Build Renderer
Write-Host "Building Renderer..." -ForegroundColor Green
& go build -o bin\renderer.exe .\cmd\renderer 2>&1 | Tee-Object -Variable rendererOutput
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Renderer binary created: bin\renderer.exe" -ForegroundColor Green
} else {
    Write-Host "✗ Renderer build failed" -ForegroundColor Red
    Write-Host $rendererOutput -ForegroundColor Red
    exit 1
}

Write-Host ""

# Build Orchestrator
Write-Host "Building Orchestrator..." -ForegroundColor Green
& go build -o bin\orchestrator.exe .\cmd\orchestrator 2>&1 | Tee-Object -Variable orchestratorOutput
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Orchestrator binary created: bin\orchestrator.exe" -ForegroundColor Green
} else {
    Write-Host "✗ Orchestrator build failed" -ForegroundColor Red
    Write-Host $orchestratorOutput -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== Build Complete ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Run the cluster with:" -ForegroundColor Yellow
Write-Host "  Manager:      .\bin\manager.exe" -ForegroundColor Gray
Write-Host "  Worker:       .\bin\worker.exe --id worker-1 --port 8081" -ForegroundColor Gray
Write-Host "  Orchestrator: .\bin\orchestrator.exe" -ForegroundColor Gray
