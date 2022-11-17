package handler

import (
	"22Fariz22/shorturl/handler/config"
	"22Fariz22/shorturl/repo"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	mu    sync.Mutex
	urls  map[string]string `json:"urls"` //map[0:http://ya.ru]
	count int               `json:"count"`
}

//type CreateShortURLRequest struct {
//	URL string `json:"url"`
//}

var value repo.CreateShortURLRequest
var valueArr []repo.CreateShortURLRequest

func NewHandler() *Handler {
	c := 0
	return &Handler{
		urls:  make(map[string]string),
		count: c,
	}
}

func (h *Handler) ShortenURL(bodyStr string) string {

	//var value CreateShortURLRequest

	h.mu.Lock()
	countStr := strconv.Itoa(h.count)
	h.mu.Unlock()

	value.URL = bodyStr

	h.mu.Lock()
	h.urls[countStr] = value.URL
	h.count++
	h.mu.Unlock()

	return countStr
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)

	cfg := config.NewConnectorConfig()

	fileName := cfg.FileStoragePath
	producer, err := repo.NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	//проверяем счетчик, если 0,то это превый запуск
	//а если первывый запуск то проверяем есть ли запись в events.log
	//если есть запись, то востанавливаем все записи в память
	if h.count == 0 {
		//RecoverEvents()

	}

	if err != nil {
		log.Printf("error: %s", err)
	} else {
		short := h.ShortenURL(string(payload))

		producer.WriteEvent(h.count, h.urls)
		//пишем в файл
		//producer, err := repo.NewProducer(cfg.FileStoragePath)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//defer producer.Close()
		//if err := producer.WriteEvent(&value); err != nil {
		//	log.Fatal(err)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(cfg.BaseURL + "/" + short))
	}

}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	i, ok := h.urls[vars]
	if ok {
		w.Header().Set("Location", i)
		http.Redirect(w, r, i, http.StatusTemporaryRedirect)
	}
}

func (h *Handler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	cfg := config.NewConnectorConfig()

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %s", err)
	} else {

		if err := json.Unmarshal(payload, &value); err != nil {
			log.Printf("error: %s", err)
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

		//пишем в файл
		producer, err := repo.NewProducer(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
		if err := producer.WriteEvent(h.count, h.urls); err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}
