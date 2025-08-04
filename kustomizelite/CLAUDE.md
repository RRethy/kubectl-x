# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in the kustomizelite module.

**IMPORTANT: When making ANY changes to the kustomizelite codebase, you MUST update this document if the changes affect:**
- Module structure or architecture
- Command definitions or CLI interface
- Data structures in api/v1/
- Business logic interfaces or implementations
- Testing patterns or development guidelines
- Error handling conventions

Always keep this documentation current and accurate to ensure effective collaboration.

## Module Overview

kustomizelite is a lightweight CLI tool for working with Kustomize configurations. It provides commands to validate and process kustomization.yaml files, including support for Helm chart inflation.

## Architecture

### Project Structure
- `api/v1/` - Kustomization data structures (Kustomization struct and related types)
- `cmd/` - Cobra command definitions (root.go, build.go)
- `pkg/kustomize/` - Core business logic with Kustomizer interface
- `pkg/helm/` - Helm chart inflation functionality
- `pkg/cli/build/` - CLI presentation layer for the build command
- `pkg/exec/` - Command execution wrapper with environment variable support
- `main.go` - Entry point

### Key Design Patterns

#### Interface-Based Design
- `pkg/kustomize/Kustomizer` - Interface for kustomization operations
- Concrete implementation is private (`kustomization` struct in `kustomization.go`)
- Constructor returns interface type (`NewKustomize`)
- Handles both regular Kustomizations and Components through the same struct

#### Dependency Injection
- Business logic (`Kustomizer`) is injected into CLI layer
- Constructed in `pkg/cli/build/build.go` and passed to `Builder`

#### Testing Strategy
- Unit tests for business logic in `pkg/kustomize/`
- `FakeKustomizer` provided for testing CLI layer
- CLI tests use fake to isolate from file system

### Current Commands

#### build
- Processes kustomization.yaml files and outputs rendered Kubernetes resources
- Supports all standard Kustomize features: resources, patches, transformers
- Supports Helm chart inflation via helmCharts field
- Outputs raw YAML suitable for piping to kubectl apply
- Accepts a single path argument (defaults to current directory)
- Supports batch builds via `-f`/`--file` flag for parallel processing

## Development Guidelines

### Adding New Commands
1. Create command definition in `cmd/`
2. Create business logic interface in `pkg/[feature]/`
3. Create CLI layer in `pkg/cli/[command]/`
4. Wire together in `pkg/cli/[command]/[command].go`

### Testing
- Always create fake implementations for interfaces
- Test business logic separately from CLI presentation
- Use table-driven tests

### Code Style
- **DO NOT add comments unless explicitly asked by the user**
- Return interfaces from constructors
- Keep business logic separate from CLI concerns
- Only add comments when the user specifically requests them

### Error Handling
When wrapping errors, use descriptive messages that explain the action being performed, not that it failed.

**Preferred:**
```go
if err := yaml.Unmarshal(content, &kustomizationMap); err != nil {
    return nil, fmt.Errorf("parsing YAML into map: %w", err)
}
```

**Avoid:**
```go
if err := yaml.Unmarshal(content, &kustomizationMap); err != nil {
    return nil, fmt.Errorf("failed to parse YAML into map: %w", err)
}
```

The error itself already indicates failure, so including "failed to" is redundant.

### Helm Support

kustomizelite supports Helm chart inflation through the `helmCharts` field in kustomization.yaml:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

helmGlobals:
  chartHome: ./charts

helmCharts:
- name: nginx
  version: 15.14.0
  releaseName: my-nginx
  namespace: web
  valuesFile: values.yaml
  additionalValuesFiles:
  - values-prod.yaml
  valuesInline:
    replicaCount: 3
  includeCRDs: true
```

#### Environment Variables
- `HELM_BINARY_PATH`: Override the helm binary path (defaults to `helm` in PATH)
  - This is read at the kustomize package level when creating a new Kustomizer instance

#### Requirements
- Helm CLI must be installed and available
- Charts must be available locally (either in chartHome or relative to kustomization.yaml)

#### Implementation Details
- Uses `helm template` command to generate YAML
- Parses multi-document YAML output
- Applies Kustomize transformations after chart inflation
- Temporary files are used for inline values and cleaned up automatically
- Helm functionality is in separate `pkg/helm` package with `Templater` interface
- Processing order: Resources → Helm Charts → Components → Transformations

### Batch Build Support

The build command supports batch processing via the `-f`/`--file` flag:

```bash
kustomizelite build -f batch.yaml
```

Batch file format:
```yaml
apiVersion: kustomizelite.io/v1
kind: BatchBuild
env:
  - name: BUNDLE_GEMFILE
    value: /path/to/Gemfile
  - name: HELM_BINARY_PATH
    value: /path/to/helm-wrapper
builds:
  - kustomization: /path/to/app1/kustomization.yaml
    output: /path/to/app1/build/output.yaml
    env:
      - name: TARGET_PATH
        value: /path/to/app1/target.yaml
  - kustomization: /path/to/app2/kustomization.yaml
    output: /path/to/app2/build/output.yaml
    env:
      - name: TARGET_PATH
        value: /path/to/app2/target.yaml
```

#### Batch Build Features
- **Parallel Execution**: All builds in the batch are executed concurrently
- **Environment Variables**: 
  - Global env vars are set for all builds
  - Build-specific env vars override global ones
  - Environment variables are passed to helm and other subcommands
- **Output Management**: Each build writes to its specified output file
- **Error Handling**: Batch continues processing all builds even if some fail

#### Implementation Details
- Uses `pkg/exec` package for environment-aware command execution
- Helm templater receives custom exec wrapper with environment variables
- Builds run in separate goroutines with synchronization via WaitGroup