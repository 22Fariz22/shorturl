package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	urls  map[string]string
	count int
}

func NewHandler() *Handler {
	c := 0
	return &Handler{
		urls:  make(map[string]string),
		count: c,
	}
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	/*
		забираем url адресс из тела
		счетчик count  добавляем в мапу count[long_url]
		возвращаем сокращенный url  )
	*/

	countStr := strconv.Itoa(h.count)

	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("error: %s", err)
	} else {
		h.urls[countStr] = string(payload)
		h.count++
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", "http://localhost:8080/"+countStr)
	}
}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
//и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	/*
		принимаем id из url-параметра
		ищем по ключу id значение url в мапе Handler
	*/

	//vars := mux.Vars(r)
	vars := chi.URLParam(r, "id")

	i, ok := h.urls[vars]
	if ok {
		w.Header().Set("Location", i)
		http.Redirect(w, r, i, http.StatusTemporaryRedirect)
	}
}
