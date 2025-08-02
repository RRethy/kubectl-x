.PHONY: test lint lint-fix build build-kubectl-x build-kubernetes-mcp fmt vet help

# Default target
help:
	@echo "Available targets:"
	@echo "  test                 - Run all tests"
	@echo "  lint                 - Run golangci-lint"
	@echo "  lint-fix             - Run golangci-lint with auto-fix"
	@echo "  build                - Build all binaries"
	@echo "  build-kubectl-x      - Build the kubectl-x binary"
	@echo "  build-kubernetes-mcp - Build the kubernetes-mcp binary"
	@echo "  fmt                  - Format Go code"
	@echo "  vet                  - Run go vet"

# Run all tests
test:
	cd kubectl-x && go test ./...
	cd kubernetes-mcp && go test ./...

# Run golangci-lint
lint:
	cd kubectl-x && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 10m
	cd kubernetes-mcp && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 10m

# Run golangci-lint with auto-fix
lint-fix:
	cd kubectl-x && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix --timeout 10m
	cd kubernetes-mcp && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix --timeout 10m

# Build all binaries
build: build-kubectl-x build-kubernetes-mcp

# Build the kubectl-x binary
build-kubectl-x:
	cd kubectl-x && go build -o ../kubectl-x .

# Build the kubernetes-mcp binary
build-kubernetes-mcp:
	cd kubernetes-mcp && go build -o ../kubernetes-mcp .

# Format Go code
fmt:
	cd kubectl-x && go fmt ./...
	cd kubernetes-mcp && go fmt ./...

# Run go vet
vet:
	cd kubectl-x && go vet ./...
	cd kubernetes-mcp && go vet ./...