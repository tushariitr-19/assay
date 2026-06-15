package main

import (
	"fmt"
	"os"

	"github.com/tushariitr-19/assay/adk"
	"github.com/tushariitr-19/assay/cli/app"

	"github.com/tushariitr-19/research-agent-adk/agents"
	"github.com/tushariitr-19/research-agent-adk/config"
	"github.com/tushariitr-19/research-agent-adk/logger"
)

func main() {
	// Build the agent (this is your agent's own startup — config + logger).
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	if err := logger.Init(cfg.Debug); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	defer logger.Sync()

	rootAgent, err := agents.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building agent: %v\n", err)
		os.Exit(2)
	}

	// The integration: wrap, register, hand off to the assay CLI.
	adapter, err := adk.New(rootAgent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	app.SetAgent(adapter)
	app.Main()
}
