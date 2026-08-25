package catalog

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"example.com/energycore/internal/clock"
	"example.com/energycore/internal/model"
	"example.com/energycore/internal/store"
)

var ErrConflict = errors.New("record revision conflict")

type Service struct {
	store *store.Store
	ids   *clock.IDs
	seq   clock.Clock
}

func NewService(repository *store.Store, ids *clock.IDs, sequence clock.Clock) *Service {
	return &Service{store: repository, ids: ids, seq: sequence}
}

func (s *Service) Register(input model.Record, actor string) (model.Record, error) {
	input = model.NormalizeRecord(input)
	if err := model.ValidateRecord(input); err != nil {
		return model.Record{}, err
	}
	input.ID = s.ids.RecordID()
	input.Revision = 1
	input.CreatedSeq = s.seq.Next()
	input.UpdatedSeq = input.CreatedSeq
	if err := s.store.PutRecord(input); err != nil {
		return model.Record{}, fmt.Errorf("persist record: %w", err)
	}
	if err := s.audit(input, "registered", actor, "core registered"); err != nil {
		return model.Record{}, err
	}
	return input, nil
}

func (s *Service) Get(id string) (model.Record, error) {
	if strings.TrimSpace(id) == "" {
		return model.Record{}, store.ErrNotFound
	}
	return s.store.GetRecord(id)
}

func (s *Service) Search(filter model.SearchFilter) ([]model.Record, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]model.Record, 0, len(records))
	for _, record := range records {
		if !filter.IncludeArchived && record.Status == model.StatusArchived {
			continue
		}
		if filter.Status != "" && record.Status != filter.Status {
			continue
		}
		if filter.Owner != "" && !strings.EqualFold(record.Owner, filter.Owner) {
			continue
		}
		if filter.Classification != "" && !strings.EqualFold(record.Classification, filter.Classification) {
			continue
		}
		if !contains(record, filter.Query) {
			continue
		}
		result = append(result, record)
	}
	sort.Slice(result, func(left, right int) bool { return result[left].ID < result[right].ID })
	return result, nil
}

func contains(record model.Record, query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return true
	}
	fields := []string{record.ID, record.Code, record.Name, record.Description, record.Owner, record.Classification}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}

func (s *Service) Update(id string, patch model.Record, expectedRevision int, actor string) (model.Record, error) {
	current, err := s.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if expectedRevision > 0 && current.Revision != expectedRevision {
		return model.Record{}, ErrConflict
	}
	if patch.Name != "" {
		current.Name = model.NormalizeText(patch.Name)
	}
	if patch.Description != "" {
		current.Description = model.NormalizeText(patch.Description)
	}
	if patch.Owner != "" {
		current.Owner = model.NormalizeText(patch.Owner)
	}
	if patch.Classification != "" {
		current.Classification = model.NormalizeText(patch.Classification)
	}
	if patch.Status != "" && patch.Status != current.Status {
		if !model.AllowedTransition(current.Status, patch.Status) {
			return model.Record{}, model.ErrInvalidTransition
		}
		current.Status = patch.Status
	}
	current.Revision++
	current.UpdatedSeq = s.seq.Next()
	if err := model.ValidateRecord(current); err != nil {
		return model.Record{}, err
	}
	if err := s.store.PutRecord(current); err != nil {
		return model.Record{}, err
	}
	if err := s.audit(current, "updated", actor, "core details changed"); err != nil {
		return model.Record{}, err
	}
	return current, nil
}

func (s *Service) Delete(id, actor string) error {
	record, err := s.Get(id)
	if err != nil {
		return err
	}
	if record.Status == model.StatusArchived {
		return errors.New("archived core cannot be deleted")
	}
	if err := s.store.DeleteRecord(id); err != nil {
		return err
	}
	if err := s.store.DeleteEvents(id); err != nil {
		return err
	}
	return nil
}

func (s *Service) audit(record model.Record, action, actor, message string) error {
	return s.store.PutEvent(model.AuditEvent{ID: s.ids.EventID(), RecordID: record.ID, Action: action, Actor: actor, Message: message, Sequence: s.seq.Next(), Revision: record.Revision})
}

func (s *Service) Events(id string) ([]model.AuditEvent, error) {
	return s.store.ListEvents(id)
}
