package repo

import (
	"22Fariz22/shorturl/handler/config"
	"22Fariz22/shorturl/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	//file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0777)

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

	b, err := ioutil.ReadAll(p.file)
	if err != nil {
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

type Handler model.HandlerModel

//функция для востановления списка urls
func (p *producer) RecoverEvents() {

	b, err := ioutil.ReadAll(p.file)
	if err != nil {
		log.Fatal(err)
	}
	if len(b) != 0 {
		var alUrls AllJSONModels

		_ = json.Unmarshal(b, &alUrls.AllUrls)

		for _, v := range alUrls.AllUrls {
			fmt.Println(v)
			for i := 0; i < v.Count; i++ {
				iStr := strconv.Itoa(i)

				//map[string]string{iStr: v.Url[iStr]}  нужный нам map
				hd := Handler{Count: 2, Urls: map[string]string{iStr: v.Url[iStr]}}
				fmt.Println(hd)
				//типа востанавливаем?
				//Handler{Urls:urlMake[iStr]= v.Url[iStr]}

				//hd.Urls[iStr] = v.Url[iStr]
				//hd.Count = v.Count

			}
		}
	}
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

func (c *consumer) Close() error {
	return c.file.Close()
}
