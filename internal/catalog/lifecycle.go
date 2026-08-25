package catalog

import (
	"errors"
	"fmt"

	"example.com/energycore/internal/model"
)

func (s *Service) Submit(id, actor string) (model.Record, error) {
	record, err := s.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusDraft && record.Status != model.StatusRejected {
		return model.Record{}, fmt.Errorf("record %s is not ready for review", id)
	}
	return s.Update(id, model.Record{Status: model.StatusInReview}, record.Revision, actor)
}

func (s *Service) Archive(id, actor string) (model.Record, error) {
	record, err := s.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusApproved && record.Status != model.StatusRejected {
		return model.Record{}, errors.New("only reviewed cores can be archived")
	}
	record.Status = model.StatusArchived
	record.Revision++
	record.ArchivedSeq = s.seq.Next()
	record.UpdatedSeq = record.ArchivedSeq
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	if err := s.audit(record, "archived", actor, "core archived"); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Restore(id, actor string) (model.Record, error) {
	record, err := s.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusArchived {
		return model.Record{}, errors.New("only archived cores can be restored")
	}
	record.Status = model.StatusDraft
	record.Revision++
	record.UpdatedSeq = s.seq.Next()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	if err := s.audit(record, "restored", actor, "core restored for review"); err != nil {
		return model.Record{}, err
	}
	return record, nil
}
