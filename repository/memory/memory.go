package memory

import (
	"github.com/22Fariz22/shorturl/storage"
)

type inMemoryRepository struct {
	memoryStorage storage.MemoryStorage
}

func New() *inMemoryRepository {
	var memoryStorage storage.MemoryStorage
	return &inMemoryRepository{
		memoryStorage: memoryStorage,
	}
}

func (f *inMemoryRepository) SaveUrl(shortID string, longURL string) error {

	return nil
}

func (f *inMemoryRepository) GetURL(shortID string) (string, error) {

	return "", nil

}
