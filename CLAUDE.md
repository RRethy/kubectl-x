# CLAUDE.md - Workspace Root

This file provides workspace-level guidance to Claude Code (claude.ai/code) when working with this repository.

**IMPORTANT**: When making workspace-level changes (new modules, workspace structure, build processes), update this file. For module-specific changes, update `kubectl-x/CLAUDE.md`.

## Development Commands

Use the provided Makefile for common development tasks:
```bash
make build                - Build all binaries
make build-kubectl-x      - Build the kubectl-x binary
make build-kubernetes-mcp - Build the kubernetes-mcp binary
make test                 - Run all tests
make lint                 - Run golangci-lint
make lint-fix             - Run golangci-lint with auto-fix
make fmt                  - Format Go code
make vet                  - Run go vet
make help                 - Show all available targets
```

### Workspace-Specific Commands
```bash
# Sync workspace dependencies
go work sync

# Install from source
go install github.com/RRethy/utils/kubectl-x@latest           # kubectl-x CLI
go install github.com/RRethy/utils/kubernetes-mcp@latest       # kubernetes-mcp CLI
```

## Architecture

### Go Workspace Structure
This is a Go workspace with two modules:
- `kubectl-x/` - Kubernetes context and namespace switching CLI
- `kubernetes-mcp/` - Readonly MCP (Model Context Protocol) server for Kubernetes

The root contains `go.work` for workspace configuration.

For detailed module implementation, architecture, and development patterns:
- See `kubectl-x/CLAUDE.md` for kubectl-x specific details
- See `kubernetes-mcp/` module for MCP server implementation

### kubernetes-mcp Module
The `kubernetes-mcp/` module provides a readonly MCP (Model Context Protocol) server that exposes kubectl functionality as tools for LLM integration.

**Key Features:**
- Readonly operations only (blocks secrets and sensitive resources)
- CLI with `serve` subcommand for stdio mode MCP communication
- Kubernetes tools: get, describe, logs, events, explain, version, cluster-info
- Built with Cobra CLI framework and mcp-go library
- Security-first design with access restrictions

**Usage:**
```bash
# Build and run
make build-kubernetes-mcp
./kubernetes-mcp serve

# Or run directly
cd kubernetes-mcp && go run . serve
```

## Workspace Development Notes

### Go Workspace Configuration
- Uses Go 1.24.4
- Multi-module workspace with `kubectl-x/` and `kubernetes-mcp/` modules
- Use `go work sync` to synchronize dependencies across workspace
- Both modules include golangci-lint as tool dependency

### Module Structure
```
.
├── go.work              # Workspace configuration
├── Makefile            # Multi-module build targets
├── kubectl-x/          # Context/namespace switching CLI
│   ├── CLAUDE.md       # Module-specific documentation
│   ├── cmd/            # CLI commands
│   ├── internal/       # Internal packages
│   └── main.go
└── kubernetes-mcp/     # MCP server for Kubernetes
    ├── cmd/            # CLI commands (root, serve)
    ├── pkg/mcp/        # MCP server implementation
    └── main.go
```

### Adding New Modules
When adding additional modules to the workspace:
1. Create new module directory with `go mod init`
2. Update `go.work` to include the new module path
3. Add build targets to root `Makefile`
4. Run `go work sync` to update workspace dependencies

### Cross-Module Dependencies
Modules are currently independent. If adding cross-module dependencies:
- Use workspace-relative imports
- Ensure proper module versioning
- Test from workspace root with `go test ./...`