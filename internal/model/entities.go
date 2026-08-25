package model

type RecordStatus string

const (
	StatusDraft    RecordStatus = "draft"
	StatusInReview RecordStatus = "in_review"
	StatusApproved RecordStatus = "approved"
	StatusRejected RecordStatus = "rejected"
	StatusArchived RecordStatus = "archived"
)

type ReviewDecision string

const (
	DecisionApprove ReviewDecision = "approve"
	DecisionReject  ReviewDecision = "reject"
	DecisionRequest ReviewDecision = "request_changes"
)

type Record struct {
	ID             string       `json:"id"`
	Code           string       `json:"code"`
	Name           string       `json:"name"`
	Status         RecordStatus `json:"status"`
	Description    string       `json:"description"`
	Owner          string       `json:"owner"`
	Classification string       `json:"classification"`
	ReviewNote     string       `json:"review_note"`
	WithdrawnNote  string       `json:"withdrawn_note"`
	Revision       int          `json:"revision"`
	CreatedSeq     int64        `json:"created_seq"`
	UpdatedSeq     int64        `json:"updated_seq"`
	ArchivedSeq    int64        `json:"archived_seq"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Message  string `json:"message"`
	Sequence int64  `json:"sequence"`
	Revision int    `json:"revision"`
}

type Workflow struct {
	ID         string   `json:"id"`
	RecordID   string   `json:"record_id"`
	Name       string   `json:"name"`
	State      string   `json:"state"`
	Steps      []string `json:"steps"`
	Completed  []string `json:"completed"`
	StartedSeq int64    `json:"started_seq"`
	UpdatedSeq int64    `json:"updated_seq"`
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	Sequence  int64  `json:"sequence"`
}

type ReviewResult struct {
	RecordID      string         `json:"record_id"`
	Decision      ReviewDecision `json:"decision"`
	Note          string         `json:"note"`
	Reviewer      string         `json:"reviewer"`
	Revision      int            `json:"revision"`
	Sequence      int64          `json:"sequence"`
	WithdrawnNote string         `json:"withdrawn_note"`
}

type SearchFilter struct {
	Query           string
	Status          RecordStatus
	Owner           string
	Classification  string
	IncludeArchived bool
}

type ImportRow struct {
	Code           string
	Name           string
	Description    string
	Owner          string
	Classification string
}

type ImportReport struct {
	ImportedIDs []string `json:"imported_ids"`
	Rejected    []string `json:"rejected"`
	Total       int      `json:"total"`
}
