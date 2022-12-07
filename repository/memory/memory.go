package memory

import (
	"context"

	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/storage"
)

type inMemoryRepository struct {
	memoryStorage storage.MemoryStorage
}

func (m *inMemoryRepository) Init() error {
	return nil
}

func New() repository.Repository {
	st := storage.New()
	return &inMemoryRepository{
		memoryStorage: st,
	}
}

func (m *inMemoryRepository) SaveURL(ctx context.Context, shortID string, longURL string, cook string) error {
	m.memoryStorage.Insert(shortID, longURL, cook)
	return nil
}

func (m *inMemoryRepository) GetURL(ctx context.Context, shortID string) (string, bool) {
	v, ok := m.memoryStorage.Get(shortID)

	return v, ok
}

func (m *inMemoryRepository) GetAll(ctx context.Context, cook string) ([]map[string]string, error) {
	return m.memoryStorage.GetAllStorageURL(cook), nil
}

func (m *inMemoryRepository) Ping(ctx context.Context) error {
	return nil
}
