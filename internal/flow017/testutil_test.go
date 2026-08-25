package flow017

import (
	"path/filepath"
	"testing"

	"example.com/energycore/internal/model"
)

func testRuntime(t *testing.T) *Runtime {
	t.Helper()
	runtime, err := OpenRuntime(filepath.Join(t.TempDir(), "energycore.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = runtime.Close() })
	return runtime
}

func testRecord() model.Record {
	return model.Record{Code: "SC-001", Name: "Helios Core", Description: "Contained fusion lattice for orbital research", Owner: "Aster Lab", Classification: "alpha"}
}
