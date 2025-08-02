# CLAUDE.md - kubectl-x Module

This file provides module-specific guidance for Claude Code when working within the kubectl-x module directory.

**IMPORTANT**: When making changes to module-specific aspects (dependencies, commands, internal packages), update this file. For workspace-level changes, update the root `CLAUDE.md`.

## Usage Examples
```bash
kubectl x ctx                    # Interactive context selection
kubectl x ctx my-context        # Switch to context with partial match
kubectl x ns                     # Interactive namespace selection
kubectl x ns my-namespace       # Switch to namespace with partial match
kubectl x cur                    # Show current context and namespace
kubectl x ctx -                  # Switch to previous context/namespace
```

## Module Commands

### Build and Test (from kubectl-x/ directory)
```bash
# Build the application
go build .

# Run all module tests
go test ./...

# Run specific package tests
go test ./internal/cmd/ctx/
go test ./internal/fzf/
go test ./internal/history/
go test ./internal/kubeconfig/
go test ./internal/kubernetes/

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Benchmark tests
go test -bench=. ./...

# Update dependencies
go mod tidy
go mod download
```

## Module Architecture

### Entry Point
- `main.go` - Simple entry point that calls `cmd.Execute()`

### Command Structure (`cmd/`)
- `root.go` - Root command setup with kubectl CLI options integration
- `ctx.go` - Context switching command definition
- `ns.go` - Namespace switching command definition  
- `cur.go` - Current status display command definition

### Internal Package Details

#### `internal/cmd/ctx/`
- `ctx.go` - Context command implementation
- `ctxer.go` - Context switching business logic
- `ctxer_test.go` - Comprehensive test suite
- **Key interfaces**: `Ctxer` for context operations

#### `internal/cmd/ns/`
- `ns.go` - Namespace command implementation
- `nser.go` - Namespace switching business logic
- `nser_test.go` - Test suite with table-driven tests
- **Key interfaces**: `Nser` for namespace operations

#### `internal/cmd/cur/`
- `cur.go` - Current status command implementation
- `curer.go` - Current status display logic
- `curer_test.go` - Status display tests
- **Key interfaces**: `Curer` for status operations

#### `internal/fzf/`
- `fzf.go` - Fuzzy finder integration with external fzf binary
- `fzf_test.go` - Tests including user cancellation scenarios
- `testing/fzf.go` - Mock implementation for testing
- **Key interfaces**: `Fzf` for interactive selection

#### `internal/history/`
- `history.go` - Command history persistence and retrieval
- `history_test.go` - History management tests
- `testing/history.go` - Mock history implementation
- **Storage**: `~/.local/share/kubectl-x/history.yaml`
- **Key interfaces**: `History` for history operations

#### `internal/kubeconfig/`
- `kubeconfig.go` - Kubeconfig file manipulation
- `kubeconfig_test.go` - Kubeconfig operation tests
- `testing/kubeconfig.go` - Mock kubeconfig implementation
- **Key interfaces**: `Kubeconfig` for kubeconfig operations

#### `internal/kubernetes/`
- `client.go` - Kubernetes API client wrapper
- `kubernetes.go` - Generic resource operations
- `testing/client.go` - Mock Kubernetes client
- **Key interfaces**: `Client` for Kubernetes operations

## Development Patterns

### Interface Implementation Pattern
Each internal package follows this pattern:
1. Define primary interface (e.g., `Ctxer`, `Nser`)
2. Implement concrete struct
3. Create constructor function with dependency injection
4. Provide mock implementation in `testing/` subdirectory

### Error Handling Pattern
- Wrap errors with context using `fmt.Errorf`
- Return user-friendly error messages
- Handle kubectl client errors consistently

### Testing Pattern
- Table-driven tests for multiple scenarios
- Mock all external dependencies
- Test both success and error conditions
- Use testify/assert for assertions

### Command Integration Pattern
- Commands are thin wrappers around internal packages
- Business logic in `internal/cmd/{command}/` packages
- Dependency injection through constructor functions
- Consistent flag handling using Cobra

## Module Dependencies

### Direct Dependencies
```go
github.com/fatih/color v1.17.0          // Terminal colors
github.com/goccy/go-yaml v1.11.3        // YAML processing
github.com/spf13/cobra v1.8.1           // CLI framework
github.com/stretchr/testify v1.9.0      // Testing utilities
k8s.io/api v0.30.2                      // Kubernetes API types
k8s.io/apimachinery v0.30.2             // Kubernetes API machinery
k8s.io/cli-runtime v0.30.2              // kubectl CLI runtime
k8s.io/client-go v0.30.2                // Kubernetes Go client
k8s.io/kubectl v0.30.2                  // kubectl utilities
k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // Kubernetes utilities
```

### External Runtime Dependencies
- `fzf` binary - Required for interactive selection
- `kubectl` - Uses kubectl's configuration and patterns

## Code Generation
- `internal/cmd/generate/generate.go` - Placeholder for potential code generation
- Currently unused but reserved for future enhancements