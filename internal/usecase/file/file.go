// Package file пакет для работы с файлом
package file

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/storage"
	"github.com/22Fariz22/shorturl/internal/usecase"

	"github.com/22Fariz22/shorturl/internal/entity"
)

//inFileRepository структура для сторада инфайл
type inFileRepository struct {
	file          io.ReadWriteCloser
	memoryStorage storage.MemoryStorage
	reader        *bufio.Reader
}

//Delete удаление из инфайла
func (f *inFileRepository) Delete(list []string, cookie string) error {
	f.memoryStorage.DeleteStorage(list, cookie)
	return nil
}

// создание списка записей  в инфайле
func (f *inFileRepository) RepoBatch(ctx context.Context, cook string, batchList []entity.PackReq) error {
	// [{1 http://mail.ru 0ATJMCH} {2 http://ya.ru 3DXH7RG} {3 http://google.ru VGGFB0D}]
	for i := range batchList {
		url := &entity.URL{
			Cookies: cook,
			ID:      batchList[i].ShortURL,
			LongURL: batchList[i].OriginalURL,
		}
		data, err := json.Marshal(url)
		if err != nil {
			log.Println(err)
			return err
		}
		f.file.Write(data)
		f.file.Write([]byte("\n"))
		f.memoryStorage.Insert(url.ID, url.LongURL, cook, false)
	}
	return nil
}

//Consumer структура консьюмера
type Consumer struct {
	File   *os.File
	reader *bufio.Reader
}

// NewConsumer  создание консьюмера
func NewConsumer(cfg config.Config) (*Consumer, error) {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		File:   file,
		reader: bufio.NewReader(file),
	}, nil
}

// New инициализация консьюмера
func New(cfg *config.Config) usecase.Repository {
	st := storage.New()

	consumer, err := NewConsumer(*cfg)
	if err != nil {
		log.Println(err)
	}

	return &inFileRepository{
		file:          consumer.File,
		memoryStorage: st,
		reader:        consumer.reader,
	}
}

// Init инициализация консьмера
func (f *inFileRepository) Init() error {
	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		txt := scanner.Text()
		var u entity.URL
		err := json.Unmarshal([]byte(txt), &u)
		if err != nil {
			return err
		}
		f.memoryStorage.Insert(u.ID, u.LongURL, u.Cookies, u.Deleted)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	return nil
}

//SaveURLсохранить запись в файле
func (f *inFileRepository) SaveURL(ctx context.Context, shortID string, longURL string, cook string) (string, error) {
	url := &entity.URL{
		Cookies: cook,
		ID:      shortID,
		LongURL: longURL,
		Deleted: false,
	}
	data, err := json.Marshal(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	f.file.Write(data)
	f.file.Write([]byte("\n"))
	f.memoryStorage.Insert(shortID, longURL, cook, false) // add to inMemory
	return "", nil
}

//GetURL поулсить запись из файла
func (f *inFileRepository) GetURL(ctx context.Context, shortID string) (entity.URL, bool) {
	v, ok := f.memoryStorage.Get(shortID)
	return v, ok
}

//GetAll получсить все записи из файла
func (f *inFileRepository) GetAll(ctx context.Context, cook string) ([]map[string]string, error) {
	return f.memoryStorage.GetAllStorageURL(cook), nil
}

//Ping заглушка метода Ping
func (f *inFileRepository) Ping(ctx context.Context) error {
	return nil
}
