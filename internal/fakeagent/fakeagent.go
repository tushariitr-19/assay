// Package fakeagent provides an in-memory Agent implementation for tests.
// It performs no network calls and needs no credentials: you construct it
// with the exact Result it should return, so tests stay deterministic.
package fakeagent

import (
	"context"

	"github.com/tushariitr-19/assay"
)

// Agent is a canned implementation of assay.Agent. It returns the configured
// Result (and optional error) regardless of input, recording the input it saw.
type Agent struct {
	// Result is what Run returns.
	Result assay.Result

	// Err, if set, is returned by Run instead of Result.
	Err error

	// LastInput captures the most recent input passed to Run (handy in tests).
	LastInput string
}

// New builds a fake agent that always returns the given result.
func New(result assay.Result) *Agent {
	return &Agent{Result: result}
}

// Run satisfies the assay.Agent interface.
func (a *Agent) Run(_ context.Context, input string) (assay.Result, error) {
	a.LastInput = input
	if a.Err != nil {
		return assay.Result{}, a.Err
	}
	return a.Result, nil
}
