package worker

import (
	"bytes"
	"fmt"
	"net/rpc"
	"os/exec"
	"sync"
	"syscall"

	"titan/pkg/logger"
	"titan/pkg/models"
	pb "titan/pkg/proto"
)

// Executor manages task execution
type Executor struct {
	mu            sync.RWMutex
	tasks         map[string]*exec.Cmd
	managerClient *rpc.Client
}

// NewExecutor creates a new executor
func NewExecutor(managerAddr string) (*Executor, error) {
	client, err := rpc.Dial("tcp", managerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to manager: %w", err)
	}

	return &Executor{
		tasks:         make(map[string]*exec.Cmd),
		managerClient: client,
	}, nil
}

// StartTask starts a new task process
func (e *Executor) StartTask(taskID, jobID, command string, env map[string]string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.tasks[taskID]; exists {
		return fmt.Errorf("task %s already running", taskID)
	}

	// Create command
	cmd := exec.Command("cmd", "/C", command) // Windows-specific shell
	
	// Set environment variables
	cmd.Env = append(cmd.Env, fmt.Sprintf("TITAN_JOB_ID=%s", jobID))
	cmd.Env = append(cmd.Env, fmt.Sprintf("TITAN_TASK_ID=%s", taskID))
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	e.tasks[taskID] = cmd

	// Monitor process in background
	go func() {
		// Report RUNNING status
		e.reportStatus(taskID, models.JobStatusRunning, "", 0)

		// Wait for completion
		err := cmd.Wait()

		// Get exit code
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
				} else {
					exitCode = exitError.ExitCode()
				}
			} else {
				exitCode = 1
			}
		}

		output := stdout.String() + stderr.String()
		status := models.JobStatusCompleted
		if exitCode != 0 {
			status = models.JobStatusFailed
		}

		// Report final status
		e.reportStatus(taskID, status, output, int32(exitCode))

		// Cleanup
		e.mu.Lock()
		delete(e.tasks, taskID)
		e.mu.Unlock()
	}()

	return nil
}

// StopTask terminates a running task
func (e *Executor) StopTask(taskID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	cmd, exists := e.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	return nil
}

// reportStatus sends a status update to the manager
func (e *Executor) reportStatus(taskID string, status models.JobStatus, output string, exitCode int32) {
	req := pb.TaskStatusUpdate{
		TaskId:   taskID,
		Status:   string(status),
		Output:   output,
		ExitCode: exitCode,
	}

	var resp pb.Ack
	err := e.managerClient.Call("ManagerService.ReportTaskStatus", req, &resp)
	if err != nil {
		logger.Error("Failed to report task status", 
			"task_id", taskID, 
			"error", err)
	}
}
