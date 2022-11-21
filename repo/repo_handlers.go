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
	URL   map[string]string `json:"url"`
}

type AllJSONModels struct {
	AllUrls []*JSONModel //`json:"all_urls"`
}

type Producer struct {
	file *os.File
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}
	return &Producer{
		file: file,
	}, nil
}

func (p *Producer) WriteEvent(cnt int, urlMap map[string]string) error {
	cfg := config.NewConnectorConfig()
	newURL := &JSONModel{}
	newURL.URL = urlMap
	newURL.Count = cnt

	b, err := ioutil.ReadAll(p.file)
	if err != nil {
		log.Fatal(err)
	}

	var alUrls AllJSONModels

	_ = json.Unmarshal(b, &alUrls.AllUrls)
	//if err != nil {
	//	log.Fatal(err)
	//}

	alUrls.AllUrls = append(alUrls.AllUrls, newURL)
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
