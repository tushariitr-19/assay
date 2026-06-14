package assay

// Case is a single evaluation scenario: an input to send the agent and the
// expectations its result is scored against.
type Case struct {
	Name        string
	Input       string
	Expectation Expectation
}

// Suite is a named collection of cases with a pass threshold in [0,1].
// A run passes when its aggregate score meets or exceeds Threshold.
type Suite struct {
	Threshold float64
	Cases     []Case
}
