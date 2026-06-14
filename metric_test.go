package assay_test

import (
	"context"
	"testing"

	"github.com/tushariitr-19/assay"
	"github.com/tushariitr-19/assay/internal/fakeagent"
)

func TestToolCorrectness_Unordered(t *testing.T) {
	// The fake agent reports that it called "topic_search".
	agent := fakeagent.New(assay.Result{
		Output:    "found 3 papers on RAG",
		ToolCalls: []assay.ToolCall{{Name: "topic_search"}},
	})

	result, err := agent.Run(context.Background(), "find papers on RAG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exp := assay.Expectation{Tools: []string{"topic_search"}}
	score, detail := assay.ToolCorrectness{}.Score(result, exp)

	if score != 1.0 {
		t.Errorf("expected score 1.0, got %v (%s)", score, detail)
	}
}

func TestToolCorrectness_WrongTool(t *testing.T) {
	agent := fakeagent.New(assay.Result{
		ToolCalls: []assay.ToolCall{{Name: "paper_analysis"}},
	})
	result, _ := agent.Run(context.Background(), "anything")

	exp := assay.Expectation{Tools: []string{"topic_search"}}
	score, _ := assay.ToolCorrectness{}.Score(result, exp)

	if score != 0.0 {
		t.Errorf("expected score 0.0 for wrong tool, got %v", score)
	}
}
