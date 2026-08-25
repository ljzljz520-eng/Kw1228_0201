package review

import (
	"fmt"
	"strings"

	"example.com/energycore/internal/model"
)

type ChecklistItem struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Passed   bool   `json:"passed"`
}

type Checklist struct {
	Items []ChecklistItem `json:"items"`
}

func NewChecklist() Checklist {
	return Checklist{Items: []ChecklistItem{{Key: "containment", Label: "Containment evidence", Required: true}, {Key: "output", Label: "Output envelope", Required: true}, {Key: "owner", Label: "Responsible owner", Required: true}, {Key: "classification", Label: "Classification", Required: true}}}
}

func (c Checklist) Mark(key string, passed bool) Checklist {
	copyChecklist := Checklist{Items: append([]ChecklistItem(nil), c.Items...)}
	for index, item := range copyChecklist.Items {
		if item.Key == key {
			copyChecklist.Items[index].Passed = passed
		}
	}
	return copyChecklist
}

func (c Checklist) Complete() bool {
	for _, item := range c.Items {
		if item.Required && !item.Passed {
			return false
		}
	}
	return len(c.Items) > 0
}

func (c Checklist) Missing() []string {
	missing := []string{}
	for _, item := range c.Items {
		if item.Required && !item.Passed {
			missing = append(missing, item.Key)
		}
	}
	return missing
}

func BuildChecklist(record model.Record) Checklist {
	checklist := NewChecklist()
	checklist = checklist.Mark("owner", strings.TrimSpace(record.Owner) != "")
	checklist = checklist.Mark("classification", strings.TrimSpace(record.Classification) != "")
	checklist = checklist.Mark("containment", strings.Contains(strings.ToLower(record.Description), "contain"))
	checklist = checklist.Mark("output", strings.Contains(strings.ToLower(record.Description), "output") || len(record.Description) >= 24)
	return checklist
}

func RequireChecklist(record model.Record) error {
	checklist := BuildChecklist(record)
	if checklist.Complete() {
		return nil
	}
	return fmt.Errorf("review checklist incomplete: %s", strings.Join(checklist.Missing(), ", "))
}

func (s *Service) ReviewWithChecklist(id string, decision model.ReviewDecision, note, reviewer string) (model.ReviewResult, Checklist, error) {
	record, err := s.catalog.Get(id)
	if err != nil {
		return model.ReviewResult{}, Checklist{}, err
	}
	checklist := BuildChecklist(record)
	if decision == model.DecisionApprove && !checklist.Complete() {
		return model.ReviewResult{}, checklist, RequireChecklist(record)
	}
	result, err := s.Review(id, decision, note, reviewer)
	return result, checklist, err
}
