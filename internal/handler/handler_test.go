// тестирование
package handler_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/cookies"
	"github.com/22Fariz22/shorturl/internal/handler"
	"github.com/22Fariz22/shorturl/internal/usecase/memory"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/22Fariz22/shorturl/pkg/logger"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// exampleTest тестирование
func exampleTest() {
	short := handler.GenUlid()
	shortURL := hex.EncodeToString([]byte(short))
	log.Printf("Short URL is %s \n", shortURL)
}

// BenchmarkGenerateShortLink бенчмарк генератора шортурлов
func BenchmarkGenerateShortLink(b *testing.B) {
	for i := 0; i < b.N; i++ {
		short := handler.GenUlid()
		shortURL := hex.EncodeToString([]byte(short))
		b.Logf("Short URL is %s \n", shortURL)
	}
}

func TestHandlerCreateShortURLJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []byte
		status int
	}{
		{
			name:   "create json",
			in:     []byte(`{"url":"https://google.ru"}`),
			status: 201,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			repo := memory.New()

			secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
			if err != nil {
				log.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(test.in))
			w := httptest.NewRecorder()

			h := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					SecretKey:     secretKey,
				},
				Workers: nil,
			}
			h.CreateShortURLJSON(w, req)

		})
	}
}

func TestHandlerCreateShortURLJSONExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []byte
		status int
	}{
		{
			name:   "create json",
			in:     []byte(`{"url":"https://google.ru"}`),
			status: 409,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			repo := memory.New()

			secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
			if err != nil {
				log.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(test.in))
			w := httptest.NewRecorder()

			h := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					SecretKey:     secretKey,
				},
				Workers: nil,
			}
			h.CreateShortURLJSON(w, req)

			req2 := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(test.in))
			w2 := httptest.NewRecorder()

			h2 := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					SecretKey:     secretKey,
				},
				Workers: nil,
			}
			h2.CreateShortURLJSON(w2, req2)
		})
	}
}

func TestHandler_Batch(t *testing.T) {
	tests := []struct {
		name   string
		status int
		in     string
	}{
		{
			name:   "batch",
			status: 201,
			in: `[
		{
			"correlation_id": "<строковый идентификатор>",
			"original_url": "<URL для сокращения>"
		}
		]`},
	}

	for _, tt := range tests {
		repo := memory.New()

		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					SecretKey:     secretKey,
				},
				Workers: nil,
			}

			//body, _ := json.Marshal(tt.in)

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader([]byte(tt.in)))
			w := httptest.NewRecorder()

			cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)

			h.Batch(w, req)

			fmt.Println("status: ", w.Code)

		})
	}
}

func TestHandler_DeleteHandler(t *testing.T) {
	l := logger.New("debug")
	ctx := context.Background()

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "test delete OK",
			status: 202,
		},
		{
			name:   "delete status 400",
			status: 400,
		},
	}

	repo := memory.New()

	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
	worker := worker.NewWorkerPool(l, repo)

	h := &handler.Handler{
		Repository: repo,
		Cfg: config.Config{
			ServerAddress: "localhost:8080",
			BaseURL:       "http://localhost:8080",
			SecretKey:     secretKey,
		},
		Workers: worker,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(`["short"]`)

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bodyReader)
			w := httptest.NewRecorder()

			cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)

			repo.SaveURL(ctx, l, "short", "https://ya.ru", "7654321") //save before delete

			h.DeleteHandler(w, req)
		})
		//t.Run(tt.name, func(t *testing.T) {
		//	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil)
		//	w := httptest.NewRecorder()
		//
		//	cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)
		//
		//	h.DeleteHandler(w, req)
		//
		//	fmt.Println("status: ", w.Code)
		//})
	}
}

func TestHandler_GetAllURL(t *testing.T) { // получилось
	l := logger.New("debug")

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "get all urls",
			status: 200,
		},
	}

	for _, tt := range tests {
		repo := memory.New()

		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					SecretKey:     secretKey,
				},
				Workers: nil,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)

			repo.SaveURL(context.Background(), l, "short", "https://ya.ru", "7654321")

			h.GetAllURL(w, req)

			fmt.Println("status: ", w.Code)

		})
	}
}

func TestHandler_GetShortURLByIDHandler(t *testing.T) {
	l := logger.New("debug")

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "get url",
			status: 307,
		},
	}

	repo := memory.New()

	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	h := &handler.Handler{
		Repository: repo,
		Cfg: config.Config{
			ServerAddress: "localhost:8080",
			BaseURL:       "http://localhost:8080",
			SecretKey:     secretKey,
		},
		Workers: nil,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/short", nil)
			w := httptest.NewRecorder()

			cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)

			repo.SaveURL(context.Background(), l, "short", "https://ya.ru", "7654321")

			h.GetShortURLByIDHandler(w, req)

			fmt.Println(w.Code)
			fmt.Println("Location: ", w.Header().Get("Location"))
		})
	}
}

func TestHandler_CreateShortURLHandler(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "CreateShortURL OK",
			status: 201,
		},
	}
	for _, tt := range tests {
		repo := memory.New()

		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
					SecretKey:     secretKey,
				},
				Workers: nil,
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("https://ya.ru")))
			w := httptest.NewRecorder()

			cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)

			h.CreateShortURLHandler(w, req)

			fmt.Println(w.Code)

			assert.EqualValues(t, tt.status, w.Code)

		})
	}
}

func TestHandler_Ping(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "ping OK",
			status: 200,
		},
	}
	for _, tt := range tests {
		repo := memory.New()

		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: repo,
				Cfg: config.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
				},
				Workers: nil,
			}

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			h.Ping(w, req)
			assert.EqualValues(t, tt.status, w.Code)
		})
	}
}
