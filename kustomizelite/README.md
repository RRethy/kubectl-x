# kustomizelite

A lightweight CLI tool for working with Kustomize configurations.

## Installation

```bash
go install github.com/RRethy/utils/kustomizelite@latest
```

Or build from source:

```bash
go build -o kustomizelite .
```

## Usage

### Build Command

Validate kustomization.yaml files:

```bash
# Validate a single file
adamize build /path/to/kustomization.yaml

# Validate multiple files
adamize build ./base/kustomization.yaml ./overlays/prod/kustomization.yaml

# Show file contents with debug flag
adamize build -v /path/to/kustomization.yaml
```

The build command validates that:
- The path points to a file (not a directory)
- The file is named exactly "kustomization.yaml"
- The file exists and is readable

## Development

### Project Structure

```
adamize/
├── api/v1/           # Kustomization data structures
├── cmd/              # CLI commands
├── pkg/
│   ├── kustomize/    # Core business logic
│   └── cli/          # CLI presentation layer
└── main.go           # Entry point
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o adamize .
```

## Contributing

This tool follows standard Go conventions and uses Cobra for CLI structure.