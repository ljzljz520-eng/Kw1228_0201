package store

import (
	"example.com/energycore/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) PutEvent(event model.AuditEvent) error {
	return s.withUpdate(func(tx *bbolt.Tx) error {
		return putJSON(tx.Bucket([]byte("events")), event.ID, event)
	})
}

func (s *Store) ListEvents(recordID string) ([]model.AuditEvent, error) {
	events := []model.AuditEvent{}
	err := s.withView(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("events")).ForEach(func(_, value []byte) error {
			var event model.AuditEvent
			if err := decode(value, &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				events = append(events, event)
			}
			return nil
		})
	})
	return events, err
}

func (s *Store) DeleteEvents(recordID string) error {
	return s.withUpdate(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("events"))
		keys := [][]byte{}
		if err := bucket.ForEach(func(key, value []byte) error {
			var event model.AuditEvent
			if err := decode(value, &event); err != nil {
				return err
			}
			if event.RecordID == recordID {
				copyKey := append([]byte(nil), key...)
				keys = append(keys, copyKey)
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}
