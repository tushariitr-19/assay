package assay

import "context"

// ServerSuite is a named set of checks with a pass threshold in [0,1].
type ServerSuite struct {
	Threshold float64
	Checks    []ServerCheck
}

// RunServer executes every check against the ToolServer and aggregates scores
// into the same RunReport shape the agent runner produces — so the existing
// reporters (WriteText / WriteJSON) work unchanged.
func RunServer(ctx context.Context, ts ToolServer, suite ServerSuite) RunReport {
	report := RunReport{Threshold: suite.Threshold}
	var total float64

	for _, check := range suite.Checks {
		score, detail := check.Run(ctx, ts)
		cr := CaseResult{
			Case:    Case{Name: check.Name()},
			Scores:  map[string]float64{check.Name(): score},
			Details: map[string]string{check.Name(): detail},
			Score:   score,
		}
		total += score
		report.Results = append(report.Results, cr)
	}

	if len(suite.Checks) > 0 {
		report.Score = total / float64(len(suite.Checks))
	}
	report.Passed = report.Score >= suite.Threshold
	return report
}
