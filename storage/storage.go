package storage

import (
	"github.com/22Fariz22/shorturl/model"
	"sync"
)

type MemoryStorage interface {
	Get(key string) (string, bool)
	Insert(key, value string) error
	GetAllStorageURL() []model.URL
}

type memoryStorage struct {
	storage map[string]string
	mutex   sync.RWMutex
}

func (m *memoryStorage) GetAllStorageURL() []model.URL {
	//fmt.Println(m.storage)
	resp := []model.URL{}
	for k, v := range m.storage {
		resp = append(resp, model.URL{
			ID:      k,
			LongURL: v,
		})
	}
	return resp
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
