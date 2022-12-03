package memory

import (
	"github.com/22Fariz22/shorturl/model"
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

func (m *inMemoryRepository) SaveURL(shortID string, longURL string) error {
	m.memoryStorage.Insert(shortID, longURL)
	return nil
}

func (m *inMemoryRepository) GetURL(shortID string) (string, bool) {
	v, ok := m.memoryStorage.Get(shortID)

	return v, ok
}

func (m *inMemoryRepository) GetAll() []model.URL {
	return m.memoryStorage.GetAllStorageURL()
}
