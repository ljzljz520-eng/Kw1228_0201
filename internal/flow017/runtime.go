package flow017

import (
	"example.com/energycore/internal/catalog"
	"example.com/energycore/internal/clock"
	"example.com/energycore/internal/importer"
	"example.com/energycore/internal/report"
	"example.com/energycore/internal/review"
	"example.com/energycore/internal/store"
)

type Runtime struct {
	Store    *store.Store
	Sequence *clock.Sequence
	IDs      *clock.IDs
	Catalog  *catalog.Service
	Review   *review.Service
	Import   *importer.Service
	Report   *report.Service
}

func NewRuntime(repository *store.Store) *Runtime {
	sequence := clock.New(0)
	ids := clock.NewIDs(sequence)
	catalogue := catalog.NewService(repository, ids, sequence)
	reviewer := review.NewService(repository, catalogue, ids, sequence)
	return &Runtime{Store: repository, Sequence: sequence, IDs: ids, Catalog: catalogue, Review: reviewer, Import: importer.NewService(catalogue), Report: report.NewService(repository, catalogue)}
}

func OpenRuntime(path string) (*Runtime, error) {
	repository, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	return NewRuntime(repository), nil
}

func (r *Runtime) Close() error {
	return r.Store.Close()
}
