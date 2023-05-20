// Package memory пакет для работы инмемори
package memory

import (
	"context"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/storage"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"github.com/22Fariz22/shorturl/pkg/logger"
)

// inMemoryRepository структура инмемори
type inMemoryRepository struct {
	memoryStorage storage.MemoryStorage
}

// Init инициализаци инмемори
func (m *inMemoryRepository) Init(l logger.Interface) error {
	return nil
}

// New создание инмемори
func New() usecase.Repository {
	st := storage.New()
	return &inMemoryRepository{
		memoryStorage: st,
	}
}

// SaveURL сохранение в инмемори
func (m *inMemoryRepository) SaveURL(ctx context.Context, l logger.Interface, shortID string, longURL string, cook string) (string, error) {
	s, err := m.memoryStorage.Insert(l, shortID, longURL, cook, false)
	return s, err
}

// GetURL получить запись из инмемори
func (m *inMemoryRepository) GetURL(ctx context.Context, l logger.Interface, shortID string) (entity.URL, bool) {
	v, ok := m.memoryStorage.Get(l, shortID)
	return v, ok
}

// GetAll получить записи из инмемори
func (m *inMemoryRepository) GetAll(ctx context.Context, l logger.Interface, cook string) ([]map[string]string, error) {
	return m.memoryStorage.GetAllStorageURL(l, cook), nil
}

// Ping заглушка
func (m *inMemoryRepository) Ping(ctx context.Context, l logger.Interface) error {
	return nil
}

// RepoBatch создание списка записей в инмемори
func (m *inMemoryRepository) RepoBatch(ctx context.Context, l logger.Interface, cook string, batchList []entity.PackReq) error {
	for i := range batchList {
		url := &entity.URL{
			ID:      batchList[i].ShortURL,
			LongURL: batchList[i].OriginalURL,
		}
		m.memoryStorage.Insert(l, url.ID, url.LongURL, cook, false)
	}
	return nil
}

// Delete удаление записи из инмемори
func (m *inMemoryRepository) Delete(l logger.Interface, list []string, cookie string) error {
	m.memoryStorage.DeleteStorage(l, list, cookie)
	return nil
}

func (m *inMemoryRepository) Stats(ctx context.Context, l logger.Interface) (int, int, error) {
	return m.memoryStorage.Stats(ctx, l)
}
