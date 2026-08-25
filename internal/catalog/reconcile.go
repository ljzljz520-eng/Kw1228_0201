package catalog

import (
	"fmt"
	"sort"
	"strings"

	"example.com/energycore/internal/model"
)

type ReconcileReport struct {
	Checked       int      `json:"checked"`
	Ready         int      `json:"ready"`
	Warnings      []string `json:"warnings"`
	DuplicateCode []string `json:"duplicate_code"`
}

func Reconcile(records []model.Record) ReconcileReport {
	report := ReconcileReport{Checked: len(records), Warnings: []string{}, DuplicateCode: []string{}}
	seen := map[string]string{}
	for _, record := range records {
		if record.Status == model.StatusApproved || record.Status == model.StatusArchived {
			report.Ready++
		}
		code := strings.ToUpper(strings.TrimSpace(record.Code))
		if code != "" {
			if previous, ok := seen[code]; ok {
				report.DuplicateCode = append(report.DuplicateCode, code)
				report.Warnings = append(report.Warnings, fmt.Sprintf("code %s appears in %s and %s", code, previous, record.ID))
			} else {
				seen[code] = record.ID
			}
		}
		if strings.TrimSpace(record.Owner) == "" {
			report.Warnings = append(report.Warnings, record.ID+" has no owner")
		}
	}
	sort.Strings(report.DuplicateCode)
	sort.Strings(report.Warnings)
	return report
}

func FindDuplicates(records []model.Record) map[string][]string {
	groups := map[string][]string{}
	for _, record := range records {
		code := strings.ToUpper(strings.TrimSpace(record.Code))
		if code == "" {
			continue
		}
		groups[code] = append(groups[code], record.ID)
	}
	duplicates := map[string][]string{}
	for code, ids := range groups {
		if len(ids) > 1 {
			duplicates[code] = append([]string(nil), ids...)
		}
	}
	return duplicates
}

func ApplyClassification(records []model.Record, classification string) []model.Record {
	result := make([]model.Record, 0, len(records))
	for _, record := range records {
		copyRecord := record.Clone()
		if strings.TrimSpace(classification) != "" {
			copyRecord.Classification = model.NormalizeText(classification)
		}
		result = append(result, copyRecord)
	}
	return result
}

func TransitionPlan(records []model.Record, target model.RecordStatus) []string {
	plan := []string{}
	for _, record := range records {
		if model.AllowedTransition(record.Status, target) && record.Status != target {
			plan = append(plan, record.ID)
		}
	}
	sort.Strings(plan)
	return plan
}

func (s *Service) ReconcileStored() (ReconcileReport, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return ReconcileReport{}, err
	}
	return Reconcile(records), nil
}
