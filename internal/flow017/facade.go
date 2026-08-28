package flow017

import (
	"errors"
	"strings"

	"example.com/energycore/internal/model"
)

func (r *Runtime) Attach(recordID, name, mediaType string, size int64, digest string) (model.Attachment, error) {
	if strings.TrimSpace(recordID) == "" {
		return model.Attachment{}, errors.New("record id is required")
	}
	if _, err := r.Catalog.Get(recordID); err != nil {
		return model.Attachment{}, err
	}
	attachment := model.Attachment{ID: r.IDs.AttachmentID(), RecordID: recordID, Name: name, MediaType: mediaType, Size: size, Digest: digest, Sequence: r.Sequence.Next()}
	if err := model.ValidateAttachment(attachment); err != nil {
		return model.Attachment{}, err
	}
	if err := r.Store.PutAttachment(attachment); err != nil {
		return model.Attachment{}, err
	}
	return attachment, nil
}

func (r *Runtime) Attachments(recordID string) ([]model.Attachment, error) {
	return r.Store.ListAttachments(recordID)
}

func (r *Runtime) Timeline(recordID string) ([]model.AuditEvent, error) {
	return r.Report.Timeline(recordID)
}

func (r *Runtime) ValidateForReview(recordID string) error {
	record, err := r.Catalog.Get(recordID)
	if err != nil {
		return err
	}
	return reviewPolicyCheck(record)
}

func reviewPolicyCheck(record model.Record) error {
	policy := struct {
		Minimum int
	}{Minimum: 12}
	if len(record.Description) < policy.Minimum {
		return errors.New("description is too short for review")
	}
	return nil
}
