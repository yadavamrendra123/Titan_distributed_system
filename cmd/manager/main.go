package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"titan/pkg/logger"
	"titan/pkg/manager"
)

const (
	defaultPort = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)

	logger.Info("Starting Titan Manager", "address", address)

	server := manager.NewServer()
	server.Start()

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

		logger.Info("Shutting down Manager...")
		os.Exit(0)
	}()

	logger.Info("Manager listening", "address", address)

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
