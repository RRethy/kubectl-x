package mcp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve() error {
	mcpServer := server.NewMCPServer(
		"kubernetes-mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithPromptCapabilities(false),
	)

	mcpServer.AddTools(
		server.ServerTool{Tool: s.createGetTool(), Handler: s.handleGet},
		server.ServerTool{Tool: s.createDescribeTool(), Handler: s.handleDescribe},
		server.ServerTool{Tool: s.createLogsTool(), Handler: s.handleLogs},
		server.ServerTool{Tool: s.createEventsTool(), Handler: s.handleEvents},
		server.ServerTool{Tool: s.createExplainTool(), Handler: s.handleExplain},
		server.ServerTool{Tool: s.createVersionTool(), Handler: s.handleVersion},
		server.ServerTool{Tool: s.createClusterInfoTool(), Handler: s.handleClusterInfo},
	)

	return server.ServeStdio(mcpServer)
}

func (s *Server) isBlockedResource(resourceType string) bool {
	blocked := []string{"secret", "secrets"}
	lower := strings.ToLower(resourceType)
	for _, b := range blocked {
		if lower == b {
			return true
		}
	}
	return false
}

func (s *Server) runKubectl(ctx context.Context, args ...string) (string, string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func (s *Server) formatOutput(stdout, stderr string, err error) (*mcp.CallToolResult, error) {
	if err != nil {
		errorMsg := fmt.Sprintf("kubectl command failed: %v", err)
		if stderr != "" {
			errorMsg = fmt.Sprintf("%s\nkubectl error: %s", errorMsg, stderr)
		}

		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(errorMsg),
			},
		}, nil
	}

	output := stdout
	if stderr != "" {
		output = fmt.Sprintf("%s\n\nWarnings:\n%s", stdout, stderr)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(output),
		},
	}, nil
}
