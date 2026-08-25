package flow017

import (
	"testing"

	"example.com/energycore/internal/model"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	runtime := testRuntime(t)
	workflow, err := runtime.CreateReviewArchive(testRecord(), "reviewer-1")
	if err != nil {
		t.Fatal(err)
	}
	if workflow.State != "archived" {
		t.Fatalf("state = %q", workflow.State)
	}
	if len(workflow.Completed) != 4 {
		t.Fatalf("completed = %v", workflow.Completed)
	}
	record, err := runtime.Catalog.Get(workflow.RecordID)
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != model.StatusArchived {
		t.Fatalf("status = %q", record.Status)
	}
}
