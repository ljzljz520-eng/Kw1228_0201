package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingCode       = errors.New("core code is required")
	ErrMissingName       = errors.New("core name is required")
	ErrMissingOwner      = errors.New("core owner is required")
	ErrInvalidTransition = errors.New("status transition is not allowed")
	ErrInvalidDecision   = errors.New("review decision is invalid")
)

func NormalizeText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func ValidateRecord(record Record) error {
	if NormalizeText(record.Code) == "" {
		return ErrMissingCode
	}
	if NormalizeText(record.Name) == "" {
		return ErrMissingName
	}
	if NormalizeText(record.Owner) == "" {
		return ErrMissingOwner
	}
	if record.Status != "" && !containsStatus(record.Status) {
		return fmt.Errorf("unknown status %q", record.Status)
	}
	return nil
}

func ValidateReview(decision ReviewDecision, note string) error {
	if !decision.Valid() {
		return ErrInvalidDecision
	}
	if decision == DecisionReject && NormalizeText(note) == "" {
		return errors.New("rejection needs a note")
	}
	return nil
}

func ValidateAttachment(attachment Attachment) error {
	if NormalizeText(attachment.Name) == "" {
		return errors.New("attachment name is required")
	}
	if attachment.Size < 0 {
		return errors.New("attachment size cannot be negative")
	}
	return nil
}

func containsStatus(status RecordStatus) bool {
	for _, candidate := range Statuses() {
		if status == candidate {
			return true
		}
	}
	return false
}

func NormalizeRecord(record Record) Record {
	record.Code = strings.ToUpper(NormalizeText(record.Code))
	record.Name = NormalizeText(record.Name)
	record.Description = NormalizeText(record.Description)
	record.Owner = NormalizeText(record.Owner)
	record.Classification = NormalizeText(record.Classification)
	record.ReviewNote = NormalizeText(record.ReviewNote)
	record.WithdrawnNote = NormalizeText(record.WithdrawnNote)
	if record.Status == "" {
		record.Status = StatusDraft
	}
	return record
}
