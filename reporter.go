package assay

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

func mark(score float64) string {
	switch {
	case score >= 1.0:
		return "✓"
	case score <= 0.0:
		return "✗"
	default:
		return "~"
	}
}

// WriteText renders the report as a human-readable table to w.
func (r RunReport) WriteText(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	fmt.Fprintln(tw, "CASE\tSCORE\tRESULT")

	for _, cr := range r.Results {
		if cr.Err != nil {
			fmt.Fprintf(tw, "%s\t0.00\tERROR: %v\n", cr.Case.Name, cr.Err)
			continue
		}
		fmt.Fprintf(tw, "%s\t%.2f\t%s\n", cr.Case.Name, cr.Score, mark(cr.Score))

		// Stable, alphabetical order so output is deterministic across runs.
		names := make([]string, 0, len(cr.Scores))
		for name := range cr.Scores {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Fprintf(tw, "  %s\t%.2f\t%s\n", name, cr.Scores[name], cr.Details[name])
		}
	}
	tw.Flush()

	status := "PASS"
	if !r.Passed {
		status = "FAIL"
	}
	fmt.Fprintf(w, "\nSUITE SCORE: %.2f   threshold %.2f   %s\n", r.Score, r.Threshold, status)
}

// jsonReport is the serializable shape of a run. We map to it explicitly so
// the JSON output stays stable even if internal types are refactored.
type jsonReport struct {
	Score     float64          `json:"score"`
	Threshold float64          `json:"threshold"`
	Passed    bool             `json:"passed"`
	Cases     []jsonCaseResult `json:"cases"`
}

type jsonCaseResult struct {
	Name   string             `json:"name"`
	Score  float64            `json:"score"`
	Scores map[string]float64 `json:"scores"`
	Error  string             `json:"error,omitempty"`
}

// WriteJSON renders the report as indented JSON to w.
func (r RunReport) WriteJSON(w io.Writer) error {
	out := jsonReport{
		Score:     r.Score,
		Threshold: r.Threshold,
		Passed:    r.Passed,
	}
	for _, cr := range r.Results {
		jc := jsonCaseResult{
			Name:   cr.Case.Name,
			Score:  cr.Score,
			Scores: cr.Scores,
		}
		if cr.Err != nil {
			jc.Error = cr.Err.Error()
		}
		out.Cases = append(out.Cases, jc)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
