package review

import (
	"errors"
	"strings"

	"example.com/energycore/internal/model"
)

type Policy struct {
	RequiredClassification []string
	MinimumDescription     int
}

func DefaultPolicy() Policy {
	return Policy{RequiredClassification: []string{"alpha", "beta", "gamma"}, MinimumDescription: 12}
}

func (p Policy) Check(record model.Record) error {
	if strings.TrimSpace(record.Description) == "" {
		return errors.New("description is required for review")
	}
	if len([]rune(record.Description)) < p.MinimumDescription {
		return errors.New("description is too short for review")
	}
	if !p.acceptsClassification(record.Classification) {
		return errors.New("classification is outside review policy")
	}
	return nil
}

func (p Policy) acceptsClassification(classification string) bool {
	if len(p.RequiredClassification) == 0 {
		return true
	}
	for _, allowed := range p.RequiredClassification {
		if strings.EqualFold(strings.TrimSpace(classification), allowed) {
			return true
		}
	}
	return false
}

func DecisionLabel(decision model.ReviewDecision) string {
	switch decision {
	case model.DecisionApprove:
		return "approved"
	case model.DecisionReject:
		return "rejected"
	case model.DecisionRequest:
		return "changes requested"
	default:
		return "unknown"
	}
}
