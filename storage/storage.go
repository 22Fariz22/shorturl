package storage

import (
	"sync"
)

type MemoryStorage interface {
	Get(key string) (string, bool)
	Insert(key, value string, cook string) error
	GetAllStorageURL(string2 string) []map[string]string
}

type memoryStorage struct {
	storage        map[string]string
	storageCookies map[string][]map[string]string // like as map([cookies]map[shortURL][longURL]
	mutex          sync.RWMutex
}

func (m *memoryStorage) GetAllStorageURL(cook string) []map[string]string {
	return m.storageCookies[cook]
}

func (m *memoryStorage) Get(key string) (string, bool) {
	v, ok := m.storage[key]
	return v, ok
}

func (m *memoryStorage) Insert(key, value string, cook string) error {
	m.storage[key] = value

	aMap := map[string]string{key: value}

	m.storageCookies[cook] = append(m.storageCookies[cook], aMap)

	return nil
}

func New() MemoryStorage {
	return &memoryStorage{
		storage:        make(map[string]string),
		storageCookies: make(map[string][]map[string]string),
	}
}
