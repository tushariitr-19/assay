package config

// Config holds the settings that control how an eval run executes and reports.
// It is populated by the CLI (from flags) and passed down to the runner and
// reporter. The core evaluation types do not depend on this — it lives at the
// edges, like the logger.
type Config struct {
	// SuitePath is the path to the YAML suite file to run.
	SuitePath string

	// JSONOut, if set, is the path where a JSON report is written.
	JSONOut string

	// Debug enables debug-level logging.
	Debug bool

	// Threshold overrides the suite's own pass threshold when >= 0.
	// A value of -1 means "use whatever the suite file specifies."
	Threshold float64
}

// Default returns a Config with sensible starting values.
func Default() *Config {
	return &Config{
		Threshold: -1, // sentinel: defer to the suite file's threshold
	}
}
