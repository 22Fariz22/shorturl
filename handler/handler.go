package handler

import (
	"compress/gzip"
	"encoding/json"
	"github.com/22Fariz22/shorturl/handler/config"
	"github.com/22Fariz22/shorturl/model_json"
	"github.com/22Fariz22/shorturl/repository"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/22Fariz22/shorturl/repo"

	"github.com/go-chi/chi/v5"
)

type HandlerModel struct {
	Repository repository.Repository
}

type Handler HandlerModel

var value repo.CreateShortURLRequest

func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		Repository: repo,
	}
}

//функция для востановления списка urls
func (h *Handler) RecoverEvents(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal("", err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	if len(b) != 0 {
		var alUrls model_json.AllJSONModels

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

	//h.Mu.Lock()
	countStr := strconv.Itoa(h.Count)
	//h.Mu.Unlock()

	value.URL = bodyStr

	//h.Mu.Lock()
	h.Urls[countStr] = value.URL
	h.Count++
	//h.Mu.Unlock()

	return countStr
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	//сокращатель
	short := h.ShortenURL(string(payload))

	h.Repository.SaveURL(short, string(payload))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.DefaultBaseURL + "/" + short))
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

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(payload, &value); err != nil {
		log.Print(err)
	}

	type Resp struct {
		Result string `json:"result"`
	}

	resp := Resp{
		Result: config.DefaultBaseURL + "/" + h.ShortenURL(value.URL),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}

	//пишем в json файл если есть FileStoragePath
	h.Repository.SaveURL(h.Count, h.Urls)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (r gzipReader) Close() error {
	if err := r.Closer.Close(); err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

// DeCompress возвращает распакованный gz
func DeCompress(next http.Handler) http.Handler {
	// переменная reader будет равна r.Body или *gzip.Reader
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}
		request.Header.Del("Content-Length")
		reader, err := gzip.NewReader(request.Body)
		if err != nil {
			io.WriteString(writer, err.Error())
			return
		}

		defer reader.Close()

		request.Body = gzipReader{
			reader,
			request.Body,
		}

		next.ServeHTTP(writer, request)
	})
}

type gzipReader struct {
	*gzip.Reader
	io.Closer
}
