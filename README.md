# kubectl-x

Kubectl plugin that provides convenient context and namespace switching utilities for Kubernetes. Features interactive fuzzy search capabilities and maintains command history for quick navigation.

## Quickstart

```bash
go install github.com/RRethy/kubectl-x@latest
```

## Docs

```
kubectl (kube-control) plugin with various useful extensions.

Usage:
  kubectl x [flags]
  kubectl x [command]

Available Commands:
  ctx         Switch context.
  cur         Print current context and namespace.
  ns          Switch namespace.
```

### `kubectl x ctx`

```
Switch context with interactive fuzzy search.

Usage:
  kubectl x ctx [context] [namespace]

Args:
  context    Partial match to filter contexts on.
             "-" to switch to the previous ctx/ns.
             If no args, opens interactive fuzzy finder.
  namespace  Partial match to filter namespaces on.

Example:
  kubectl x ctx                        # Interactive context selection
  kubectl x ctx my-context             # Switch to context with partial match
  kubectl x ctx my-context my-namespace # Switch context and namespace
  kubectl x ctx -                      # Switch to previous context/namespace
```

### `kubectl x ns`

```
Switch namespace with interactive fuzzy search.

Usage:
  kubectl x ns [namespace]

Args:
  namespace  Partial match to filter namespaces on.
             "-" to switch to the previous namespace.
             If no args, opens interactive fuzzy finder.

Example:
  kubectl x ns                # Interactive namespace selection
  kubectl x ns my-namespace   # Switch to namespace with partial match
  kubectl x ns -              # Switch to previous namespace
```

### `kubectl x cur`

```
Print current context and namespace.

Usage:
  kubectl x cur

Example:
  kubectl x cur
```
