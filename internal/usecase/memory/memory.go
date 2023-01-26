package memory

import (
	"context"
	"fmt"
	"github.com/22Fariz22/shorturl/internal/storage"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"log"

	"github.com/22Fariz22/shorturl/internal/entity"
)

type inMemoryRepository struct {
	memoryStorage storage.MemoryStorage
}

func (m *inMemoryRepository) Init() error {
	return nil
}

func New() usecase.Repository {
	st := storage.New()
	return &inMemoryRepository{
		memoryStorage: st,
	}
}

func (m *inMemoryRepository) SaveURL(ctx context.Context, shortID string, longURL string, cook string) (string, error) {
	s, err := m.memoryStorage.Insert(shortID, longURL, cook, false)
	fmt.Println("su in memory", s)
	return s, err
}

func (m *inMemoryRepository) GetURL(ctx context.Context, shortID string) (entity.URL, bool) {
	v, ok := m.memoryStorage.Get(shortID)
	fmt.Println("ok in mem Geturl:", ok)
	return v, ok
}

func (m *inMemoryRepository) GetAll(ctx context.Context, cook string) ([]map[string]string, error) {
	return m.memoryStorage.GetAllStorageURL(cook), nil
}

func (m *inMemoryRepository) Ping(ctx context.Context) error {
	return nil
}

func (m *inMemoryRepository) RepoBatch(ctx context.Context, cook string, batchList []entity.PackReq) error {
	for i := range batchList {
		url := &entity.URL{
			ID:      batchList[i].ShortURL,
			LongURL: batchList[i].OriginalURL,
		}
		m.memoryStorage.Insert(url.ID, url.LongURL, cook, false)
	}
	return nil
}

func (m *inMemoryRepository) Delete(list []string, cookie string) error {
	log.Print("del in mem")
	m.memoryStorage.DeleteStorage(list, cookie)
	return nil
}
