package memory

import (
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/storage"
)

type inMemoryRepository struct {
	memoryStorage storage.MemoryStorage
}

func (m *inMemoryRepository) Init() error {
	//TODO implement me
	panic("implement me")
}

func New() repository.Repository {
	st := storage.New()
	return &inMemoryRepository{
		memoryStorage: st,
	}
}

func (m *inMemoryRepository) SaveURL(shortID string, longURL string, cook string) error {
	m.memoryStorage.Insert(shortID, longURL, cook)
	return nil
}

func (m *inMemoryRepository) GetURL(shortID string) (string, bool) {
	v, ok := m.memoryStorage.Get(shortID)

	return v, ok
}

func (m *inMemoryRepository) GetAll(cook string) []map[string]string {
	return m.memoryStorage.GetAllStorageURL(cook)
}
