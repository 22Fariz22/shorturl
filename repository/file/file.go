package file

import (
	"bufio"
	"encoding/json"
	"github.com/22Fariz22/shorturl/handler/config"
	"github.com/22Fariz22/shorturl/model"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/storage"
	"io"
	"log"
	"os"
)

type inFileRepository struct {
	file          io.ReadWriteCloser
	memoryStorage storage.MemoryStorage
	reader        *bufio.Reader
}

type Consumer struct {
	File   *os.File
	Cfg    *config.Config
	reader *bufio.Reader
}

func NewConsumer() (*Consumer, error) {
	file, err := os.OpenFile(config.DefaultFileStoragePath, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		File:   file,
		Cfg:    config.NewConnectorConfig(),
		reader: bufio.NewReader(file),
	}, nil
}

func New() repository.Repository {
	st := storage.New()

	consumer, err := NewConsumer()
	if err != nil {
		log.Fatal(err)
	}

	return &inFileRepository{
		file:          consumer.File,
		memoryStorage: st,
		reader:        consumer.reader,
	}
}

func (f *inFileRepository) Init() error {
	scanner := bufio.NewScanner(f.file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		txt := scanner.Text()
		var u model.URL
		err := json.Unmarshal([]byte(txt), &u)
		if err != nil {
			return err
		}
		f.memoryStorage.Insert(u.ID, u.LongURL)
		log.Println(u)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (f *inFileRepository) SaveURL(shortID string, longURL string) error {
	url := &model.URL{
		ID:      shortID,
		LongURL: longURL,
	}
	data, err := json.Marshal(url)
	if err != nil {
		log.Println(err)
		return err
	}
	f.file.Write([]byte("\n"))
	f.file.Write(data)
	f.memoryStorage.Insert(shortID, longURL)
	return nil
}

func (f *inFileRepository) GetURL(shortID string) (string, bool) {
	v, ok := f.memoryStorage.Get(shortID)
	return v, ok
}
