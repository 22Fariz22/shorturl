// Package handler принимает и отправляет
package handler

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/22Fariz22/shorturl/pkg/logger"

	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/cookies"
	"github.com/22Fariz22/shorturl/internal/usecase"

	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

const CtxTimeOut = 5 * time.Second

//go:generate make total

// Handler структура хэндлер
type Handler struct {
	Repository usecase.Repository
	Cfg        config.Config
	Workers    *worker.Pool
	l          logger.Interface
}

type reqURL struct {
	URL string `json:"url"`
}

var rURL reqURL

// NewHandler создает хэндлер
func NewHandler(repo usecase.Repository, cfg *config.Config, workers *worker.Pool, l logger.Interface) *Handler {
	return &Handler{
		Repository: repo,
		Cfg:        *cfg,
		Workers:    workers,
		l:          l,
	}
}

// Stats возвращает количество сокращённых URL в сервисе и количество пользователей в сервисе
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	type stats struct {
		Urls  int `json:"urls"`
		Users int `json:"users"`
	}

	var st stats

	addr := r.RemoteAddr

	ipStr, _, err := net.SplitHostPort(addr)
	if err != nil {
		h.l.Info("err net.SplitHostPort: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		h.l.Info("nil from net.ParseIP: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, ipnet, err := net.ParseCIDR(h.Cfg.TrustedSubnet)
	if err != nil {
		h.l.Info("err net.ParseCIDR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ipnet.Contains(ip) {
		urls, users, err := h.Repository.Stats(context.Background(), h.l)
		if err != nil {
			h.l.Info("h.Repository.Stats: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		st.Urls = urls
		st.Users = users
	} else {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	outJSON, err := json.Marshal(st)
	if err != nil {
		h.l.Info("err net.ParseIP: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(outJSON)

}

// DeleteHandler удаляет запись
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.Cfg.SecretKey)
	}

	var list []string

	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		h.l.Info("", err)
		//status 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	h.Workers.AddJob(ctx, h.l, list, r.Cookies()[0].Value)

	//status 202
	w.WriteHeader(http.StatusAccepted)
}

// GetAllURL получить все записи
func (h *Handler) GetAllURL(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.Cfg.SecretKey)
	}

	type resp struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var res []resp

	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	list, err := h.Repository.GetAll(ctx, h.l, r.Cookies()[0].Value)
	//status 204
	if err != nil {
		h.l.Info("", err)
		w.WriteHeader(http.StatusNoContent)
	}

	for i := range list {
		for k, v := range list[i] {
			res = append(res, resp{
				ShortURL:    h.Cfg.BaseURL + "/" + k,
				OriginalURL: v,
			})
		}
	}
	res1, err := json.Marshal(res)
	if err != nil {
		h.l.Info("", err)
	}

	w.Header().Set("Content-Type", "application/json")

	//status 204
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		//status 200
		w.WriteHeader(http.StatusOK)
	}

	w.Write(res1)
}

// CreateShortUrlHandler эндпоинт POST / принимает в теле запроса строку URL для сокращения
func (h *Handler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.Cfg.SecretKey)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.l.Info("", err)
	}

	//сокращатель
	short := GenUlid()

	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, h.l, short, string(payload), r.Cookies()[0].Value)

	//status 409
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(h.Cfg.BaseURL + "/" + s))
		return
	}

	//status 201
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.Cfg.BaseURL + "/" + short))
}

// GetShortURLByIDHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.Cfg.SecretKey)
	}

	vars := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	url, ok := h.Repository.GetURL(ctx, h.l, vars)

	//status 400
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//status 401
	if url.Deleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	w.Header().Set("Location", url.LongURL)

	//status 307
	http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
}

// Batch создает записи конвеером
func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.Cfg.SecretKey)
	}

	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	var batchResp []entity.PackReq

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.l.Info("", err)
		http.Error(w, "", 500)
		return
	}
	if err := json.Unmarshal(payload, &batchResp); err != nil {
		h.l.Info("", err)
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
			ShortURL:      h.Cfg.BaseURL + "/" + short,
		}

		listReq = append(listReq, req)

		listResp = append(listResp, resp)
	}

	err = h.Repository.RepoBatch(ctx, h.l, r.Cookies()[0].Value, listReq)
	if err != nil {
		h.l.Info("", err)
		return
	}

	res, err := json.Marshal(listResp)
	if err != nil {
		h.l.Info("", err)
	}
	w.Header().Set("Content-Type", "application/json")

	//status 201
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// CreateShortURLJSON создает запись из тела json
func (h *Handler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.Cfg.SecretKey)
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.l.Info("", err)
		http.Error(w, "", 500)
	}

	if err := json.Unmarshal(payload, &rURL); err != nil {
		h.l.Info("", err)
		return
	}

	short := GenUlid()

	type respURL struct {
		Result string `json:"result"`
	}

	resp := respURL{
		Result: h.Cfg.BaseURL + "/" + short,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		h.l.Info("", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, h.l, short, rURL.URL, r.Cookies()[0].Value)

	if err != nil {
		resp1 := respURL{
			Result: h.Cfg.BaseURL + "/" + s,
		}
		res1, err := json.Marshal(resp1)
		if err != nil {
			h.l.Info("", err)
		}
		w.Header().Set("Content-Type", "application/json")
		//status 409
		w.WriteHeader(http.StatusConflict)
		w.Write(res1)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// Close закрывает соединение
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

// Ping проверка соединения
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), CtxTimeOut)
	defer cancel()

	err := h.Repository.Ping(ctx, h.l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GenUlid генератор шортурлов
func GenUlid() string {
	//time.Sleep(1 * time.Second)
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	moreShorter := id.String()[len(id.String())-7:]
	return moreShorter
}
