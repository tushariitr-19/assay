package assay

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type yamlServerSuite struct {
	Threshold float64           `yaml:"threshold"`
	Checks    []yamlServerCheck `yaml:"checks"`
}

type yamlServerCheck struct {
	Name           string         `yaml:"name"`
	Type           string         `yaml:"type"` // "tools_list" or "tool_call"
	ExpectTools    []string       `yaml:"expect_tools"`
	Tool           string         `yaml:"tool"`
	Args           map[string]any `yaml:"args"`
	ExpectNoError  bool           `yaml:"expect_no_error"`
	ResultContains []string       `yaml:"result_contains"`
	MaxLatency     string         `yaml:"max_latency"`
}

// LoadServerSuite reads a YAML file into a ServerSuite of checks.
func LoadServerSuite(path string) (ServerSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ServerSuite{}, fmt.Errorf("reading suite: %w", err)
	}
	var ys yamlServerSuite
	if err := yaml.Unmarshal(data, &ys); err != nil {
		return ServerSuite{}, fmt.Errorf("parsing suite: %w", err)
	}

	suite := ServerSuite{Threshold: ys.Threshold}
	for i, yc := range ys.Checks {
		switch yc.Type {
		case "tools_list":
			suite.Checks = append(suite.Checks, ToolsListCheck{Expected: yc.ExpectTools})
		case "tool_call":
			var maxMS int64
			if yc.MaxLatency != "" {
				d, err := time.ParseDuration(yc.MaxLatency)
				if err != nil {
					return ServerSuite{}, fmt.Errorf("check %d (%s): invalid max_latency %q: %w", i, yc.Name, yc.MaxLatency, err)
				}
				maxMS = d.Milliseconds()
			}
			suite.Checks = append(suite.Checks, ToolCallCheck{
				Tool:           yc.Tool,
				Args:           yc.Args,
				ExpectNoError:  yc.ExpectNoError,
				ResultContains: yc.ResultContains,
				MaxLatencyMS:   maxMS,
			})
		default:
			return ServerSuite{}, fmt.Errorf("check %d (%s): unknown type %q", i, yc.Name, yc.Type)
		}
	}
	return suite, nil
}
