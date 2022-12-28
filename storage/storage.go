package storage

import (
	"fmt"
	"github.com/22Fariz22/shorturl/model"
	"sync"
)

type MemoryStorage interface {
	Get(key string) (model.URL, bool)
	Insert(key, value string, cook string, deleted bool) error
	GetAllStorageURL(string2 string) []map[string]string
	DeleteStorage([]string, string) error
}

///переделать в нормальную структуру
type memoryStorage struct {
	storage     map[string]model.URL // список мап sortURL:model.URL
	storageList []map[string]model.URL
	//storage        map[string]string  //old
	storageCookies map[string][]map[string]string // like as map([cookies]map[shortURL][longURL]
	mutex          sync.RWMutex
}

func (m *memoryStorage) DeleteStorage(listShorts []string, cookies string) error {
	for _, v := range listShorts {
		if url, ok := m.storage[v]; ok {
			url.Deleted = true
			m.storage[v] = url
		}
	}
	return nil
}

func (m *memoryStorage) GetAllStorageURL(cook string) []map[string]string {
	return m.storageCookies[cook]
}

func (m *memoryStorage) Get(key string) (model.URL, bool) {
	//m.storage это список

	v, ok := m.storage[key]
	fmt.Println(v, ok)
	if !ok {
		return v, false
	}
	return v, ok

}

func (m *memoryStorage) Insert(key string, value string, cook string, deleted bool) error { //u.ID, u.LongURL, u.Cookies,u.Deleted
	//m.storage[key] = value
	//aMap := map[string]string{key: value}
	//m.storageCookies[cook] = append(m.storageCookies[cook], aMap)

	url := model.URL{
		Cookies: cook,
		ID:      key,
		LongURL: value,
		Deleted: deleted,
	}
	m.storage[key] = url

	//mp := map[string]model.URL{key: url}
	//
	//m.storage = append(m.storage, mp)

	return nil
}

func New() MemoryStorage {
	return &memoryStorage{
		storage: map[string]model.URL{},
		//storage:        make(map[string]string),
		storageCookies: make(map[string][]map[string]string),
	}
}
