package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"titan/pkg/logger"
	"titan/pkg/worker"
)

const (
	defaultManagerAddr = "localhost:8080"
)

func main() {
	workerID := flag.String("id", "", "Worker ID (required)")
	port := flag.String("port", "8081", "Worker port")
	managerAddr := flag.String("manager", defaultManagerAddr, "Manager address")
	flag.Parse()

	if *workerID == "" {
		logger.Error("Worker ID is required (use --id flag)")
		os.Exit(1)
	}

	address := fmt.Sprintf("localhost:%s", *port)

	logger.Info("Starting Titan Worker",
		"worker_id", *workerID,
		"address", address,
		"manager", *managerAddr)

	server, err := worker.NewServer(*workerID, address, *managerAddr)
	if err != nil {
		logger.Error("Failed to create worker server", "error", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		logger.Error("Failed to start worker", "error", err)
		os.Exit(1)
	}

	// Create net/rpc server
	rpcServer := server.GetRPCServer()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down Worker...", "worker_id", *workerID)
		os.Exit(0)
	}()

	logger.Info("Worker listening", "worker_id", *workerID, "address", address)

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Accept error", "error", err)
			continue
		}
		go rpcServer.ServeConn(conn)
	}
}
