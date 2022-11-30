package storage

import (
	"sync"
)

type MemoryStorage interface {
	Get(key string) (string, bool)
	Insert(key, value string) error
}

type memoryStorage struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func (m *memoryStorage) Get(key string) (string, bool) {
	v, ok := m.storage[key]
	return v, ok
}

func (m *memoryStorage) Insert(key, value string) error {
	m.storage[key] = value
	return nil
}

func New() MemoryStorage {
	return &memoryStorage{
		storage: make(map[string]string),
	}
}
