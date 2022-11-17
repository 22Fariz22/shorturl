package repo

import (
	"22Fariz22/shorturl/handler/config"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type JSONModel struct {
	count int               `json:"count"`
	urls  map[string]string `json:"urls"`
}
type AllJSONModels struct {
	allUrls []*JSONModel
}

type CreateShortURLRequestArray struct {
	URLs []CreateShortURLRequestArray `json:"urls"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteEvent(cnt int, urlMap map[string]string) error {
	cfg := config.NewConnectorConfig()
	newURL := &JSONModel{}
	newURL.urls = urlMap
	newURL.count = cnt

	b, _ := ioutil.ReadAll(p.file)
	var alUrls AllJSONModels
	err := json.Unmarshal(b, &alUrls.allUrls)
	if err != nil {
		log.Fatal(err)
	}

	alUrls.allUrls = append(alUrls.allUrls, newURL)
	newURLBytes, err := json.MarshalIndent(&alUrls.allUrls, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(cfg.FileStoragePath, newURLBytes, 0666)
	return nil
}

func (p *producer) Close() error {
	return p.file.Close()
}

//********
//дальше код для того чтобы приложение при перезапуске прочитала или востановила ранее сокращенные урлы
type consumer struct {
	file    *os.File
	Decoder *json.Decoder
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		Decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadEvent() (*CreateShortURLRequest, error) {
	event := &CreateShortURLRequest{}
	if err := c.Decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

//функция для востановления списка urls
func (c *consumer) RecoverEvents() {

}

func (c *consumer) Close() error {
	return c.file.Close()
}
