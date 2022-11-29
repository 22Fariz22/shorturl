package repo

import (
	"encoding/json"
	"github.com/22Fariz22/shorturl/model_json"
	"io/ioutil"
	"log"
	"os"

	"github.com/22Fariz22/shorturl/handler/config"
)

//type JSONModel struct {
//	Count int               `json:"count"`
//	URL   map[string]string `json:"url"`
//}
//
//type AllJSONModels struct {
//	AllUrls []*JSONModel
//}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type Producer struct {
	File *os.File
	Cfg  *config.Config
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Producer{
		File: file,
		Cfg:  config.NewConnectorConfig(),
	}, nil
}

func (p *Producer) WriteEvent(id int, urlMap map[string]string) error {

	newURL := &model_json.JSONModel{}
	newURL.URL = urlMap
	newURL.ID = id

	b, err := ioutil.ReadAll(p.file)
	if err != nil {
		log.Print(err)
	}

	var alUrls model_json.AllJSONModels

	_ = json.Unmarshal(b, &alUrls.AllUrls)
	//if err != nil {
	//	log.Fatal(err)
	//}

	alUrls.AllUrls = append(alUrls.AllUrls, newURL)
	newURLBytes, err := json.MarshalIndent(&alUrls.AllUrls, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(p.Cfg.FileStoragePath, newURLBytes, 0666)
	if err != nil {
		log.Print(err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}
