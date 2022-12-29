package handler

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/cookies"
	"github.com/22Fariz22/shorturl/model"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/worker"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

//type HandlerModel struct {
//	Repository repository.Repository
//	count      int // кажется не нужная штука. проверить и удалить
//	cfg        config.Config
//}

type Handler struct {
	Repository repository.Repository
	count      int // кажется не нужная штука. проверить и удалить
	cfg        config.Config
	workers    *worker.WorkerPool
}

type reqURL struct {
	URL string `json:"url"`
}

var rURL reqURL

func NewHandler(repo repository.Repository, cfg *config.Config, workers *worker.WorkerPool) *Handler {
	return &Handler{
		Repository: repo,
		cfg:        *cfg,
		workers:    workers,
	}
}

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

	h.workers.AddJob(list, r.Cookies()[0].Value)
	//h.Repository.Delete(ctx, list, r.Cookies()[0].Value)

	w.WriteHeader(http.StatusAccepted)
}

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
	//shr := uuid.New().NodeID()
	//short := hex.EncodeToString(shr)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	s, err := h.Repository.SaveURL(ctx, short, string(payload), r.Cookies()[0].Value)
	fmt.Println("out err", s, err)

	if err != nil {
		fmt.Println("in err", s, err)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(h.cfg.BaseURL + "/" + s))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.cfg.BaseURL + "/" + short))
}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
func (h *Handler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	vars := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	url, ok := h.Repository.GetURL(ctx, vars)

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

func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		cookies.SetCookieHandler(w, r, h.cfg.SecretKey)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

		//shr := uuid.New().NodeID()
		//short := hex.EncodeToString(shr)

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
	//shr := uuid.New().NodeID()
	//short := hex.EncodeToString(shr)

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

/*
func (dbs *DBStorage) LoadLong(ctx context.Context, shortlink string) (string, string, error) {
	var longlink string
	var exist string
	rows := dbs.storage.QueryRow(ctx, "SELECT longlink, exist FROM linklisttable WHERE shortlink = $1",
		shortlink)
	err := rows.Scan(&longlink, &exist) if err != nil { fmt.Println("Error in LoadLong (DBStorage):", err)
		return "", "", err
		}
		return longlink, exist, nil
}

func (dbs *DBStorage) Store(ctx context.Context, userID string, longlink string, shortlink string) error {
	_, err := dbs.storage.Exec(ctx, "INSERT INTO linklisttable (userhexid, longlink, shortlink, exist)"+
		" VALUES ($1, $2, $3, $4) ",
		userID, longlink, shortlink, "Yes")
	if err != nil {
		fmt.Println("Error in Store (DBStorage):", err)
		return err
	}
	return nil
}

*/
