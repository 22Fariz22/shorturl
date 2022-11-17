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
	Count int               `json:"count"`
	Url   map[string]string `json:"url"`
}
type AllJSONModels struct {
	AllUrls []*JSONModel
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
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
	newURL.Url = urlMap
	newURL.Count = cnt

	//open file
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var alUrls AllJSONModels

	err = json.Unmarshal(b, &alUrls.AllUrls)
	if err != nil {
		log.Fatal(err)
	}

	alUrls.AllUrls = append(alUrls.AllUrls, newURL)
	newURLBytes, err := json.MarshalIndent(&alUrls.AllUrls, "", " ")
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
