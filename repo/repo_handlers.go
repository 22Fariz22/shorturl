package repo

import (
	"22Fariz22/shorturl/handler/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type JSONModel struct {
	Count int               `json:"count"`
	URL   map[string]string `json:"url"`
}

type AllJSONModels struct {
	AllUrls []*JSONModel
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(cnt int, urlMap map[string]string) error {
	cfg := config.NewConnectorConfig()
	newURL := &JSONModel{}
	newURL.URL = urlMap
	newURL.Count = cnt
	fmt.Println("p.file", p.file)
	b, err := ioutil.ReadAll(p.file)
	if err != nil {
		fmt.Println("ReadAll")
		log.Fatal(err)
	}
	defer p.file.Close()

	var alUrls AllJSONModels

	_ = json.Unmarshal(b, &alUrls.AllUrls)

	alUrls.AllUrls = append(alUrls.AllUrls[len(alUrls.AllUrls):], newURL)
	newURLBytes, err := json.MarshalIndent(&alUrls.AllUrls, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(cfg.FileStoragePath, newURLBytes, 0666)
	return nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}

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

func (c *consumer) Close() error {
	return c.file.Close()
}
