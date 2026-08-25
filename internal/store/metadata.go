package store

import (
	"errors"
	"strconv"

	"go.etcd.io/bbolt"
)

var ErrMetadataMissing = errors.New("metadata key is missing")

func (s *Store) SetMetadata(key, value string) error {
	if key == "" {
		return errors.New("metadata key is required")
	}
	return s.withUpdate(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("metadata")).Put([]byte(key), []byte(value))
	})
}

func (s *Store) GetMetadata(key string) (string, error) {
	var value string
	err := s.withView(func(tx *bbolt.Tx) error {
		data := tx.Bucket([]byte("metadata")).Get([]byte(key))
		if data == nil {
			return ErrMetadataMissing
		}
		value = string(data)
		return nil
	})
	return value, err
}

func (s *Store) DeleteMetadata(key string) error {
	return s.withUpdate(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("metadata"))
		if bucket.Get([]byte(key)) == nil {
			return ErrMetadataMissing
		}
		return bucket.Delete([]byte(key))
	})
}

func (s *Store) SetSequence(value int64) error {
	return s.SetMetadata("sequence", strconv.FormatInt(value, 10))
}

func (s *Store) Sequence() (int64, error) {
	value, err := s.GetMetadata("sequence")
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

func (s *Store) Health() error {
	return s.withView(func(tx *bbolt.Tx) error {
		if tx.Bucket([]byte("records")) == nil {
			return errors.New("records bucket missing")
		}
		return nil
	})
}
