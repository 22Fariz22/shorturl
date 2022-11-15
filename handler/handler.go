package handler

import (
	"encoding/json"
	"fmt"
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
	Count int `json:"count"`
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

func NewHandler() *Handler {
	c := 0
	return &Handler{
		urls:  make(map[string]string),
		Count: c,
	}
}

//RequestURL принимает URL
func (c *CreateShortURLRequest) RequestURL(w http.ResponseWriter, r *http.Request) string {
	return fmt.Sprintf("здесь будет body с URL")
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	var value CreateShortURLRequest

	if err != nil {
		log.Printf("error: %s", err)
	} else {
		h.mu.Lock()
		countStr := strconv.Itoa(h.Count)
		h.mu.Unlock()
		value.URL = string(payload)
		h.mu.Lock()
		h.urls[countStr] = value.URL
		h.Count++
		h.mu.Unlock()

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
		var value CreateShortURLRequest

		if err := json.Unmarshal(payload, &value); err != nil {
			log.Printf("error: %s", err)
		}

		type Resp struct {
			Result string `json:"result"`
		}
		resp := Resp{
			Result: "http://localhost:8080/" + countStr,
		}
		res, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		h.mu.Lock()
		h.urls[countStr] = value.URL
		h.Count++
		h.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}
