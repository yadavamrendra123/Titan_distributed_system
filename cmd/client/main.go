package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"os"

	pb "titan/pkg/proto"
)

func main() {
	managerAddr := flag.String("manager", "localhost:8080", "Manager address")
	command := flag.String("command", "", "Command to run")
	list := flag.Bool("list", false, "List all jobs")
	status := flag.String("status", "", "Get status of job ID")
	flag.Parse()

	client, err := rpc.Dial("tcp", *managerAddr)
	if err != nil {
		fmt.Printf("Error connecting to manager: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	if *list {
		var req pb.ListJobsRequest
		var resp pb.ListJobsResponse
		err = client.Call("ManagerService.ListJobs", req, &resp)
		if err != nil {
			fmt.Printf("Error listing jobs: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Jobs:\n")
		for _, job := range resp.Jobs {
			fmt.Printf("- %s [%s] Worker: %s ExitCode: %d\n", job.JobId, job.Status, job.WorkerId, job.ExitCode)
		}
		return
	}

	if *status != "" {
		req := pb.JobStatusRequest{JobId: *status}
		var resp pb.JobStatusResponse
		err = client.Call("ManagerService.GetJobStatus", req, &resp)
		if err != nil {
			fmt.Printf("Error getting job status: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Job ID: %s\n", resp.JobId)
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Worker: %s\n", resp.WorkerId)
		fmt.Printf("Exit Code: %d\n", resp.ExitCode)
		fmt.Printf("Output:\n%s\n", resp.Output)
		return
	}

	if *command != "" {
		req := pb.JobRequest{
			Command: *command,
			Env:     make(map[string]string),
		}
		var resp pb.JobResponse
		err = client.Call("ManagerService.SubmitJob", req, &resp)
		if err != nil {
			fmt.Printf("Error submitting job: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Job submitted successfully!\n")
		fmt.Printf("Job ID: %s\n", resp.JobId)
		fmt.Printf("Status: %s\n", resp.Status)
		return
	}

	fmt.Println("Usage:")
	fmt.Println("  Submit job: client.exe --command \"echo hello\"")
	fmt.Println("  List jobs:  client.exe --list")
	fmt.Println("  Job status: client.exe --status <JOB_ID>")
}
