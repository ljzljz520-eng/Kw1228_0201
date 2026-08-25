package model

import (
	"sort"
	"strings"
)

func (r Record) DisplayLabel() string {
	if r.Code == "" {
		return r.Name
	}
	if r.Name == "" {
		return r.Code
	}
	return r.Code + " - " + r.Name
}

func (r Record) Keywords() []string {
	values := []string{r.Code, r.Name, r.Owner, r.Classification}
	seen := map[string]bool{}
	keywords := []string{}
	for _, value := range values {
		for _, word := range strings.Fields(strings.ToLower(value)) {
			if !seen[word] {
				seen[word] = true
				keywords = append(keywords, word)
			}
		}
	}
	sort.Strings(keywords)
	return keywords
}

func (r Record) Clone() Record {
	return Record{ID: r.ID, Code: r.Code, Name: r.Name, Status: r.Status, Description: r.Description, Owner: r.Owner, Classification: r.Classification, ReviewNote: r.ReviewNote, WithdrawnNote: r.WithdrawnNote, Revision: r.Revision, CreatedSeq: r.CreatedSeq, UpdatedSeq: r.UpdatedSeq, ArchivedSeq: r.ArchivedSeq}
}

func (r Record) IsVisible(includeArchived bool) bool {
	return includeArchived || r.Status != StatusArchived
}

func (w Workflow) NextStep() string {
	if len(w.Completed) >= len(w.Steps) {
		return ""
	}
	return w.Steps[len(w.Completed)]
}

func (w Workflow) IsComplete() bool {
	return len(w.Steps) > 0 && len(w.Completed) == len(w.Steps)
}

func (w Workflow) CompleteStep(step string) Workflow {
	copyWorkflow := w
	for _, done := range w.Completed {
		if done == step {
			return copyWorkflow
		}
	}
	copyWorkflow.Completed = append(append([]string(nil), w.Completed...), step)
	if copyWorkflow.IsComplete() {
		copyWorkflow.State = "complete"
	}
	return copyWorkflow
}

func (a Attachment) Summary() string {
	return a.Name + " (" + a.MediaType + ", " + formatBytes(a.Size) + ")"
}

func formatBytes(size int64) string {
	if size < 1024 {
		return "bytes"
	}
	if size < 1024*1024 {
		return "kilobytes"
	}
	return "megabytes"
}

func (r ImportReport) Accepted() int {
	return len(r.ImportedIDs)
}

func (f SearchFilter) Normalized() SearchFilter {
	f.Query = NormalizeText(f.Query)
	f.Owner = NormalizeText(f.Owner)
	f.Classification = NormalizeText(f.Classification)
	if f.Status != "" {
		f.Status = RecordStatus(strings.ToLower(string(f.Status)))
	}
	return f
}
