package handler

import (
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/22Fariz22/shorturl/cookies"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/oklog/ulid/v2"

	"github.com/go-chi/chi/v5"
)

type HandlerModel struct {
	Repository repository.Repository
	count      int
	cfg        config.Config
	conn       *pgx.Conn
}

type Handler HandlerModel

type reqURL struct {
	URL string `json:"url"`
}

var rURL reqURL
var (
	db *sql.DB
)

func NewHandler(repo repository.Repository, cfg *config.Config) *Handler {
	count := 0
	return &Handler{
		Repository: repo,
		count:      count,
		cfg:        *cfg,
	}
}
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), h.cfg.DatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = conn.Ping(ctx)
	status := http.StatusOK

	if err != nil {
		status = http.StatusInternalServerError
	}
	//status := h.Repository.Ping()
	w.WriteHeader(status)

}

func GenUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	moreShorter := id.String()[len(id.String())-7:]
	return moreShorter
}

//вернуть все свои URL
func (h *Handler) GetAllURL(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	type resp struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var res []resp

	list := h.Repository.GetAll(r.Cookies()[0].Value)

	fmt.Println(list)
	for i := range list {
		for k, v := range list[i] {
			res = append(res, resp{
				ShortURL:    h.cfg.BaseURL + "/" + k,
				OriginalURL: v,
			})
		}
	}
	res1, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")

	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(res1)
}

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	//сокращатель
	short := GenUlid()

	h.Repository.SaveURL(short, string(payload), r.Cookies()[0].Value)

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
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	fmt.Println("CreateShortURLJSON len(r.Cookies()):", len(r.Cookies()))

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
	h.Repository.SaveURL(short, rURL.URL, r.Cookies()[0].Value)

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
