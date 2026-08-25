package flow017

import (
	"errors"
	"fmt"
	"strings"

	"example.com/energycore/internal/model"
)

type QueueItem struct {
	Record model.Record
	Age    int64
}

type Snapshot struct {
	Records     []model.Record
	Workflows   []model.Workflow
	Events      []model.AuditEvent
	Attachments []model.Attachment
}

func (r *Runtime) ReviewQueue() ([]QueueItem, error) {
	records, err := r.Review.Pending()
	if err != nil {
		return nil, err
	}
	items := make([]QueueItem, 0, len(records))
	for _, record := range records {
		items = append(items, QueueItem{Record: record, Age: r.Sequence.Current() - record.UpdatedSeq})
	}
	return items, nil
}

func (r *Runtime) ApproveQueue(ids []string, reviewer string) ([]model.ReviewResult, []string) {
	results := []model.ReviewResult{}
	issues := []string{}
	for _, id := range ids {
		result, _, err := r.Review.ReviewWithChecklist(id, model.DecisionApprove, "queue approval", reviewer)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		results = append(results, result)
	}
	return results, issues
}

func (r *Runtime) Snapshot(recordID string) (Snapshot, error) {
	records, err := r.Catalog.Search(model.SearchFilter{Query: recordID, IncludeArchived: true})
	if err != nil {
		return Snapshot{}, err
	}
	if recordID != "" && len(records) == 0 {
		return Snapshot{}, errors.New("record not found for snapshot")
	}
	workflows, err := r.Store.ListWorkflows(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	events, err := r.Store.ListEvents(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := r.Store.ListAttachments(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Records: records, Workflows: workflows, Events: events, Attachments: attachments}, nil
}

func (r *Runtime) Publish(recordID, actor string) (model.Record, error) {
	record, err := r.Catalog.Get(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status == model.StatusArchived {
		return model.Record{}, errors.New("archived record is already published")
	}
	if record.Status == model.StatusDraft || record.Status == model.StatusRejected {
		if _, err := r.Catalog.Submit(recordID, actor); err != nil {
			return model.Record{}, err
		}
	}
	if record.Status != model.StatusApproved {
		if _, err := r.Review.Review(recordID, model.DecisionApprove, "published from operations", actor); err != nil {
			return model.Record{}, err
		}
	}
	return r.Catalog.Get(recordID)
}

func (r *Runtime) ArchivePublished(ids []string, actor string) (int, []string) {
	count := 0
	issues := []string{}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			issues = append(issues, "empty record id")
			continue
		}
		if _, err := r.Catalog.Archive(id, actor); err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		count++
	}
	return count, issues
}

func (r *Runtime) ResetReviewState(recordID string) error {
	if _, err := r.Catalog.Get(recordID); err != nil {
		return err
	}
	return r.Review.ClearReview(recordID)
}
