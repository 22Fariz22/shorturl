package handler

import (
	"22Fariz22/shorturl/handler/config"
	"22Fariz22/shorturl/model"
	"22Fariz22/shorturl/repo"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler model.HandlerModel

var value repo.CreateShortURLRequest

func NewHandler() *Handler {
	c := 0
	return &Handler{
		Urls:  make(map[string]string),
		Count: c,
	}
}

type JSONModel struct {
	Count int               `json:"count"`
	URL   map[string]string `json:"url"`
}

type AllJSONModels struct {
	AllUrls []*JSONModel
}

//функция для востановления списка urls
func (h *Handler) RecoverEvents() {
	cfg := config.NewConnectorConfig()
	fileName := cfg.FileStoragePath
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal("", err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	if len(b) != 0 {
		var alUrls AllJSONModels

		err = json.Unmarshal(b, &alUrls.AllUrls)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range alUrls.AllUrls {
			for i := 0; i < v.Count+1; i++ {
				iStr := strconv.Itoa(i)
				//типа востанавливаем
				h.Urls[iStr] = v.URL[iStr]
				h.Count = v.Count

			}
		}
	}
}

func (h *Handler) ShortenURL(bodyStr string) string {

	h.Mu.Lock()
	countStr := strconv.Itoa(h.Count)
	h.Mu.Unlock()

	value.URL = bodyStr

	h.Mu.Lock()
	h.Urls[countStr] = value.URL
	h.Count++
	h.Mu.Unlock()

	return countStr
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConnectorConfig()

	//сокращатель
	short := h.ShortenURL(string(payload))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(cfg.BaseURL + "/" + short))

}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	i, ok := h.Urls[vars]
	if ok {
		w.Header().Set("Location", i)
		http.Redirect(w, r, i, http.StatusTemporaryRedirect)
	}
}

func (h *Handler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	cfg := config.NewConnectorConfig()

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(payload, &value); err != nil {
		log.Fatal(err)
	}

	type Resp struct {
		Result string `json:"result"`
	}
	resp := Resp{
		Result: cfg.BaseURL + "/" + h.ShortenURL(value.URL),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
