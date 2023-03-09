// Package handler принимает и отправляет
package handler

import (
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/cookies"
	"github.com/22Fariz22/shorturl/internal/usecase"

	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

const ctxTimeOut = 5 * time.Second

// Handler структура хэндлер
type Handler struct {
	Repository usecase.Repository
	cfg        config.Config
	workers    *worker.Pool
}

type reqURL struct {
	URL string `json:"url"`
}

var rURL reqURL

// NewHandler создает хэндлер
func NewHandler(repo usecase.Repository, cfg *config.Config, workers *worker.Pool) *Handler {
	return &Handler{
		Repository: repo,
		cfg:        *cfg,
		workers:    workers,
	}
}

// DeleteHandler удаляет запись
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	var list []string

	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	h.workers.AddJob(ctx, list, r.Cookies()[0].Value)

	w.WriteHeader(http.StatusAccepted)
}

//GetAllURL получить все записи
func (h *Handler) GetAllURL(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	type resp struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var res []resp

	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	list, err := h.Repository.GetAll(ctx, r.Cookies()[0].Value)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
	}

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

//CreateShortUrlHandler эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	//сокращатель
	short := GenUlid()

	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, short, string(payload), r.Cookies()[0].Value)
	fmt.Println("su in handler(0)", s, err)

	if err != nil {
		fmt.Println("print in handler(1)", s, err)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(h.cfg.BaseURL + "/" + s))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.cfg.BaseURL + "/" + short))
}

//GetShortURLByIDHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	vars := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	url, ok := h.Repository.GetURL(ctx, vars)
	log.Print("in handler Get url,ok:", url, ok)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url.Deleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	w.Header().Set("Location", url.LongURL)
	http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
}

//Batch создает записи конвеером
func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	var batchResp []entity.PackReq

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "", 500)
		return
	}
	if err := json.Unmarshal(payload, &batchResp); err != nil { // [{1 mail.ru} {2 ya.ru} {3 google.ru}]
		log.Print(err)
		return
	}

	var listReq []entity.PackReq

	var listResp []entity.PackResponse

	for _, resp := range batchResp {
		short := GenUlid()

		req := entity.PackReq{
			CorrelationID: resp.CorrelationID,
			OriginalURL:   resp.OriginalURL,
			ShortURL:      short,
		}

		resp := entity.PackResponse{
			CorrelationID: resp.CorrelationID,
			ShortURL:      h.cfg.BaseURL + "/" + short,
		}

		listReq = append(listReq, req)

		listResp = append(listResp, resp)
	}

	err = h.Repository.RepoBatch(ctx, r.Cookies()[0].Value, listReq)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := json.Marshal(listResp)
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

//CreateShortURLJSON создает запись из тела json
func (h *Handler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "", 500)
	}

	if err := json.Unmarshal(payload, &rURL); err != nil {
		log.Print(err)
		return
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

	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, short, rURL.URL, r.Cookies()[0].Value)
	fmt.Println("out err", s, err)

	if err != nil {
		fmt.Println("in err", s, err)

		resp1 := respURL{
			Result: h.cfg.BaseURL + "/" + s,
		}
		res1, err := json.Marshal(resp1)
		if err != nil {
			log.Print(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(res1)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

//Close закрывает соединение
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

// структура для сжатия
type gzipReader struct {
	*gzip.Reader
	io.Closer
}

//Ping проверка соединения
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeOut)
	defer cancel()

	err := h.Repository.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//GenUlid генератор шортурлов
func GenUlid() string {
	//time.Sleep(1 * time.Second)
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	moreShorter := id.String()[len(id.String())-7:]
	return moreShorter
}

func Example_GenUlid() {
	short := GenUlid()
	shortURL := hex.EncodeToString([]byte(short))
	log.Printf("Short URL is %s \n", shortURL)
}
