// Package storage юзекейс для мемори сторада
package storage

import (
	"context"
	"errors"
	"github.com/22Fariz22/shorturl/pkg/logger"
	"sync"

	"github.com/22Fariz22/shorturl/internal/entity"
)

// MemoryStorage интерфейс для стоража инмемори
type MemoryStorage interface {
	Get(l logger.Interface, key string) (entity.URL, bool)
	Insert(l logger.Interface, key, value string, cook string, deleted bool) (string, error)
	GetAllStorageURL(l logger.Interface, string2 string) []map[string]string
	DeleteStorage(logger.Interface, []string, string) error
	Stats(ctx context.Context, l logger.Interface) (int, int, error)
}

// /memoryStorage структура для стоража инмемори
type memoryStorage struct {
	storage map[string]entity.URL // список мап sortURL:entity.URL
	mutex   sync.RWMutex
}

// DeleteStorage удаление записи
func (m *memoryStorage) DeleteStorage(l logger.Interface, listShorts []string, cookies string) error {
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
func (m *memoryStorage) GetAllStorageURL(l logger.Interface, cook string) []map[string]string {
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
func (m *memoryStorage) Get(l logger.Interface, key string) (entity.URL, bool) {
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
func (m *memoryStorage) Insert(l logger.Interface, key string, value string, cook string, deleted bool) (string, error) {
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

// Stats
func (m *memoryStorage) Stats(ctx context.Context, l logger.Interface) (int, int, error) {
	mpUsers := make(map[string]int)
	countUrls := 0

	for _, x := range m.storage {
		mpUsers[x.Cookies] += 1
		countUrls += 1
	}

	return countUrls, len(mpUsers), nil
}

// New создание структуры для инмемори типа
func New() MemoryStorage {
	return &memoryStorage{
		storage: map[string]entity.URL{},
	}
}
