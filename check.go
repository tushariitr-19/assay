package assay

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ServerCheck scores one aspect of a ToolServer. Each check performs its own
// operation against the server (list or call) and returns a score in [0,1].
type ServerCheck interface {
	Name() string
	Run(ctx context.Context, ts ToolServer) (score float64, detail string)
}

// ToolsListCheck verifies the server exposes the expected tools.
type ToolsListCheck struct {
	Expected []string
}

func (ToolsListCheck) Name() string { return "tools_list" }

func (c ToolsListCheck) Run(ctx context.Context, ts ToolServer) (float64, string) {
	listed, err := ts.ListTools(ctx)
	if err != nil {
		return 0, fmt.Sprintf("list failed: %v", err)
	}
	have := make(map[string]bool, len(listed))
	for _, t := range listed {
		have[t] = true
	}
	found := 0
	for _, want := range c.Expected {
		if have[want] {
			found++
		}
	}
	if len(c.Expected) == 0 {
		return 1, "no tools expected"
	}
	score := float64(found) / float64(len(c.Expected))
	return score, fmt.Sprintf("%d/%d expected tools exposed", found, len(c.Expected))
}

// ToolCallCheck calls one tool and asserts on the outcome.
type ToolCallCheck struct {
	Tool           string
	Args           map[string]any
	ExpectNoError  bool     // outcome.IsError must be false
	ResultContains []string // substrings expected in outcome.Text
	MaxLatencyMS   int64    // 0 = no latency check
}

func (c ToolCallCheck) Name() string { return "tool_call:" + c.Tool }

func (c ToolCallCheck) Run(ctx context.Context, ts ToolServer) (float64, string) {
	out, err := ts.CallTool(ctx, c.Tool, c.Args)
	if err != nil {
		return 0, fmt.Sprintf("call failed: %v", err)
	}

	var checks, passed int

	if c.ExpectNoError {
		checks++
		if !out.IsError {
			passed++
		}
	}
	for _, sub := range c.ResultContains {
		checks++
		if strings.Contains(out.Text, sub) {
			passed++
		}
	}
	if c.MaxLatencyMS > 0 {
		checks++
		if out.Latency <= time.Duration(c.MaxLatencyMS)*time.Millisecond {
			passed++
		}
	}

	if checks == 0 {
		return 1, "no assertions"
	}
	return float64(passed) / float64(checks), fmt.Sprintf("%d/%d assertions passed", passed, checks)
}
