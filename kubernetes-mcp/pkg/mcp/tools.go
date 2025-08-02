package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) createGetTool() mcp.Tool {
	return mcp.NewTool("get",
		mcp.WithDescription("Get Kubernetes resources from the current context/namespace using kubectl"),
		mcp.WithString("resource_type", mcp.Required(), mcp.Description("The type of Kubernetes resource to get (e.g., pods, deployments, services)")),
		mcp.WithString("resource_name", mcp.Description("Optional specific resource name to get. If not provided, lists all resources of the given type")),
		mcp.WithString("namespace", mcp.Description("Namespace to get resources from (default: current namespace)")),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithString("selector", mcp.Description("Label selector to filter results (e.g., 'app=nginx')")),
		mcp.WithString("output", mcp.Description("Output format: 'json', 'yaml', 'wide', or default table format")),
		mcp.WithBoolean("all_namespaces", mcp.Description("Get resources from all namespaces (equivalent to kubectl get --all-namespaces or -A)")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleGet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	resourceType, ok := args["resource_type"].(string)
	if !ok {
		return nil, fmt.Errorf("resource_type parameter required")
	}

	if s.isBlockedResource(resourceType) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Access to resource type '%s' is blocked for security reasons", resourceType)),
			},
		}, nil
	}

	cmdArgs := []string{"get", resourceType}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	if allNamespaces, ok := args["all_namespaces"].(bool); ok && allNamespaces {
		cmdArgs = append(cmdArgs, "--all-namespaces")
	} else if namespace, ok := args["namespace"].(string); ok && namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	if resourceName, ok := args["resource_name"].(string); ok && resourceName != "" {
		cmdArgs = append(cmdArgs, resourceName)
	}

	if selector, ok := args["selector"].(string); ok && selector != "" {
		cmdArgs = append(cmdArgs, "-l", selector)
	}

	if output, ok := args["output"].(string); ok && output != "" {
		cmdArgs = append(cmdArgs, "-o", output)
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}

func (s *Server) createDescribeTool() mcp.Tool {
	return mcp.NewTool("describe",
		mcp.WithDescription("Describe Kubernetes resources to get detailed information including events"),
		mcp.WithString("resource_type", mcp.Required(), mcp.Description("The type of Kubernetes resource to describe (e.g., pods, deployments, services)")),
		mcp.WithString("resource_name", mcp.Required(), mcp.Description("The name of the resource to describe")),
		mcp.WithString("namespace", mcp.Description("Namespace of the resource (default: current namespace)")),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleDescribe(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	resourceType, ok := args["resource_type"].(string)
	if !ok {
		return nil, fmt.Errorf("resource_type parameter required")
	}

	if s.isBlockedResource(resourceType) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Access to resource type '%s' is blocked for security reasons", resourceType)),
			},
		}, nil
	}

	resourceName, ok := args["resource_name"].(string)
	if !ok {
		return nil, fmt.Errorf("resource_name parameter required")
	}

	cmdArgs := []string{"describe", resourceType, resourceName}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	if namespace, ok := args["namespace"].(string); ok && namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}

func (s *Server) createLogsTool() mcp.Tool {
	return mcp.NewTool("logs",
		mcp.WithDescription("Get logs from a pod container"),
		mcp.WithString("pod_name", mcp.Required(), mcp.Description("The name of the pod")),
		mcp.WithString("namespace", mcp.Description("Namespace of the pod (default: current namespace)")),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithString("container", mcp.Description("Container name (if multiple containers in pod)")),
		mcp.WithNumber("tail", mcp.Description("Number of lines to show from the end of the logs (default: 100)")),
		mcp.WithString("since", mcp.Description("Show logs since this duration (e.g., 5m, 1h)")),
		mcp.WithBoolean("previous", mcp.Description("Show logs from previous instance of container")),
		mcp.WithBoolean("timestamps", mcp.Description("Include timestamps in log output")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	podName, ok := args["pod_name"].(string)
	if !ok {
		return nil, fmt.Errorf("pod_name parameter required")
	}

	cmdArgs := []string{"logs", podName}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	if namespace, ok := args["namespace"].(string); ok && namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	if container, ok := args["container"].(string); ok && container != "" {
		cmdArgs = append(cmdArgs, "-c", container)
	}

	tail := 100.0
	if tailVal, ok := args["tail"].(float64); ok {
		tail = tailVal
	}
	if tail > 0 {
		cmdArgs = append(cmdArgs, "--tail", fmt.Sprintf("%d", int(tail)))
	}

	if since, ok := args["since"].(string); ok && since != "" {
		cmdArgs = append(cmdArgs, "--since", since)
	}

	if previous, ok := args["previous"].(bool); ok && previous {
		cmdArgs = append(cmdArgs, "--previous")
	}

	if timestamps, ok := args["timestamps"].(bool); ok && timestamps {
		cmdArgs = append(cmdArgs, "--timestamps")
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}

func (s *Server) createEventsTool() mcp.Tool {
	return mcp.NewTool("events",
		mcp.WithDescription("Get events from Kubernetes with filtering options"),
		mcp.WithString("namespace", mcp.Description("Namespace to get events from (default: current namespace)")),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithString("for", mcp.Description("Filter events for a specific resource (e.g., pod/my-pod, deployment/my-deployment)")),
		mcp.WithBoolean("all_namespaces", mcp.Description("Get events from all namespaces (equivalent to kubectl get events --all-namespaces or -A)")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleEvents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	cmdArgs := []string{"get", "events", "--sort-by=.lastTimestamp"}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	if allNamespaces, ok := args["all_namespaces"].(bool); ok && allNamespaces {
		cmdArgs = append(cmdArgs, "--all-namespaces")
	} else if namespace, ok := args["namespace"].(string); ok && namespace != "" {
		cmdArgs = append(cmdArgs, "-n", namespace)
	}

	if forResource, ok := args["for"].(string); ok && forResource != "" {
		cmdArgs = append(cmdArgs, "--field-selector", fmt.Sprintf("involvedObject.name=%s", forResource))
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}

func (s *Server) createExplainTool() mcp.Tool {
	return mcp.NewTool("explain",
		mcp.WithDescription("Explain Kubernetes resource types and their fields"),
		mcp.WithString("resource", mcp.Required(), mcp.Description("Resource type to explain (e.g., pods, deployments, pods.spec.containers)")),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithBoolean("recursive", mcp.Description("Print all fields recursively")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleExplain(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	resource, ok := args["resource"].(string)
	if !ok {
		return nil, fmt.Errorf("resource parameter required")
	}

	cmdArgs := []string{"explain", resource}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	if recursive, ok := args["recursive"].(bool); ok && recursive {
		cmdArgs = append(cmdArgs, "--recursive")
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}

func (s *Server) createVersionTool() mcp.Tool {
	return mcp.NewTool("version",
		mcp.WithDescription("Get version information for kubectl client and Kubernetes cluster"),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithString("output", mcp.Description("Output format: 'json', 'yaml', or default")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	cmdArgs := []string{"version"}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	if output, ok := args["output"].(string); ok && output != "" {
		cmdArgs = append(cmdArgs, "-o", output)
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}

func (s *Server) createClusterInfoTool() mcp.Tool {
	return mcp.NewTool("cluster-info",
		mcp.WithDescription("Get cluster information including master and services locations"),
		mcp.WithString("context", mcp.Description("Kubernetes context to use (default: current context)")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

func (s *Server) handleClusterInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	cmdArgs := []string{"cluster-info"}

	if contextName, ok := args["context"].(string); ok && contextName != "" {
		cmdArgs = append([]string{"--context", contextName}, cmdArgs...)
	}

	stdout, stderr, err := s.runKubectl(ctx, cmdArgs...)
	return s.formatOutput(stdout, stderr, err)
}
