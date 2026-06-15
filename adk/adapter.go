// Package adk adapts a Google ADK agent so it can be evaluated by assay.
package adk

import (
	"context"
	"strings"
	"time"

	"github.com/tushariitr-19/assay"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

// Adapter wraps an ADK agent so it satisfies assay.Agent.
type Adapter struct {
	runner *runner.Runner
}

// New builds an assay-compatible adapter around any ADK agent.
func New(a agent.Agent) (*Adapter, error) {
	r, err := runner.New(runner.Config{
		AppName:           "assay",
		Agent:             a,
		SessionService:    session.InMemoryService(),
		AutoCreateSession: true,
	})
	if err != nil {
		return nil, err
	}
	return &Adapter{runner: r}, nil
}

// Run drives the agent for one input and returns an assay.Result.
func (ad *Adapter) Run(ctx context.Context, input string) (assay.Result, error) {
	start := time.Now()
	msg := genai.NewContentFromText(input, genai.RoleUser)

	var calls []assay.ToolCall
	var text strings.Builder

	for ev, err := range ad.runner.Run(ctx, "assay", "assay-session", msg, agent.RunConfig{}) {
		if err != nil {
			return assay.Result{}, err
		}
		if ev == nil || ev.Content == nil {
			continue
		}
		for _, part := range ev.Content.Parts {
			if part.FunctionCall != nil {
				calls = append(calls, assay.ToolCall{Name: part.FunctionCall.Name, Args: part.FunctionCall.Args})
			}
			if part.Text != "" {
				text.WriteString(part.Text)
			}
		}
	}

	return assay.Result{Output: text.String(), ToolCalls: calls, Latency: time.Since(start)}, nil
}
