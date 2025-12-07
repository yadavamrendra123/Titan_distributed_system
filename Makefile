.PHONY: proto build clean test run-manager run-worker

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
MANAGER_BINARY=bin/manager
WORKER_BINARY=bin/worker

# Protocol Buffers
PROTO_DIR=proto
PROTO_OUT=pkg/proto

all: proto build

# Generate Go code from Protocol Buffers
proto:
	@echo "Generating gRPC code from protobuf..."
	@mkdir -p $(PROTO_OUT)
	protoc --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/titan.proto
	@echo "✓ Protocol buffers compiled"

# Build both binaries
build:
	@echo "Building Manager..."
	@mkdir -p bin
	$(GOBUILD) -o $(MANAGER_BINARY) ./cmd/manager
	@echo "✓ Manager built at $(MANAGER_BINARY)"
	@echo "Building Worker..."
	$(GOBUILD) -o $(WORKER_BINARY) ./cmd/worker
	@echo "✓ Worker built at $(WORKER_BINARY)"

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -rf $(PROTO_OUT)

# Run manager
run-manager:
	$(MANAGER_BINARY)

# Run worker (pass ID and PORT as env vars)
run-worker:
	$(WORKER_BINARY) --id $(ID) --port $(PORT)

# Initialize Go modules
init:
	$(GOMOD) init github.com/yourusername/titan
	$(GOMOD) tidy

# Download dependencies
deps:
	$(GOGET) google.golang.org/grpc
	$(GOGET) google.golang.org/protobuf
	$(GOGET) github.com/google/uuid
	$(GOMOD) tidy
