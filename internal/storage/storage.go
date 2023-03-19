// Package storage юзекейс для мемори сторада
package storage

import (
	"errors"
	"log"
	"sync"

	"github.com/22Fariz22/shorturl/internal/entity"
)

// MemoryStorage интерфейс для стоража инмемори
type MemoryStorage interface {
	Get(key string) (entity.URL, bool)
	Insert(key, value string, cook string, deleted bool) (string, error)
	GetAllStorageURL(string2 string) []map[string]string
	DeleteStorage([]string, string) error
}

// /memoryStorage структура для стоража инмемори
type memoryStorage struct {
	storage map[string]entity.URL // список мап sortURL:entity.URL
	mutex   sync.RWMutex
}

// DeleteStorage удаление записи
func (m *memoryStorage) DeleteStorage(listShorts []string, cookies string) error {
	log.Print("del in stor")

	for _, v := range listShorts {
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

// GetAllStorageURL получить все записи
func (m *memoryStorage) GetAllStorageURL(cook string) []map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	list := make([]map[string]string, 0)

	for i, ok := range m.storage {
		if ok.Cookies == cook {
			mp := make(map[string]string)
			mp[m.storage[i].ID] = m.storage[i].LongURL
			list = append(list, mp)
		}
	}
	return list
}

// Get получить запись
func (m *memoryStorage) Get(key string) (entity.URL, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, x := range m.storage {
		if x.ID == key {
			return x, true
		}
	}
	return entity.URL{}, false
}

// Insert вставить запись
func (m *memoryStorage) Insert(key string, value string, cook string, deleted bool) (string, error) {
	//ErrAlreadyExists вывод ошибки существования
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
	return v.ID, ErrAlreadyExists // если такого еще нет в мапе,то ничего не вернет. а если есть, то вернет shorturl.
}

// New создание структуры для инмемори типа
func New() MemoryStorage {
	return &memoryStorage{
		storage: map[string]entity.URL{},
	}
}
