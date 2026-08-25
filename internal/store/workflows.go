package store

import (
	"example.com/energycore/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) PutWorkflow(workflow model.Workflow) error {
	return s.withUpdate(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket([]byte("workflows")), workflow.ID, workflow)
	})
}

func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	var workflow model.Workflow
	err := s.withView(func(tx *bbolt.Tx) error {
		return getJSON(tx.Bucket([]byte("workflows")), id, &workflow)
	})
	return workflow, err
}

func (s *Store) ListWorkflows(recordID string) ([]model.Workflow, error) {
	workflows := []model.Workflow{}
	err := s.withView(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("workflows")).ForEach(func(_, value []byte) error {
			var workflow model.Workflow
			if err := decode(value, &workflow); err != nil {
				return err
			}
			if recordID == "" || workflow.RecordID == recordID {
				workflows = append(workflows, workflow)
			}
			return nil
		})
	})
	return workflows, err
}
