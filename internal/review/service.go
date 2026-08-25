package review

import (
	"errors"
	"fmt"
	"strings"

	"example.com/energycore/internal/catalog"
	"example.com/energycore/internal/clock"
	"example.com/energycore/internal/model"
	"example.com/energycore/internal/store"
)

type Service struct {
	store     *store.Store
	catalog   *catalog.Service
	ids       *clock.IDs
	seq       clock.Clock
	last      model.ReviewResult
	lastValid bool
}

func NewService(repository *store.Store, catalogue *catalog.Service, ids *clock.IDs, sequence clock.Clock) *Service {
	return &Service{store: repository, catalog: catalogue, ids: ids, seq: sequence}
}

func (s *Service) Review(id string, decision model.ReviewDecision, note, reviewer string) (model.ReviewResult, error) {
	if err := model.ValidateReview(decision, note); err != nil {
		return model.ReviewResult{}, err
	}
	record, err := s.catalog.Get(id)
	if err != nil {
		return model.ReviewResult{}, err
	}
	if record.Status != model.StatusInReview {
		return model.ReviewResult{}, fmt.Errorf("record %s is not awaiting review", id)
	}
	note = model.NormalizeText(note)
	nextStatus := model.StatusApproved
	if decision == model.DecisionReject {
		nextStatus = model.StatusRejected
	} else if decision == model.DecisionRequest {
		nextStatus = model.StatusDraft
	}
	updated, err := s.catalog.Update(id, model.Record{Status: nextStatus, ReviewNote: note}, record.Revision, reviewer)
	if err != nil {
		return model.ReviewResult{}, err
	}
	result := model.ReviewResult{RecordID: id, Decision: decision, Note: note, Reviewer: strings.TrimSpace(reviewer), Revision: updated.Revision, Sequence: s.seq.Next()}
	if decision == model.DecisionReject {
		result.WithdrawnNote = note
		updated.WithdrawnNote = note
		if err := s.store.PutRecord(updated); err != nil {
			return model.ReviewResult{}, err
		}
	}
	s.last = result
	s.lastValid = true
	if err := s.store.PutEvent(model.AuditEvent{ID: s.ids.EventID(), RecordID: id, Action: "reviewed", Actor: reviewer, Message: note, Sequence: s.seq.Next(), Revision: updated.Revision}); err != nil {
		return model.ReviewResult{}, err
	}
	return result, nil
}

func (s *Service) GetResult(id string) (model.ReviewResult, error) {
	record, err := s.catalog.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) && s.lastValid {
			return s.last, nil
		}
		return model.ReviewResult{}, err
	}
	events, err := s.store.ListEvents(id)
	if err != nil {
		return model.ReviewResult{}, err
	}
	result := model.ReviewResult{RecordID: id, Revision: record.Revision, WithdrawnNote: record.WithdrawnNote}
	for _, event := range events {
		if event.Action == "reviewed" {
			result.Note = event.Message
			result.Reviewer = event.Actor
			result.Sequence = event.Sequence
			if record.Status == model.StatusRejected {
				result.Decision = model.DecisionReject
			} else if record.Status == model.StatusApproved {
				result.Decision = model.DecisionApprove
			} else {
				result.Decision = model.DecisionRequest
			}
		}
	}
	return result, nil
}

func (s *Service) Pending() ([]model.Record, error) {
	return s.catalog.Search(model.SearchFilter{Status: model.StatusInReview})
}

func (s *Service) ClearReview(id string) error {
	record, err := s.catalog.Get(id)
	if err != nil {
		return err
	}
	record.ReviewNote = ""
	record.WithdrawnNote = ""
	return s.store.PutRecord(record)
}
