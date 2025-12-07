package main

import (
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
	"time"

	pb "titan/pkg/proto"
)

// Orchestrator manages the distributed rendering job
func main() {
	managerAddr := "localhost:8080"
	outputDir := "fractals"
	
	// Image parameters (4K resolution)
	fullWidth := 4096
	fullHeight := 4096
	totalIter := 1000
	
	// Viewport (Zoomed into a nice spiral)
	minX, minY := -0.748, 0.1
	maxX, maxY := -0.744, 0.104
	
	// Grid (4x4 = 16 workers/jobs needed)
	rows, cols := 4, 4
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	// Connect to Manager
	client, err := rpc.Dial("tcp", managerAddr)
	if err != nil {
		fmt.Printf("Error connecting to manager: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Printf("=== Starting Distributed Fractal Rendering ===\n")
	fmt.Printf("Resolution: %dx%d\n", fullWidth, fullHeight)
	fmt.Printf("Tiles: %dx%d (%d jobs)\n", rows, cols, rows*cols)
	
	startTime := time.Now()
	var wg sync.WaitGroup
	
	tileWidth := fullWidth / cols
	tileHeight := fullHeight / rows
	
	dx := (maxX - minX) / float64(cols)
	dy := (maxY - minY) / float64(rows)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			wg.Add(1)
			
			// Calculate tile coordinates
			tileMinX := minX + float64(c)*dx
			tileMaxX := tileMinX + dx
			tileMinY := minY + float64(r)*dy
			tileMaxY := tileMinY + dy
			
			outputFile := filepath.Join(outputDir, fmt.Sprintf("tile_%d_%d.png", r, c))
			
			// Construct command
			// Note: We use absolute path to renderer.exe assuming it's in the same bin dir as worker
			// In a real cluster, we'd deploy the binary. Here we rely on shared storage (localhost).
			cwd, _ := os.Getwd()
			rendererPath := filepath.Join(cwd, "bin", "renderer.exe")
			
			cmd := fmt.Sprintf("\"%s\" -minx %f -miny %f -maxx %f -maxy %f -w %d -h %d -iter %d -out \"%s\"",
				rendererPath, tileMinX, tileMinY, tileMaxX, tileMaxY, tileWidth, tileHeight, totalIter, outputFile)
			
			go func(r, c int, cmd string) {
				defer wg.Done()
				
				fmt.Printf("[Tile %d,%d] Submitting job...\n", r, c)
				
				req := pb.JobRequest{
					Command: cmd,
					Env:     nil,
				}
				var resp pb.JobResponse
				err := client.Call("ManagerService.SubmitJob", req, &resp)
				if err != nil {
					fmt.Printf("[Tile %d,%d] Failed to submit: %v\n", r, c, err)
					return
				}
				
				// Poll for completion
				for {
					statusReq := pb.JobStatusRequest{JobId: resp.JobId}
					var statusResp pb.JobStatusResponse
					client.Call("ManagerService.GetJobStatus", statusReq, &statusResp)
					
					if statusResp.Status == "COMPLETED" {
						fmt.Printf("[Tile %d,%d] ✅ Finished (Worker: %s)\n", r, c, statusResp.WorkerId)
						break
					} else if statusResp.Status == "FAILED" {
						fmt.Printf("[Tile %d,%d] ❌ Failed: %s\n", r, c, statusResp.Output)
						break
					}
					
					time.Sleep(500 * time.Millisecond)
				}
			}(r, c, cmd)
		}
	}
	
	wg.Wait()
	fmt.Printf("All jobs complete in %v\n", time.Since(startTime))
	fmt.Printf("Output available in .\\%s\\\n", outputDir)
}
