package clock

import "fmt"

type IDs struct {
	sequence Clock
}

func NewIDs(sequence Clock) *IDs {
	return &IDs{sequence: sequence}
}

func (i *IDs) RecordID() string {
	return fmt.Sprintf("core-%06d", i.sequence.Next())
}

func (i *IDs) EventID() string {
	return fmt.Sprintf("event-%06d", i.sequence.Next())
}

func (i *IDs) WorkflowID() string {
	return fmt.Sprintf("workflow-%06d", i.sequence.Next())
}

func (i *IDs) AttachmentID() string {
	return fmt.Sprintf("attachment-%06d", i.sequence.Next())
}
