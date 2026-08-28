package report

import (
	"sort"
	"strings"

	"example.com/energycore/internal/model"
)

type HistoryEntry struct {
	Sequence int64  `json:"sequence"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Message  string `json:"message"`
	Revision int    `json:"revision"`
}

type History struct {
	RecordID string         `json:"record_id"`
	Label    string         `json:"label"`
	Entries  []HistoryEntry `json:"entries"`
}

func BuildHistory(record model.Record, events []model.AuditEvent) History {
	entries := make([]HistoryEntry, 0, len(events))
	for _, event := range events {
		entries = append(entries, HistoryEntry{Sequence: event.Sequence, Action: event.Action, Actor: event.Actor, Message: event.Message, Revision: event.Revision})
	}
	sort.SliceStable(entries, func(left, right int) bool { return entries[left].Sequence < entries[right].Sequence })
	return History{RecordID: record.ID, Label: record.DisplayLabel(), Entries: entries}
}

func (h History) LatestAction() string {
	if len(h.Entries) == 0 {
		return ""
	}
	return h.Entries[len(h.Entries)-1].Action
}

func (h History) Actors() []string {
	seen := map[string]bool{}
	actors := []string{}
	for _, entry := range h.Entries {
		actor := strings.TrimSpace(entry.Actor)
		if actor != "" && !seen[actor] {
			seen[actor] = true
			actors = append(actors, actor)
		}
	}
	sort.Strings(actors)
	return actors
}

func (h History) HasAction(action string) bool {
	for _, entry := range h.Entries {
		if entry.Action == action {
			return true
		}
	}
	return false
}

func (h History) RevisionSpan() int {
	if len(h.Entries) == 0 {
		return 0
	}
	minimum := h.Entries[0].Revision
	maximum := minimum
	for _, entry := range h.Entries[1:] {
		if entry.Revision < minimum {
			minimum = entry.Revision
		}
		if entry.Revision > maximum {
			maximum = entry.Revision
		}
	}
	return maximum - minimum
}

func (s *Service) RecordHistory(id string) (History, error) {
	record, err := s.catalog.Get(id)
	if err != nil {
		return History{}, err
	}
	events, err := s.catalog.Events(id)
	if err != nil {
		return History{}, err
	}
	return BuildHistory(record, events), nil
}
