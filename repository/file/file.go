package file

import (
	"bufio"
	"github.com/22Fariz22/shorturl/handler/config"
	"github.com/22Fariz22/shorturl/model"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/storage"
	"io"
	"log"
	"os"
)

type inFileRepository struct {
	file          io.WriteCloser
	memoryStorage storage.MemoryStorage
	reader        *bufio.Reader
}

type Consumer struct {
	File   *os.File
	Cfg    *config.Config
	reader *bufio.Reader
}

func NewConsumer() (*Consumer, error) {
	file, err := os.OpenFile(config.DefaultFileStoragePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (i *inFileRepository) SaveURL(shortID string, longURL string) error {
	url := &model.URL{
		ID:      shortID,
		LongURL: longURL,
	}

	panic("implement me")
}

func (i *inFileRepository) GetURL(shortID string) (string, error) {
	//TODO implement me
	panic("implement me")
}
