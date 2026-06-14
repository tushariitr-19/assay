package assay

import (
	"context"
	"time"
)

// Agent is the only contract the assay core depends on. Any agent — ADK, MCP,
// or anything else — is evaluated by wrapping it in an adapter that satisfies
// this interface. The core never imports a framework; it only knows this.
type Agent interface {
	Run(ctx context.Context, input string) (Result, error)
}

// Result is what an agent returns from a single run: its final text output,
// the trace of tools it called, and how long the run took.
type Result struct {
	Output    string
	ToolCalls []ToolCall
	Latency   time.Duration
}

// ToolCall records a single tool invocation: the tool's name and the
// arguments the agent chose to pass it.
type ToolCall struct {
	Name string
	Args map[string]any
}
