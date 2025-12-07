package worker

import (
	"net/rpc"
	"time"

	"titan/pkg/logger"
	pb "titan/pkg/proto"
)

// Heartbeater manages periodic heartbeats to the manager
type Heartbeater struct {
	workerID      string
	managerClient *rpc.Client
	interval      time.Duration
	stopChan      chan struct{}
}

// NewHeartbeater creates a new heartbeater
func NewHeartbeater(workerID string, client *rpc.Client, interval time.Duration) *Heartbeater {
	return &Heartbeater{
		workerID:      workerID,
		managerClient: client,
		interval:      interval,
		stopChan:      make(chan struct{}),
	}
}

// Start begins sending heartbeats
func (h *Heartbeater) Start() {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()
	
	logger.Info("Heartbeater started", "worker_id", h.workerID, "interval", h.interval.String())
	
	for {
		select {
		case <-ticker.C:
			h.sendHeartbeat()
		case <-h.stopChan:
			logger.Info("Heartbeater stopped", "worker_id", h.workerID)
			return
		}
	}
}

// Stop stops the heartbeater
func (h *Heartbeater) Stop() {
	close(h.stopChan)
}

// sendHeartbeat sends a single heartbeat to the manager
func (h *Heartbeater) sendHeartbeat() {
	req := pb.HeartbeatRequest{
		WorkerId:  h.workerID,
		Timestamp: time.Now().Unix(),
		CurrentUsage: pb.ResourceUsage{
			UsedCpuMillicores: 100, // Mock usage
			UsedMemoryMb:      256, // Mock usage
		},
	}
	
	var resp pb.HeartbeatResponse
	err := h.managerClient.Call("ManagerService.Heartbeat", req, &resp)
	if err != nil {
		logger.Error("Failed to send heartbeat", "error", err)
	}
}
