package catalog

import (
	"errors"
	"fmt"
	"strings"

	"example.com/energycore/internal/model"
)

type Change struct {
	Field string
	From  string
	To    string
}

type ChangeSet struct {
	RecordID string
	Changes  []Change
}

func Diff(before, after model.Record) ChangeSet {
	changes := []Change{}
	if before.Name != after.Name {
		changes = append(changes, Change{Field: "name", From: before.Name, To: after.Name})
	}
	if before.Description != after.Description {
		changes = append(changes, Change{Field: "description", From: before.Description, To: after.Description})
	}
	if before.Owner != after.Owner {
		changes = append(changes, Change{Field: "owner", From: before.Owner, To: after.Owner})
	}
	if before.Classification != after.Classification {
		changes = append(changes, Change{Field: "classification", From: before.Classification, To: after.Classification})
	}
	if before.Status != after.Status {
		changes = append(changes, Change{Field: "status", From: string(before.Status), To: string(after.Status)})
	}
	return ChangeSet{RecordID: after.ID, Changes: changes}
}

func (c ChangeSet) HasChanges() bool {
	return len(c.Changes) > 0
}

func (c ChangeSet) Summary() string {
	parts := make([]string, 0, len(c.Changes))
	for _, change := range c.Changes {
		parts = append(parts, fmt.Sprintf("%s: %s -> %s", change.Field, change.From, change.To))
	}
	return strings.Join(parts, "; ")
}

func (s *Service) ValidateBulk(records []model.Record) []string {
	issues := []string{}
	seen := map[string]bool{}
	for index, record := range records {
		if err := model.ValidateRecord(record); err != nil {
			issues = append(issues, fmt.Sprintf("row %d: %v", index+1, err))
		}
		key := strings.ToUpper(strings.TrimSpace(record.Code))
		if key != "" && seen[key] {
			issues = append(issues, fmt.Sprintf("row %d: duplicate code %s", index+1, key))
		}
		seen[key] = true
	}
	return issues
}

func (s *Service) ArchiveBatch(ids []string, actor string) ([]model.Record, []string) {
	archived := []model.Record{}
	issues := []string{}
	for _, id := range ids {
		record, err := s.Archive(id, actor)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		archived = append(archived, record)
	}
	return archived, issues
}

func (s *Service) EnsureEditable(id string) (model.Record, error) {
	record, err := s.Get(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status.IsTerminal() {
		return model.Record{}, errors.New("terminal record is not editable")
	}
	return record, nil
}

func (s *Service) UpdateDescription(id, description, actor string) (model.Record, ChangeSet, error) {
	before, err := s.Get(id)
	if err != nil {
		return model.Record{}, ChangeSet{}, err
	}
	after, err := s.Update(id, model.Record{Description: description}, before.Revision, actor)
	if err != nil {
		return model.Record{}, ChangeSet{}, err
	}
	return after, Diff(before, after), nil
}
