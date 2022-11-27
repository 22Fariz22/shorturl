package storage

import "sync"

type MemoryStorage interface {
	Get(key string) (string, error)
	Insert(key, value string) error
}

type memoryStorage struct {
	storage map[string]string
	mutex   sync.RWMutex
}
