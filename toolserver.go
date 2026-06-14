package assay

import (
	"context"
	"time"
)

// ToolServer is the contract for evaluating a tool-exposing server (e.g. an
// MCP server) directly, with no LLM. An adapter (e.g. the mcp package) wraps a
// live server connection and satisfies this interface; the core never imports
// an MCP SDK, so it stays protocol-agnostic.
type ToolServer interface {
	// ListTools returns the names of the tools the server exposes.
	ListTools(ctx context.Context) ([]string, error)

	// CallTool invokes one tool by name with the given arguments.
	CallTool(ctx context.Context, name string, args map[string]any) (ToolOutcome, error)
}

// ToolOutcome is the normalized result of a single tool call, flattened from
// whatever the underlying protocol returns into what checks need to inspect.
type ToolOutcome struct {
	IsError bool          // the server reported the tool call as failed
	Text    string        // concatenated text content from the response
	Latency time.Duration // wall-clock time for this call
}
