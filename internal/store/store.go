package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"example.com/energycore/internal/model"
	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("entity not found")

var bucketNames = [][]byte{
	[]byte("records"),
	[]byte("events"),
	[]byte("workflows"),
	[]byte("attachments"),
	[]byte("metadata"),
}

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepathDir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	store := &Store{db: db, path: path}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func filepathDir(path string) string {
	for index := len(path) - 1; index >= 0; index-- {
		if path[index] == '/' {
			if index == 0 {
				return "/"
			}
			return path[:index]
		}
	}
	return "."
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string {
	return s.path
}

func encode(value any) ([]byte, error) {
	return json.Marshal(value)
}

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	return json.Unmarshal(data, target)
}

func putJSON(bucket *bbolt.Bucket, key string, value any) error {
	data, err := encode(value)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), data)
}

func getJSON(bucket *bbolt.Bucket, key string, target any) error {
	data := bucket.Get([]byte(key))
	if data == nil {
		return ErrNotFound
	}
	return decode(data, target)
}

func (s *Store) Delete(bucketName []byte, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket.Get([]byte(key)) == nil {
			return ErrNotFound
		}
		return bucket.Delete([]byte(key))
	})
}

func (s *Store) withUpdate(fn func(*bbolt.Tx) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(fn)
}

func (s *Store) withView(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.View(fn)
}

func (s *Store) PutRecord(record model.Record) error {
	return s.withUpdate(func(tx *bbolt.Tx) error { return putJSON(tx.Bucket([]byte("records")), record.ID, record) })
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.withView(func(tx *bbolt.Tx) error { return getJSON(tx.Bucket([]byte("records")), id, &record) })
	return record, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.Delete([]byte("records"), id)
}

func (s *Store) ListRecords() ([]model.Record, error) {
	items := []model.Record{}
	err := s.withView(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, value []byte) error {
			var record model.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			items = append(items, record)
			return nil
		})
	})
	return items, err
}
