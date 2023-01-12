package storage

import (
	"errors"
	"fmt"
	"github.com/22Fariz22/shorturl/model"
	"log"
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
	//storageCookies map[string][]map[string]string // like as map([cookies]map[shortURL][longURL]
	mutex sync.RWMutex
}

type storList struct {
}

func (m *memoryStorage) DeleteStorage(listShorts []string, cookies string) error {
	log.Print("del in stor")
	//m.mutex.RLock()
	//defer m.mutex.RUnlock()

	for _, v := range listShorts {
		//m.mutex.Lock()
		//if url, ok := m.storage[v]; ok {
		//	fmt.Println("in stor Del v in range:", url, ok)
		//	url.Deleted = true
		//	m.storage[v] = url
		//}
		fmt.Println("v", v)
		for k := range m.storage {
			if m.storage[k].ID == v && m.storage[k].Cookies == cookies {
				//delete(m.storage, k)
				m.mutex.RLock()
				m.storage[k] = model.URL{
					Cookies:       cookies,
					ID:            v,
					LongURL:       k,
					CorrelationID: m.storage[k].CorrelationID,
					Deleted:       true,
				}
				m.mutex.RUnlock()
			}
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
	//m.storage это мапа
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	//v, ok := m.storage[key]
	//fmt.Println(v, ok)
	//if !ok {
	//	return v, false
	//}
	//return v, ok

	for _, x := range m.storage {
		if x.ID == key {
			fmt.Println("get x in storage", x)
			return x, true
		}
	}
	return model.URL{}, false

}

func (m *memoryStorage) Insert(key string, value string, cook string, deleted bool) (string, error) {
	//m.storage[key] = value
	//aMap := map[string]string{key: value}
	//m.storageCookies[cook] = append(m.storageCookies[cook], aMap)

	var ErrAlreadyExists = errors.New("this URL already exists")

	url := &model.URL{
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
		storage: map[string]model.URL{},
		//storage:        make(map[string]string),
		//storageCookies: make(map[string][]map[string]string),
	}
}
