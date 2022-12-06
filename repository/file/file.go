package file

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/model"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/storage"
)

type inFileRepository struct {
	file          io.ReadWriteCloser
	memoryStorage storage.MemoryStorage
	reader        *bufio.Reader
}

type Consumer struct {
	File *os.File
	//Cfg    *config.Config
	reader *bufio.Reader
}

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

func New(cfg *config.Config) repository.Repository {
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

func (f *inFileRepository) Init() error {
	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		txt := scanner.Text()
		var u model.URL
		err := json.Unmarshal([]byte(txt), &u)
		if err != nil {
			return err
		}
		f.memoryStorage.Insert(u.ID, u.LongURL, u.Cookies)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	return nil
}

func (f *inFileRepository) SaveURL(shortID string, longURL string, cook string) error {
	url := &model.URL{
		Cookies: cook,
		ID:      shortID,
		LongURL: longURL,
	}
	data, err := json.Marshal(url)
	if err != nil {
		log.Println(err)
		return err
	}
	f.file.Write(data)
	f.file.Write([]byte("\n"))
	f.memoryStorage.Insert(shortID, longURL, cook)
	return nil
}

func (f *inFileRepository) GetURL(shortID string) (string, bool) {
	v, ok := f.memoryStorage.Get(shortID)
	return v, ok
}

func (f *inFileRepository) GetAll(cook string) []map[string]string {
	return f.memoryStorage.GetAllStorageURL(cook)

}
func (f *inFileRepository) Ping() int {
	return http.StatusInternalServerError
}
