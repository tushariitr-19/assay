package assay

import (
	"fmt"
	"regexp"
	"strings"
)

// Metric scores one aspect of a Result against a case's expectations.
// Every metric — including future ones like an LLM judge — satisfies this
// same interface, which is what keeps the runner agnostic to what it's scoring.
type Metric interface {
	// Name identifies the metric in reports.
	Name() string

	// Score returns a value in [0,1] plus a short human-readable detail.
	Score(result Result, exp Expectation) (score float64, detail string)
}

// Expectation describes what a single case expects. Fields left empty are
// simply not checked by the metrics that read them.
type Expectation struct {
	Tools        []string // expected tool names (for ToolCorrectness)
	OutputSubstr []string // substrings expected in Output
	OutputRegex  []string // regex patterns Output should match
	MaxLatencyMS int64    // latency ceiling in milliseconds (0 = no check)
}

// ToolCorrectness scores whether the agent called the expected tools.
// In its simplest form it checks set membership (did each expected tool
// get called, ignoring order and extras).
type ToolCorrectness struct {
	Ordered bool // if true, call order must match expectation order
}

func (ToolCorrectness) Name() string { return "tool_correctness" }

func (m ToolCorrectness) Score(result Result, exp Expectation) (float64, string) {
	if len(exp.Tools) == 0 {
		return 1.0, "no tools expected"
	}

	called := make([]string, len(result.ToolCalls))
	for i, tc := range result.ToolCalls {
		called[i] = tc.Name
	}

	if m.Ordered {
		matched := 0
		for i, want := range exp.Tools {
			if i < len(called) && called[i] == want {
				matched++
			}
		}
		score := float64(matched) / float64(len(exp.Tools))
		return score, fmt.Sprintf("ordered: %d/%d matched", matched, len(exp.Tools))
	}

	// Unordered: how many expected tools appear anywhere in the calls.
	calledSet := make(map[string]bool, len(called))
	for _, c := range called {
		calledSet[c] = true
	}
	found := 0
	for _, want := range exp.Tools {
		if calledSet[want] {
			found++
		}
	}
	score := float64(found) / float64(len(exp.Tools))
	return score, fmt.Sprintf("%d/%d expected tools called", found, len(exp.Tools))
}

// OutputContains scores how many expected substrings appear in Output.
type OutputContains struct{}

func (OutputContains) Name() string { return "output_contains" }

func (OutputContains) Score(result Result, exp Expectation) (float64, string) {
	if len(exp.OutputSubstr) == 0 {
		return 1.0, "no substrings expected"
	}
	found := 0
	for _, want := range exp.OutputSubstr {
		if strings.Contains(result.Output, want) {
			found++
		}
	}
	score := float64(found) / float64(len(exp.OutputSubstr))
	return score, fmt.Sprintf("%d/%d substrings present", found, len(exp.OutputSubstr))
}

// OutputRegex scores how many expected patterns match Output.
type OutputRegex struct{}

func (OutputRegex) Name() string { return "output_regex" }

func (OutputRegex) Score(result Result, exp Expectation) (float64, string) {
	if len(exp.OutputRegex) == 0 {
		return 1.0, "no patterns expected"
	}
	matched := 0
	for _, pat := range exp.OutputRegex {
		re, err := regexp.Compile(pat)
		if err != nil {
			return 0, fmt.Sprintf("invalid pattern %q: %v", pat, err)
		}
		if re.MatchString(result.Output) {
			matched++
		}
	}
	score := float64(matched) / float64(len(exp.OutputRegex))
	return score, fmt.Sprintf("%d/%d patterns matched", matched, len(exp.OutputRegex))
}

// MaxLatency is a pass/fail check on run duration.
type MaxLatency struct{}

func (MaxLatency) Name() string { return "max_latency" }

func (MaxLatency) Score(result Result, exp Expectation) (float64, string) {
	if exp.MaxLatencyMS == 0 {
		return 1.0, "no latency ceiling"
	}
	actualMS := result.Latency.Milliseconds()
	if actualMS <= exp.MaxLatencyMS {
		return 1.0, fmt.Sprintf("%dms <= %dms", actualMS, exp.MaxLatencyMS)
	}
	return 0.0, fmt.Sprintf("%dms > %dms (too slow)", actualMS, exp.MaxLatencyMS)
}
