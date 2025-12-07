# Titan - Quick Test Demo Script
# This script demonstrates the system by testing basic job submission

Write-Host "=== Titan Distributed Job Scheduler - Quick Test ===" -ForegroundColor Cyan
Write-Host ""

# Check if binaries exist
if (!(Test-Path "bin\manager.exe") -or !(Test-Path "bin\worker.exe")) {
    Write-Host "❌ Binaries not found! Please build first:" -ForegroundColor Red
    Write-Host "   go build -o bin\manager.exe .\cmd\manager" -ForegroundColor Yellow
    Write-Host "   go build -o bin\worker.exe .\cmd\worker" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Binaries found!" -ForegroundColor Green
Write-Host ""
Write-Host "This demo will:" -ForegroundColor Yellow
Write-Host "  1. Start the Manager" -ForegroundColor Gray
Write-Host "  2. Start 2 Workers" -ForegroundColor Gray
Write-Host "  3. Run for 30 seconds so you can test it" -ForegroundColor Gray
Write-Host "  4. Automatically clean up" -ForegroundColor Gray
Write-Host ""
Write-Host "To submit jobs while running, open another terminal and use:" -ForegroundColor Cyan
Write-Host '  .\bin\client.exe --command "echo Hello"' -ForegroundColor Gray
Write-Host ""
Write-Host "Press any key to start..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Write-Host ""

# Cleanup function
function Cleanup {
    Write-Host ""
    Write-Host "=== Shutting down ===" -ForegroundColor Yellow
    Get-Job | Stop-Job
    Get-Job | Remove-Job
    Write-Host "✅ All processes stopped" -ForegroundColor Green
}

# Register cleanup on Ctrl+C
$null = Register-EngineEvent PowerShell.Exiting -Action { Cleanup }

try {
    # Start Manager
    Write-Host "▶ Starting Manager on port 8080..." -ForegroundColor Green
    $managerJob = Start-Job -ScriptBlock {
        Set-Location "c:\Users\LEGION\OneDrive\Desktop\system architecture"
        .\bin\manager.exe
    }
    Start-Sleep -Seconds 2

    # Start Worker 1
    Write-Host "▶ Starting Worker 1 on port 8081..." -ForegroundColor Green
    $worker1Job = Start-Job -ScriptBlock {
        Set-Location "c:\Users\LEGION\OneDrive\Desktop\system architecture"
        .\bin\worker.exe --id worker-1 --port 8081
    }
    Start-Sleep -Seconds 1

    # Start Worker 2
    Write-Host "▶ Starting Worker 2 on port 8082..." -ForegroundColor Green
    $worker2Job = Start-Job -ScriptBlock {
        Set-Location "c:\Users\LEGION\OneDrive\Desktop\system architecture"
        .\bin\worker.exe --id worker-2 --port 8082
    }
    Start-Sleep -Seconds 2

    Write-Host ""
    Write-Host "🚀 Cluster is running!" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Manager: http://localhost:8080" -ForegroundColor Gray
    Write-Host "Worker 1: http://localhost:8081" -ForegroundColor Gray
    Write-Host "Worker 2: http://localhost:8082" -ForegroundColor Gray
    Write-Host ""
    Write-Host "The cluster will run for 30 seconds..." -ForegroundColor Yellow
    Write-Host "Press Ctrl+C to stop early" -ForegroundColor Gray
    Write-Host ""

    # Show logs
    Write-Host "=== Recent Manager Logs ===" -ForegroundColor Cyan
    Receive-Job -Job $managerJob
    Write-Host ""
    Write-Host "=== Recent Worker Logs ===" -ForegroundColor Cyan
    Receive-Job -Job $worker1Job
    Receive-Job -Job $worker2Job

    Write-Host ""
    Write-Host "Cluster running... (30 seconds remaining)" -ForegroundColor Yellow
    
    # Run for 30 seconds
    for ($i = 30; $i -gt 0; $i--) {
        Start-Sleep -Seconds 1
        if ($i % 5 -eq 0) {
            Write-Host "  $i seconds remaining..." -ForegroundColor Gray
        }
    }

    Write-Host ""
    Write-Host "✅ Demo complete!" -ForegroundColor Green

} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
} finally {
    Cleanup
}

Write-Host ""
Write-Host "To run manually, see QUICKSTART.md" -ForegroundColor Cyan
