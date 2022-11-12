package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	mu    sync.Mutex        `json:"-"`
	urls  map[string]string `json:"-"`
	Count int               `json:"count"`
	Url   string            `json:"url"`
}

func NewHandler() *Handler {
	c := 0
	return &Handler{
		urls:  make(map[string]string),
		Count: c,
	}
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {

	countStr := strconv.Itoa(h.Count)

	payload, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("error: %s", err)
	} else {

		h.mu.Lock()
		defer h.mu.Unlock()
		h.urls[countStr] = string(payload)
		h.Count++

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + countStr))
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
	countStr := strconv.Itoa(h.Count)

	payload, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("error: %s", err)
	} else {
		h.mu.Lock()
		defer h.mu.Unlock()

		value := Handler{}

		if err := json.Unmarshal(payload, &value); err != nil {
			log.Printf("error: %s", err)
		}

		type Resp struct {
			URL      string `json:"url"`
			ShortURL string `json:"short_url"`
		}
		resp := Resp{
			URL:      value.Url,
			ShortURL: countStr,
		}
		res, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		h.urls[countStr] = value.Url
		h.Count++

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}
