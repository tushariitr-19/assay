// Package mcp adapts a Model Context Protocol server so it can be evaluated by
// assay. It wraps an MCP client session and satisfies assay.ToolServer.
package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/tushariitr-19/assay"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server is an assay.ToolServer backed by a live MCP client session.
type Server struct {
	session *mcpsdk.ClientSession
}

// Connect launches an MCP server via the given command (stdio transport),
// performs the handshake, and returns a ready-to-evaluate Server.
// Call Close when done.
func Connect(ctx context.Context, command string, args ...string) (*Server, error) {
	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "assay", Version: "v0.1.0"}, nil)

	transport := &mcpsdk.CommandTransport{Command: exec.Command(command, args...)}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to MCP server: %w", err)
	}
	return &Server{session: session}, nil
}

// Close shuts down the session and the underlying server process.
func (s *Server) Close() error {
	return s.session.Close()
}

// ListTools satisfies assay.ToolServer.
func (s *Server) ListTools(ctx context.Context) ([]string, error) {
	res, err := s.session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(res.Tools))
	for _, t := range res.Tools {
		names = append(names, t.Name)
	}
	return names, nil
}

// CallTool satisfies assay.ToolServer, flattening the MCP result into a ToolOutcome.
func (s *Server) CallTool(ctx context.Context, name string, args map[string]any) (assay.ToolOutcome, error) {
	start := time.Now()
	res, err := s.session.CallTool(ctx, &mcpsdk.CallToolParams{Name: name, Arguments: args})
	latency := time.Since(start)
	if err != nil {
		return assay.ToolOutcome{}, err
	}

	var text strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*mcpsdk.TextContent); ok {
			text.WriteString(tc.Text)
		}
	}

	return assay.ToolOutcome{
		IsError: res.IsError,
		Text:    text.String(),
		Latency: latency,
	}, nil
}
