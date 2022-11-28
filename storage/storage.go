package storage

import (
	"sync"
)

type MemoryStorage interface {
	Get(key string) (string, error)
	Insert(key, value string) error
}

type memoryStorage struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func (m *memoryStorage) Get(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *memoryStorage) Insert(key, value string) error {
	//TODO implement me
	panic("implement me")
}

func New() MemoryStorage {
	return &memoryStorage{}
}
