package storage

import (
	"errors"
	"fmt"
	"github.com/22Fariz22/shorturl/model"
	"sync"
)

type MemoryStorage interface {
	Get(key string) (model.URL, bool)
	Insert(key, value string, cook string, deleted bool) (string, error)
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

type storList struct {
}

func (m *memoryStorage) DeleteStorage(listShorts []string, cookies string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, v := range listShorts {
		//m.mutex.Lock()
		if url, ok := m.storage[v]; ok {
			url.Deleted = true
			m.storage[v] = url
		}
		//m.mutex.RUnlock()

	}
	return nil
}

func (m *memoryStorage) GetAllStorageURL(cook string) []map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	list := make([]map[string]string, 1)

	for i, ok := range m.storage { //i = shortURL ok=model.URL
		if ok.Cookies == cook {
			//m.mutex.Lock()
			mp := make(map[string]string)
			mp[m.storage[i].ID] = m.storage[i].LongURL
			list = append(list, mp)
			//m.mutex.RUnlock()

		}
	}
	return list
}

func (m *memoryStorage) Get(key string) (model.URL, bool) {
	//m.storage это список
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	v, ok := m.storage[key]
	fmt.Println(v, ok)
	if !ok {
		return v, false
	}
	return v, ok

}

func (m *memoryStorage) Insert(key string, value string, cook string, deleted bool) (string, error) {
	//m.storage[key] = value
	//aMap := map[string]string{key: value}
	//m.storageCookies[cook] = append(m.storageCookies[cook], aMap)

	//m.mutex.RLock()
	//defer m.mutex.RUnlock()
	var ErrAlreadyExists = errors.New("this URL already exists")

	url := &model.URL{
		Cookies: cook,
		ID:      key,
		LongURL: value,
		Deleted: deleted,
	}

	v, ok := m.storage[value]
	if !ok {
		m.storage[value] = *url
		return "", nil
	} else {
		fmt.Println("long in storage", v.LongURL)
		fmt.Println("su in storage", v.ID)
		return v.ID, ErrAlreadyExists
	}

	//mp := map[string]model.URL{key: url}
	//
	//m.storage = append(m.storage, mp)
	return "", nil
}

func New() MemoryStorage {
	return &memoryStorage{
		storage: map[string]model.URL{},
		//storage:        make(map[string]string),
		storageCookies: make(map[string][]map[string]string),
	}
}
