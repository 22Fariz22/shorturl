package memory

import (
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/storage"
)

type inMemoryRepository struct {
	memoryStorage storage.MemoryStorage
}

func New() repository.Repository {
	st := storage.New()
	return &inMemoryRepository{
		memoryStorage: st,
	}
}

func (f *inMemoryRepository) SaveURL(shortID string, longURL string) error {

	return nil
}

func (f *inMemoryRepository) GetURL(shortID string) (string, error) {
	return "", nil
}
