package flow017

import (
	"testing"

	"example.com/energycore/internal/model"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	runtime := testRuntime(t)
	created, err := runtime.Catalog.Register(testRecord(), "operator")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := runtime.SearchUpdatePublish(created.Code, structPatch(), "operator")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != "approved" || updated.Name != "Helios Core Revised" {
		t.Fatalf("updated = %+v", updated)
	}
}

func structPatch() model.Record {
	return model.Record{Name: "Helios Core Revised"}
}
