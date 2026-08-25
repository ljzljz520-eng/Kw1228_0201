package store

import (
	"path/filepath"
	"testing"

	"example.com/energycore/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := model.Record{ID: "core-persist", Code: "P-1", Name: "Persistent Core", Owner: "Archive Lab", Status: model.StatusDraft, Revision: 1}
	if err := first.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Code != record.Code || loaded.Name != record.Name {
		t.Fatalf("loaded = %+v", loaded)
	}
}
