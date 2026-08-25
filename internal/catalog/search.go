package catalog

import (
	"sort"
	"strings"

	"example.com/energycore/internal/model"
)

type Group struct {
	Key    string
	Count  int
	Latest model.Record
}

func GroupByClassification(records []model.Record) []Group {
	groups := map[string]Group{}
	for _, record := range records {
		key := strings.TrimSpace(record.Classification)
		if key == "" {
			key = "unclassified"
		}
		group := groups[key]
		group.Key = key
		group.Count++
		if group.Latest.ID == "" || record.UpdatedSeq > group.Latest.UpdatedSeq {
			group.Latest = record
		}
		groups[key] = group
	}
	result := make([]Group, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	sort.Slice(result, func(left, right int) bool { return result[left].Key < result[right].Key })
	return result
}

func SortByRevision(records []model.Record) []model.Record {
	result := append([]model.Record(nil), records...)
	sort.SliceStable(result, func(left, right int) bool {
		if result[left].Revision == result[right].Revision {
			return result[left].ID < result[right].ID
		}
		return result[left].Revision > result[right].Revision
	})
	return result
}

func SummarizeStatuses(records []model.Record) map[model.RecordStatus]int {
	counts := map[model.RecordStatus]int{}
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}
