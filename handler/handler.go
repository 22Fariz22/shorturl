package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	urls map[string]string
}

func NewHandler() *Handler {
	return &Handler{
		urls: make(map[string]string),
	}
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	/*
		забираем url адресс из тела
		генерим число и добавляем в мапу id[url]
		возвращаем сокращенный url
	*/

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000
	count := rand.Intn(max-min+1) + min
	countStr := strconv.Itoa(count)

	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("error: %s", err)
	} else {

		h.urls[countStr] = string(payload)
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

	vars := mux.Vars(r)

	i, ok := h.urls[vars["id"]]
	if ok {
		r.Header.Set("Location", i)
		http.Redirect(w, r, i, http.StatusTemporaryRedirect)
	}
}
