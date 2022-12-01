package handler

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/oklog/ulid/v2"

	"github.com/go-chi/chi/v5"
)

type HandlerModel struct {
	Repository repository.Repository
	count      int
	cfg        config.Config
}

type Handler HandlerModel

type reqURL struct {
	URL string `json:"url"`
}

var rURL reqURL

func NewHandler(repo repository.Repository, cfg *config.Config) *Handler {
	count := 0
	return &Handler{
		Repository: repo,
		count:      count,
		cfg:        *cfg,
	}
}

func GenUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	//сокращатель
	short := GenUlid()

	h.Repository.SaveURL(short, string(payload))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.cfg.BaseURL + "/" + short))
}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")
	i, ok := h.Repository.GetURL(vars)
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

	if err := json.Unmarshal(payload, &rURL); err != nil {
		log.Print(err)
	}

	short := GenUlid()

	type respURL struct {
		Result string `json:"result"`
	}

	resp := respURL{
		Result: h.cfg.BaseURL + "/" + short,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}

	//пишем в json файл если есть FileStoragePath
	h.Repository.SaveURL(short, rURL.URL)

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
