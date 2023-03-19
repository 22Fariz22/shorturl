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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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
		name string
		in   []byte
		want string
	}{
		{
			name: "create json",
			in:   []byte(`{"url":"https://google.ru"}`),
			want: `{"result":"http://localhost:8080/`,
		},
	}

	cfg := config.NewConfig()

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {

			repo := memory.New()

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(test.in))
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.NewHandler(repo, cfg, nil).CreateShortURLJSON(w, req)

			response, err := io.ReadAll(w.Result().Body)
			uulid := string(response)[33:]

			require.NoError(t, err)
			assert.JSONEq(t, test.want+uulid+``, string(response)) // как вставить  want геренрирующий
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

func TestHandler_DeleteHandler(t *testing.T) { //status 400. как исправить bodyReader?
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "test delete",
			status: 202,
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

			//bodyReader := strings.NewReader(`["short"]`)

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			repo.SaveURL(context.Background(), "short", "https://ya.ru", "7654321") //save before delete

			h.DeleteHandler(w, req)

			fmt.Println("status: ", w.Code)

		})
	}
}

func TestHandler_GetAllURL(t *testing.T) { // получилось
	//bodyReader := strings.NewReader("")

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

			repo.SaveURL(context.Background(), "short", "https://ya.ru", "7654321")

			h.GetAllURL(w, req)

			fmt.Println("status: ", w.Code)

		})
	}
}

func TestHandler_GetShortURLByIDHandler(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "get url",
			status: 307,
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

			req := httptest.NewRequest(http.MethodGet, "/short", nil)
			w := httptest.NewRecorder()

			cookies.SetCookieHandler(w, req, h.Cfg.SecretKey)

			repo.SaveURL(context.Background(), "short", "https://ya.ru", "7654321")

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
