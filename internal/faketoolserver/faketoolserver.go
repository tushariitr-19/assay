// Package faketoolserver provides an in-memory ToolServer for tests.
// No network, no MCP SDK: you configure exactly what it lists and returns.
package faketoolserver

import (
	"context"

	"github.com/tushariitr-19/assay"
)

// Server is a canned assay.ToolServer for deterministic tests.
type Server struct {
	Tools    []string                     // what ListTools returns
	Outcomes map[string]assay.ToolOutcome // tool name -> canned CallTool result
	ListErr  error                        // if set, ListTools returns this
	CallErr  error                        // if set, CallTool returns this
}

// ListTools returns the configured tool names.
func (s *Server) ListTools(_ context.Context) ([]string, error) {
	if s.ListErr != nil {
		return nil, s.ListErr
	}
	return s.Tools, nil
}

// CallTool returns the canned outcome for the named tool.
func (s *Server) CallTool(_ context.Context, name string, _ map[string]any) (assay.ToolOutcome, error) {
	if s.CallErr != nil {
		return assay.ToolOutcome{}, s.CallErr
	}
	return s.Outcomes[name], nil
}
