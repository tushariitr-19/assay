// Package app is the importable entry point for the assay CLI. Users who want
// to evaluate their own agent import this package, call SetAgent with their
// wrapped agent, and then Main(). The default assay binary calls Main() with no
// agent set (MCP-only).
package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tushariitr-19/assay"
	"github.com/tushariitr-19/assay/mcp"
)

// registered holds the single agent set by SetAgent, if any.
var registered assay.Agent

// SetAgent registers the agent that `assay agent` will evaluate. Call it before
// Main() from your own main(), after wrapping your agent in an adapter.
func SetAgent(a assay.Agent) {
	registered = a
}

// Main runs the assay CLI and exits the process.
func Main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) < 2 {
		usage()
		return 2
	}
	switch os.Args[1] {
	case "mcp":
		return runMCP(os.Args[2:])
	case "agent":
		return runAgent(os.Args[2:])
	default:
		usage()
		return 2
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  assay mcp   --server \"<command>\" --suite <checks.yaml>")
	fmt.Fprintln(os.Stderr, "  assay agent --suite <cases.yaml>   (agent must be set via SetAgent)")
}

func suiteFlag(argv []string) string {
	for i := 0; i < len(argv); i++ {
		if argv[i] == "--suite" && i+1 < len(argv) {
			return argv[i+1]
		}
	}
	return ""
}

func runMCP(argv []string) int {
	var serverCmd string
	for i := 0; i < len(argv); i++ {
		if argv[i] == "--server" && i+1 < len(argv) {
			serverCmd = argv[i+1]
			i++
		}
	}
	suitePath := suiteFlag(argv)
	if serverCmd == "" || suitePath == "" {
		usage()
		return 2
	}

	suite, err := assay.LoadServerSuite(suitePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	parts := strings.Fields(serverCmd)
	ctx := context.Background()

	server, err := mcp.Connect(ctx, parts[0], parts[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting: %v\n", err)
		return 2
	}
	defer server.Close()

	report := assay.RunServer(ctx, server, suite)
	report.WriteText(os.Stdout)
	if !report.Passed {
		return 1
	}
	return 0
}

func runAgent(argv []string) int {
	if registered == nil {
		fmt.Fprintln(os.Stderr, "error: no agent set — call app.SetAgent(adapter) before Main()")
		fmt.Fprintln(os.Stderr, "see examples/research-agent for the pattern")
		return 2
	}
	suitePath := suiteFlag(argv)
	if suitePath == "" {
		usage()
		return 2
	}

	suite, err := assay.LoadSuite(suitePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	report := assay.Runner{}.Run(context.Background(), registered, suite)
	report.WriteText(os.Stdout)
	if !report.Passed {
		return 1
	}
	return 0
}
