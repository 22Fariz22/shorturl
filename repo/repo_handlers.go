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
	Cfg  *config.Config
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}
	return &Producer{
		file: file,
		Cfg:  config.NewConnectorConfig(),
	}, nil
}

func (p *Producer) WriteEvent(cnt int, urlMap map[string]string) error {

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
	ioutil.WriteFile(p.Cfg.FileStoragePath, newURLBytes, 0666)
	return nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}
