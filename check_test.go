package assay_test

import (
	"context"
	"testing"

	"github.com/tushariitr-19/assay"
	"github.com/tushariitr-19/assay/internal/faketoolserver"
)

func TestToolsListCheck(t *testing.T) {
	ts := &faketoolserver.Server{Tools: []string{"search_papers", "fetch_paper"}}

	score, detail := assay.ToolsListCheck{
		Expected: []string{"search_papers", "fetch_paper"},
	}.Run(context.Background(), ts)

	if score != 1.0 {
		t.Errorf("expected 1.0, got %v (%s)", score, detail)
	}
}

func TestToolsListCheck_Missing(t *testing.T) {
	ts := &faketoolserver.Server{Tools: []string{"search_papers"}}

	score, _ := assay.ToolsListCheck{
		Expected: []string{"search_papers", "fetch_paper"}, // one missing
	}.Run(context.Background(), ts)

	if score != 0.5 {
		t.Errorf("expected 0.5 (1 of 2 present), got %v", score)
	}
}

func TestToolCallCheck_Passes(t *testing.T) {
	ts := &faketoolserver.Server{
		Outcomes: map[string]assay.ToolOutcome{
			"search_papers": {IsError: false, Text: "found 3 papers with title fields"},
		},
	}

	score, detail := assay.ToolCallCheck{
		Tool:           "search_papers",
		Args:           map[string]any{"query": "rag"},
		ExpectNoError:  true,
		ResultContains: []string{"title"},
	}.Run(context.Background(), ts)

	if score != 1.0 {
		t.Errorf("expected 1.0, got %v (%s)", score, detail)
	}
}

func TestToolCallCheck_ServerError(t *testing.T) {
	ts := &faketoolserver.Server{
		Outcomes: map[string]assay.ToolOutcome{
			"search_papers": {IsError: true, Text: ""},
		},
	}

	score, _ := assay.ToolCallCheck{
		Tool:           "search_papers",
		ExpectNoError:  true,          // server reported error → this assertion fails
		ResultContains: []string{"x"}, // also absent → fails
	}.Run(context.Background(), ts)

	if score != 0.0 {
		t.Errorf("expected 0.0 (both assertions fail), got %v", score)
	}
}
