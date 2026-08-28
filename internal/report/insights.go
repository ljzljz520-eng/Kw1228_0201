package report

import (
	"sort"
	"strings"

	"example.com/energycore/internal/model"
)

type OwnerCount struct {
	Owner string `json:"owner"`
	Count int    `json:"count"`
}

type RiskFlag struct {
	RecordID string   `json:"record_id"`
	Reasons  []string `json:"reasons"`
}

func Owners(records []model.Record) []OwnerCount {
	counts := map[string]int{}
	for _, record := range records {
		owner := strings.TrimSpace(record.Owner)
		if owner == "" {
			owner = "unassigned"
		}
		counts[owner]++
	}
	result := make([]OwnerCount, 0, len(counts))
	for owner, count := range counts {
		result = append(result, OwnerCount{Owner: owner, Count: count})
	}
	sort.Slice(result, func(left, right int) bool { return result[left].Owner < result[right].Owner })
	return result
}

func RiskFlags(records []model.Record) []RiskFlag {
	flags := []RiskFlag{}
	for _, record := range records {
		reasons := []string{}
		if record.Status == model.StatusInReview {
			reasons = append(reasons, "awaiting review")
		}
		if strings.TrimSpace(record.Description) == "" {
			reasons = append(reasons, "missing description")
		}
		if record.Status == model.StatusRejected {
			reasons = append(reasons, "previously rejected")
		}
		if len(reasons) > 0 {
			flags = append(flags, RiskFlag{RecordID: record.ID, Reasons: reasons})
		}
	}
	return flags
}

func StatusOrder() []model.RecordStatus {
	return []model.RecordStatus{model.StatusDraft, model.StatusInReview, model.StatusApproved, model.StatusRejected, model.StatusArchived}
}

func RenderMarkdown(summary Summary) string {
	lines := []string{"# Energy Core Summary", "", "| Status | Count |", "| --- | ---: |"}
	for _, status := range StatusOrder() {
		lines = append(lines, "| "+status.String()+" | "+itoa(summary.ByStatus[status])+" |")
	}
	lines = append(lines, "", "Total records: "+itoa(summary.Total), "Pending review: "+itoa(summary.Pending))
	return strings.Join(lines, "\n")
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	result := ""
	for value > 0 {
		result = string(rune('0'+value%10)) + result
		value /= 10
	}
	return result
}
