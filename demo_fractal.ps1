# Titan - Distributed Fractal Rendering Demo
# This script demonstrates the system by rendering a high-res Mandelbrot set

Write-Host "=== Titan Distributed Fractal Rendering Demo ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "This demo will:" -ForegroundColor Yellow
Write-Host "  1. Start a 4-node Worker Cluster" -ForegroundColor Gray
Write-Host "  2. Orchestrate a distributed rendering job (16 tiles)" -ForegroundColor Gray
Write-Host "  3. Each worker will compute a piece of the fractal" -ForegroundColor Gray
Write-Host "  4. You will see CPU usage spike and tiles appear" -ForegroundColor Gray
Write-Host ""

# Check binaries
if (!(Test-Path "bin\renderer.exe") -or !(Test-Path "bin\orchestrator.exe")) {
    Write-Host "❌ Binaries not found! Run build.ps1 first." -ForegroundColor Red
    exit 1
}

# Clean output dir
if (Test-Path "fractals") {
    Remove-Item "fractals" -Recurse -Force
}

# Cleanup function
function Cleanup {
    Write-Host ""
    Write-Host "=== Shutting down ===" -ForegroundColor Yellow
    Get-Job | Stop-Job
    Get-Job | Remove-Job
    Write-Host "✅ All processes stopped" -ForegroundColor Green
}

$null = Register-EngineEvent PowerShell.Exiting -Action { Cleanup }

try {
    # Start Manager
    Write-Host "▶ Starting Manager..." -ForegroundColor Green
    $managerJob = Start-Job -ScriptBlock {
        Set-Location "c:\Users\LEGION\OneDrive\Desktop\system architecture"
        .\bin\manager.exe
    }
    Start-Sleep -Seconds 2

    # Start 4 Workers
    for ($i=1; $i -le 4; $i++) {
        $port = 8080 + $i
        Write-Host "▶ Starting Worker $i on port $port..." -ForegroundColor Green
        Start-Job -ScriptBlock {
            param($id, $p)
            Set-Location "c:\Users\LEGION\OneDrive\Desktop\system architecture"
            .\bin\worker.exe --id $id --port $p
        } -ArgumentList "worker-$i", $port
    }
    Start-Sleep -Seconds 3

    Write-Host ""
    Write-Host "🚀 Cluster Ready (4 Workers)" -ForegroundColor Cyan
    Write-Host "Starting Orchestrator..." -ForegroundColor Yellow
    Write-Host ""

    # Run Orchestrator
    # We run this in the foreground so we can see its output
    .\bin\orchestrator.exe

    Write-Host ""
    Write-Host "✅ Demo complete!" -ForegroundColor Green
    Write-Host "Check the 'fractals' folder for the generated images." -ForegroundColor Cyan

} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
} finally {
    Write-Host "Press Enter to exit and cleanup..."
    Read-Host
    Cleanup
}
