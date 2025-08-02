# CLAUDE.md - Workspace Root

This file provides workspace-level guidance to Claude Code (claude.ai/code) when working with this repository.

**IMPORTANT**: When making workspace-level changes (new modules, workspace structure, build processes), update this file. For module-specific changes, update `kubectl-x/CLAUDE.md`.

## Development Commands

Use the provided Makefile for common development tasks:
```bash
make build     # Build the kubectl-x binary
make test      # Run all tests
make lint      # Run golangci-lint
make lint-fix  # Run golangci-lint with auto-fix
make fmt       # Format Go code
make vet       # Run go vet
make help      # Show all available targets
```

### Workspace-Specific Commands
```bash
# Sync workspace dependencies
go work sync

# Install from source
go install github.com/RRethy/kubectl-x@latest
```

## Architecture

### Go Workspace Structure
This is a Go workspace with a single module `kubectl-x/` containing all source code. The root contains `go.work` for workspace configuration.

For detailed module implementation, architecture, and development patterns, see `kubectl-x/CLAUDE.md`.

## Workspace Development Notes

### Go Workspace Configuration
- Uses Go 1.24.4
- Single module workspace with `kubectl-x/` module
- Use `go work sync` to synchronize dependencies across workspace

### Adding New Modules
If adding additional modules to the workspace:
1. Create new module directory
2. Update `go.work` to include the new module path
3. Run `go work sync` to update workspace dependencies

### Cross-Module Dependencies
Currently single module, but if adding modules that depend on each other:
- Use workspace-relative imports
- Ensure proper module versioning
- Test from workspace root with `go test ./...`