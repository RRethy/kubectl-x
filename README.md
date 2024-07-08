# kubectl-x

Kubectl plugin with various helpers I find useful.

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
  shell       TODO.
  chat        TODO.
```

### `kubectl x ctx`

```
Switch context.

Usage:
  kubectl x ctx [context] [namespace]

Args:
  context    Partial match to filter contexts on.
             "-" to switch to the previous ctx/ns.
  namespace  Partial match to filter namespaces on.

Example:
  kubectl-pi ctx
  kubectl-pi ctx my-context
  kubectl-pi ctx my-context my-namespace
```

### `kubectl x ns`

```
Switch namespace.

Usage:
  kubectl x ns [namespace]

Args:
  namespace  Partial match to filter namespaces on.
             "-" to switch to the previous namespace.

Example:
  kubectl-pi ns
  kubectl-pi ns my-namespace
```

### `kubectl x cur`

```
Print current context and namespace.

Usage:
  kubectl x cur

Example:
  kubectl x cur
```
