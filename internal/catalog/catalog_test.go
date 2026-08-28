package catalog

import (
	"path/filepath"
	"testing"

	"example.com/energycore/internal/clock"
	"example.com/energycore/internal/model"
	"example.com/energycore/internal/store"
)

func TestCatalogSearchAndConflict(t *testing.T) {
	repository, err := store.Open(filepath.Join(t.TempDir(), "catalog.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	sequence := clock.New(0)
	service := NewService(repository, clock.NewIDs(sequence), sequence)
	record, err := service.Register(model.Record{Code: "A-1", Name: "Alpha", Owner: "North", Classification: "alpha"}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	items, err := service.Search(model.SearchFilter{Query: "alpha"})
	if err != nil || len(items) != 1 {
		t.Fatalf("items = %v err=%v", items, err)
	}
	if _, err := service.Update(record.ID, model.Record{Name: "Changed"}, record.Revision+1, "operator"); err != ErrConflict {
		t.Fatalf("conflict err = %v", err)
	}
}
