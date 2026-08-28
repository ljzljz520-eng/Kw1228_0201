package review

import (
	"path/filepath"
	"testing"

	"example.com/energycore/internal/catalog"
	"example.com/energycore/internal/clock"
	"example.com/energycore/internal/model"
	"example.com/energycore/internal/store"
)

func TestReviewPolicyAndDecision(t *testing.T) {
	repository, err := store.Open(filepath.Join(t.TempDir(), "review.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	sequence := clock.New(0)
	ids := clock.NewIDs(sequence)
	catalogue := catalog.NewService(repository, ids, sequence)
	service := NewService(repository, catalogue, ids, sequence)
	record, err := catalogue.Register(model.Record{Code: "R-1", Name: "Reviewable", Description: "Long enough containment description", Owner: "Review Lab", Classification: "beta"}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := catalogue.Submit(record.ID, "operator"); err != nil {
		t.Fatal(err)
	}
	result, err := service.Review(record.ID, model.DecisionApprove, "clear", "reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if result.Decision != model.DecisionApprove || result.Note != "clear" {
		t.Fatalf("result = %+v", result)
	}
}
