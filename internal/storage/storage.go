package storage

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/22Fariz22/shorturl/internal/entity"
)

type MemoryStorage interface {
	Get(key string) (entity.URL, bool)
	Insert(key, value string, cook string, deleted bool) (string, error)
	GetAllStorageURL(string2 string) []map[string]string
	DeleteStorage([]string, string) error
}

///переделать в нормальную структуру
type memoryStorage struct {
	storage map[string]entity.URL // список мап sortURL:entity.URL
	mutex   sync.RWMutex
}

func (m *memoryStorage) DeleteStorage(listShorts []string, cookies string) error {
	log.Print("del in stor")

	for _, v := range listShorts {
		fmt.Println("v", v)
		for k := range m.storage {
			if m.storage[k].ID == v && m.storage[k].Cookies == cookies {
				m.storage[k] = entity.URL{
					Cookies:       cookies,
					ID:            v,
					LongURL:       k,
					CorrelationID: m.storage[k].CorrelationID,
					Deleted:       true,
				}
			}
		}
	}
	return nil
}

func (m *memoryStorage) GetAllStorageURL(cook string) []map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	list := make([]map[string]string, 1)

	for i, ok := range m.storage {
		if ok.Cookies == cook {
			mp := make(map[string]string)
			mp[m.storage[i].ID] = m.storage[i].LongURL
			list = append(list, mp)
		}
	}
	return list
}

func (m *memoryStorage) Get(key string) (entity.URL, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, x := range m.storage {
		if x.ID == key {
			fmt.Println("get x in storage", x)
			return x, true
		}
	}
	return entity.URL{}, false
}

func (m *memoryStorage) Insert(key string, value string, cook string, deleted bool) (string, error) {
	var ErrAlreadyExists = errors.New("this URL already exists")

	url := &entity.URL{
		Cookies: cook,
		ID:      key,
		LongURL: value,
		Deleted: deleted,
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.storage[value]
	if !ok {
		m.storage[value] = *url
		return "", nil
	}
	fmt.Println("long in storage", v.LongURL)
	fmt.Println("su in storage", v.ID)
	return v.ID, ErrAlreadyExists
}

func New() MemoryStorage {
	return &memoryStorage{
		storage: map[string]entity.URL{},
	}
}