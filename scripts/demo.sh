#!/bin/bash

# Demo script for Titan - Distributed Job Scheduler
# This script demonstrates the system by starting a cluster and submitting jobs

set -e

echo "=== Titan Demo ==="
echo ""

# Check if binaries exist
if [ ! -f "bin/manager" ] || [ ! -f "bin/worker" ]; then
    echo "Error: Binaries not found. Please run 'make build' first."
    exit 1
fi

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo "Error: grpcurl is required but not installed."
    echo "Install with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "=== Cleaning up ==="
    kill $(jobs -p) 2>/dev/null || true
    wait 2>/dev/null || true
    echo "Demo complete"
}

trap cleanup EXIT

# Start Manager
echo "Starting Manager on port 8080..."
./bin/manager &
MANAGER_PID=$!
sleep 2

# Start 3 Workers
echo "Starting Worker 1 on port 8081..."
./bin/worker --id worker-1 --port 8081 &
sleep 1

echo "Starting Worker 2 on port 8082..."
./bin/worker --id worker-2 --port 8082 &
sleep 1

echo "Starting Worker 3 on port 8083..."
./bin/worker --id worker-3 --port 8083 &
sleep 2

echo ""
echo "=== Cluster is ready ==="
echo ""

# Submit test jobs
echo "Submitting jobs..."

echo "Job 1: Simple echo"
grpcurl -plaintext -d '{"command":"echo Hello from Titan"}' \
    localhost:8080 titan.ManagerService/SubmitJob

echo ""
echo "Job 2: Sleep and echo"
grpcurl -plaintext -d '{"command":"sleep 2 && echo Job completed"}' \
    localhost:8080 titan.ManagerService/SubmitJob

echo ""
echo "Job 3: Date command"
grpcurl -plaintext -d '{"command":"date"}' \
    localhost:8080 titan.ManagerService/SubmitJob

echo ""
echo "Waiting for jobs to complete..."
sleep 5

# List all jobs
echo ""
echo "=== Job Status ==="
grpcurl -plaintext localhost:8080 titan.ManagerService/ListJobs

echo ""
echo "Demo will continue running for 10 more seconds..."
echo "Check the logs above to see task execution"
sleep 10
