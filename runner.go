package assay

import "context"

// CaseResult is the outcome of running one case: its per-metric scores and
// the averaged case score. Err is set if the agent itself failed.
type CaseResult struct {
	Case    Case
	Scores  map[string]float64 // metric name -> score
	Details map[string]string  // metric name -> human-readable detail
	Score   float64            // mean of the metric scores
	Err     error
}

// RunReport is the outcome of running a whole suite.
type RunReport struct {
	Results   []CaseResult
	Score     float64 // mean of all case scores
	Threshold float64
	Passed    bool
}

// Runner evaluates an Agent against a Suite using a fixed set of metrics.
type Runner struct {
	Metrics []Metric
}

// DefaultMetrics returns the standard v1 metric set.
func DefaultMetrics() []Metric {
	return []Metric{
		ToolCorrectness{},
		OutputContains{},
		OutputRegex{},
		MaxLatency{},
	}
}

// Run executes every case in the suite and aggregates the scores.
func (r Runner) Run(ctx context.Context, agent Agent, suite Suite) RunReport {
	metrics := r.Metrics
	if len(metrics) == 0 {
		metrics = DefaultMetrics()
	}

	report := RunReport{Threshold: suite.Threshold}
	var total float64

	for _, c := range suite.Cases {
		cr := CaseResult{
			Case:    c,
			Scores:  make(map[string]float64, len(metrics)),
			Details: make(map[string]string, len(metrics)),
		}

		result, err := agent.Run(ctx, c.Input)
		if err != nil {
			// Agent failure = the whole case scores zero, but the run continues.
			cr.Err = err
			report.Results = append(report.Results, cr)
			continue
		}

		var sum float64
		for _, m := range metrics {
			score, detail := m.Score(result, c.Expectation)
			cr.Scores[m.Name()] = score
			cr.Details[m.Name()] = detail
			sum += score
		}
		cr.Score = sum / float64(len(metrics))
		total += cr.Score
		report.Results = append(report.Results, cr)
	}

	if len(suite.Cases) > 0 {
		report.Score = total / float64(len(suite.Cases))
	}
	report.Passed = report.Score >= suite.Threshold
	return report
}
