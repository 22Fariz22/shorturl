package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Handler struct {
	urls  map[string]string
	count int
}

type Url struct {
	url string `json:"url"`
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
	var url Url
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(url)

	//countStr := strconv.Itoa(h.count)
	//
	//defer r.Body.Close()
	////payload, err := io.ReadAll(r.Body)
	//
	//payload, err := ioutil.ReadAll(r.Body)
	//fmt.Println(payload, err)
	//
	//if err != nil {
	//	log.Printf("error: %s", err)
	//} else {
	//	h.urls[countStr] = string(payload)
	//	h.count++
	//	w.WriteHeader(http.StatusCreated)
	//	fmt.Fprintf(w, "%s", "http://localhost:8080/"+countStr)
	//}
}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
//и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	vars := chi.URLParam(r, "id")

	i, ok := h.urls[vars]
	if ok {
		w.Header().Set("Location", i)
		http.Redirect(w, r, i, http.StatusTemporaryRedirect)
	}
}
