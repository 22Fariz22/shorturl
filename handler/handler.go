package handler

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/22Fariz22/shorturl/model"

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

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := h.Repository.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GenUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	moreShorter := id.String()[len(id.String())-7:]
	return moreShorter
}

//
func (h *Handler) GetAllURL(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	type resp struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var res []resp

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, short, string(payload), r.Cookies()[0].Value)
	//fmt.Println("value SaveUrl in handler:", s)
	//fmt.Println("err SaveUrl in handler:", err)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(h.cfg.BaseURL + "/" + short))
		return
	}
	w.WriteHeader(http.StatusConflict)
	w.Write([]byte(h.cfg.BaseURL + "/" + s))

	//w.WriteHeader(http.StatusCreated)
	//w.Write([]byte(h.cfg.BaseURL + "/" + short))
}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	vars := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	i, ok := h.Repository.GetURL(ctx, vars, r.Cookies()[0].Value)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", i)
	http.Redirect(w, r, i, http.StatusTemporaryRedirect)
}

func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	var batchResp []model.PackReq

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

	var listReq []model.PackReq

	var listResp []model.PackResponse

	for i := range batchResp {
		short := GenUlid()

		req := model.PackReq{
			CorrelationID: batchResp[i].CorrelationID,
			OriginalURL:   batchResp[i].OriginalURL,
			ShortURL:      short,
		}

		resp := model.PackResponse{
			CorrelationID: batchResp[i].CorrelationID,
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, short, rURL.URL, r.Cookies()[0].Value)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
		return
	}

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

	////пишем в json файл если есть FileStoragePath
	//h.Repository.SaveURL(ctx, short, rURL.URL, r.Cookies()[0].Value)
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	//w.Write(res)
	//

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
