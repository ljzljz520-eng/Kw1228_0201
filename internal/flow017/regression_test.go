package flow017

import (
	"testing"

	"example.com/energycore/internal/model"
)

func Test1228BusinessRegression(t *testing.T) {
	runtime := testRuntime(t)
	record, err := runtime.Catalog.Register(testRecord(), "operator")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Catalog.Submit(record.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Review.Review(record.ID, model.DecisionReject, "withdraw containment note", "reviewer"); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Catalog.Delete(record.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	result, err := runtime.Review.GetResult(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.WithdrawnNote != "" || result.Note != "" {
		t.Fatalf("review result retains invalid note: %+v", result)
	}
}
