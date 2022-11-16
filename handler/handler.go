package handler

import (
	"22Fariz22/shorturl/handler/config"
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
	urls  map[string]string
	count int
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

func NewHandler() *Handler {
	c := 0
	return &Handler{
		urls:  make(map[string]string),
		count: c,
	}
}

func (h *Handler) ShortenURL(bodyStr string) string {

	var value CreateShortURLRequest

	h.mu.Lock()
	countStr := strconv.Itoa(h.count)
	h.mu.Unlock()

	value.URL = bodyStr

	h.mu.Lock()
	h.urls[countStr] = value.URL
	h.count++
	h.mu.Unlock()

	//return "http://localhost:8080/" + countStr
	return countStr
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)

	cfg := config.NewConnectorConfig()

	if err != nil {
		log.Printf("error: %s", err)
	} else {
		short := h.ShortenURL(string(payload))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(cfg.BaseURL + short))
	}
}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
//и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
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

		var value CreateShortURLRequest

		if err := json.Unmarshal(payload, &value); err != nil {
			log.Printf("error: %s", err)
		}

		type Resp struct {
			Result string `json:"result"`
		}
		resp := Resp{
			Result: cfg.BaseURL + h.ShortenURL(value.URL),
		}
		res, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}
