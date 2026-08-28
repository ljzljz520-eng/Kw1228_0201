package store

import (
	"example.com/energycore/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) PutAttachment(attachment model.Attachment) error {
	return s.withUpdate(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket([]byte("attachments")), attachment.ID, attachment)
	})
}

func (s *Store) ListAttachments(recordID string) ([]model.Attachment, error) {
	attachments := []model.Attachment{}
	err := s.withView(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("attachments")).ForEach(func(_, value []byte) error {
			var attachment model.Attachment
			if err := decode(value, &attachment); err != nil {
				return err
			}
			if recordID == "" || attachment.RecordID == recordID {
				attachments = append(attachments, attachment)
			}
			return nil
		})
	})
	return attachments, err
}
