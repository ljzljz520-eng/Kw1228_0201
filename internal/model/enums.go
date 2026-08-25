package model

func (s RecordStatus) String() string {
	if s == "" {
		return string(StatusDraft)
	}
	return string(s)
}

func (s RecordStatus) IsTerminal() bool {
	switch s {
	case StatusArchived, StatusRejected:
		return true
	default:
		return false
	}
}

func (d ReviewDecision) Valid() bool {
	switch d {
	case DecisionApprove, DecisionReject, DecisionRequest:
		return true
	default:
		return false
	}
}

func AllowedTransition(from, to RecordStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusDraft:
		return to == StatusInReview || to == StatusArchived
	case StatusInReview:
		return to == StatusApproved || to == StatusRejected || to == StatusDraft
	case StatusApproved:
		return to == StatusArchived || to == StatusInReview
	case StatusRejected:
		return to == StatusDraft || to == StatusArchived
	case StatusArchived:
		return false
	default:
		return false
	}
}

func Statuses() []RecordStatus {
	return []RecordStatus{StatusDraft, StatusInReview, StatusApproved, StatusRejected, StatusArchived}
}

func WorkflowStates() []string {
	return []string{"created", "reviewing", "confirmed", "archived", "reported"}
}
