package report

import (
	"sort"

	"example.com/energycore/internal/catalog"
	"example.com/energycore/internal/model"
	"example.com/energycore/internal/store"
)

type Summary struct {
	Total       int                        `json:"total"`
	ByStatus    map[model.RecordStatus]int `json:"by_status"`
	ByClass     []catalog.Group            `json:"by_classification"`
	Owners      []OwnerCount               `json:"owners"`
	Risks       []RiskFlag                 `json:"risks"`
	Pending     int                        `json:"pending_review"`
	AuditEvents int                        `json:"audit_events"`
	Actions     map[string]int             `json:"audit_actions"`
}

type Service struct {
	store   *store.Store
	catalog *catalog.Service
}

func NewService(repository *store.Store, catalogue *catalog.Service) *Service {
	return &Service{store: repository, catalog: catalogue}
}

func (s *Service) Summary() (Summary, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return Summary{}, err
	}
	events, err := s.store.ListEvents("")
	if err != nil {
		return Summary{}, err
	}
	statuses := catalog.SummarizeStatuses(records)
	pending := statuses[model.StatusInReview]
	return Summary{Total: len(records), ByStatus: statuses, ByClass: catalog.GroupByClassification(records), Owners: Owners(records), Risks: RiskFlags(records), Pending: pending, AuditEvents: len(events), Actions: AuditActionCounts(events)}, nil
}

func (s *Service) Timeline(id string) ([]model.AuditEvent, error) {
	events, err := s.catalog.Events(id)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(events, func(left, right int) bool { return events[left].Sequence < events[right].Sequence })
	return events, nil
}

func FormatStatus(status model.RecordStatus) string {
	if status == "" {
		return "draft"
	}
	return status.String()
}
