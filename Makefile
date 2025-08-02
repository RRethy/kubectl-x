.PHONY: test lint lint-fix build fmt vet help

# Default target
help:
	@echo "Available targets:"
	@echo "  test      - Run all tests"
	@echo "  lint      - Run golangci-lint"
	@echo "  lint-fix  - Run golangci-lint with auto-fix"
	@echo "  build     - Build the kubectl-x binary"
	@echo "  fmt       - Format Go code"
	@echo "  vet       - Run go vet"

# Run all tests
test:
	cd kubectl-x && go test ./...

# Run golangci-lint
lint:
	cd kubectl-x && go run github.com/golangci/golangci-lint/cmd/golangci-lint run

# Run golangci-lint with auto-fix
lint-fix:
	cd kubectl-x && go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix

# Build the binary
build:
	cd kubectl-x && go build -o ../kubectl-x .

# Format Go code
fmt:
	cd kubectl-x && go fmt ./...

# Run go vet
vet:
	cd kubectl-x && go vet ./...