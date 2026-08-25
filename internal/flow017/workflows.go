package flow017

import (
	"errors"
	"fmt"

	"example.com/energycore/internal/model"
)

func (r *Runtime) CreateReviewArchive(input model.Record, actor string) (model.Workflow, error) {
	record, err := r.Catalog.Register(input, actor)
	if err != nil {
		return model.Workflow{}, err
	}
	workflow := model.Workflow{ID: r.IDs.WorkflowID(), RecordID: record.ID, Name: "create-review-archive", State: "created", Steps: []string{"create", "review", "confirm", "archive"}, Completed: []string{"create"}, StartedSeq: r.Sequence.Next(), UpdatedSeq: r.Sequence.Current()}
	if err := r.Store.PutWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	if _, err := r.Catalog.Submit(record.ID, actor); err != nil {
		return model.Workflow{}, err
	}
	workflow.State = "reviewing"
	workflow.Completed = append(workflow.Completed, "review")
	workflow.UpdatedSeq = r.Sequence.Next()
	if _, err := r.Review.Review(record.ID, model.DecisionApprove, "validated containment and output", actor); err != nil {
		return model.Workflow{}, err
	}
	workflow.State = "confirmed"
	workflow.Completed = append(workflow.Completed, "confirm")
	workflow.UpdatedSeq = r.Sequence.Next()
	if _, err := r.Catalog.Archive(record.ID, actor); err != nil {
		return model.Workflow{}, err
	}
	workflow.State = "archived"
	workflow.Completed = append(workflow.Completed, "archive")
	workflow.UpdatedSeq = r.Sequence.Next()
	if err := r.Store.PutWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}

func (r *Runtime) SearchUpdatePublish(query string, patch model.Record, actor string) (model.Record, error) {
	records, err := r.Catalog.Search(model.SearchFilter{Query: query, IncludeArchived: true})
	if err != nil {
		return model.Record{}, err
	}
	if len(records) == 0 {
		return model.Record{}, errors.New("no matching core")
	}
	record := records[0]
	updated, err := r.Catalog.Update(record.ID, patch, record.Revision, actor)
	if err != nil {
		return model.Record{}, err
	}
	if updated.Status == model.StatusDraft || updated.Status == model.StatusRejected {
		if _, err := r.Catalog.Submit(updated.ID, actor); err != nil {
			return model.Record{}, err
		}
	}
	if _, err := r.Review.Review(updated.ID, model.DecisionApprove, "published after update", actor); err != nil {
		return model.Record{}, err
	}
	final, err := r.Catalog.Get(updated.ID)
	if err != nil {
		return model.Record{}, err
	}
	return final, nil
}

func (r *Runtime) ImportReport(input, actor string) (model.ImportReport, ReportSummary, error) {
	reportValue, err := r.Import.ImportCSV(input, actor)
	if err != nil {
		return model.ImportReport{}, ReportSummary{}, err
	}
	summary, err := r.Report.Summary()
	if err != nil {
		return model.ImportReport{}, ReportSummary{}, err
	}
	return reportValue, ReportSummary{Total: summary.Total, Imported: len(reportValue.ImportedIDs), Rejected: len(reportValue.Rejected)}, nil
}

type ReportSummary struct {
	Total    int
	Imported int
	Rejected int
}

func (r *Runtime) Workflow(id string) (model.Workflow, error) {
	if id == "" {
		return model.Workflow{}, fmt.Errorf("workflow id is required")
	}
	return r.Store.GetWorkflow(id)
}
