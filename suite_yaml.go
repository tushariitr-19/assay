package assay

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// yamlSuite/yamlCase mirror the on-disk YAML shape. We decode into these,
// then convert to the in-memory Suite — so the file format can evolve
// independently of the core types.
type yamlSuite struct {
	Threshold float64    `yaml:"threshold"`
	Cases     []yamlCase `yaml:"cases"`
}

type yamlCase struct {
	Name           string   `yaml:"name"`
	Input          string   `yaml:"input"`
	ExpectTools    []string `yaml:"expect_tools"`
	OutputContains []string `yaml:"output_contains"`
	OutputRegex    []string `yaml:"output_regex"`
	MaxLatency     string   `yaml:"max_latency"` // duration string e.g. "8s"
}

// LoadSuite reads and parses a YAML suite file into a Suite.
func LoadSuite(path string) (Suite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Suite{}, fmt.Errorf("reading suite: %w", err)
	}

	var ys yamlSuite
	if err := yaml.Unmarshal(data, &ys); err != nil {
		return Suite{}, fmt.Errorf("parsing suite: %w", err)
	}

	suite := Suite{Threshold: ys.Threshold}
	for i, yc := range ys.Cases {
		exp := Expectation{
			Tools:        yc.ExpectTools,
			OutputSubstr: yc.OutputContains,
			OutputRegex:  yc.OutputRegex,
		}
		if yc.MaxLatency != "" {
			d, err := time.ParseDuration(yc.MaxLatency)
			if err != nil {
				return Suite{}, fmt.Errorf("case %d (%s): invalid max_latency %q: %w",
					i, yc.Name, yc.MaxLatency, err)
			}
			exp.MaxLatencyMS = d.Milliseconds()
		}
		suite.Cases = append(suite.Cases, Case{
			Name:        yc.Name,
			Input:       yc.Input,
			Expectation: exp,
		})
	}
	return suite, nil
}
