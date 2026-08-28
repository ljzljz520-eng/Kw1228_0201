package report

import (
	"sort"
	"strings"

	"example.com/energycore/internal/model"
)

type SummaryDelta struct {
	TotalChange   int                        `json:"total_change"`
	StatusChanges map[model.RecordStatus]int `json:"status_changes"`
	PendingChange int                        `json:"pending_change"`
}

func Compare(before, after Summary) SummaryDelta {
	changes := map[model.RecordStatus]int{}
	for _, status := range model.Statuses() {
		changes[status] = after.ByStatus[status] - before.ByStatus[status]
	}
	return SummaryDelta{TotalChange: after.Total - before.Total, StatusChanges: changes, PendingChange: after.Pending - before.Pending}
}

func AuditActionCounts(events []model.AuditEvent) map[string]int {
	counts := map[string]int{}
	for _, event := range events {
		counts[event.Action]++
	}
	return counts
}

func RenderOwners(owners []OwnerCount) string {
	items := append([]OwnerCount(nil), owners...)
	sort.Slice(items, func(left, right int) bool { return items[left].Owner < items[right].Owner })
	lines := []string{"owner,count"}
	for _, item := range items {
		lines = append(lines, item.Owner+","+number(item.Count))
	}
	return strings.Join(lines, "\n")
}

func number(value int) string {
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

func (s Summary) IsEmpty() bool {
	return s.Total == 0 && s.Pending == 0 && len(s.ByStatus) == 0
}

func (s Summary) StatusCount(status model.RecordStatus) int {
	return s.ByStatus[status]
}
