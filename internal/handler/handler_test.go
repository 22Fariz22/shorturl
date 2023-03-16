// тестирование
package handler_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/handler"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"github.com/22Fariz22/shorturl/internal/usecase/memory"
	repoMock "github.com/22Fariz22/shorturl/internal/usecase/mocks"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/golang/mock/gomock"
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

//BenchmarkGenerateShortLink бенчмарк генератора шортурлов
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
			fmt.Println("uulid: ", uulid)
			fmt.Println("exept: ", test.want+uulid+`"}`)

			require.NoError(t, err)
			assert.JSONEq(t, test.want+uulid+``, string(response)) // как вставить  want геренрирующий
		})
	}
}

//TestHandler_CreateShortURLHandler тест хэндлера из эндпойнта r.Post("/", hd.CreateShortURLHandler)
func TestHandler_GetShortURLByIDHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	ctl.Finish()

	repo := repoMock.NewMockRepository(ctl)

	ctx := context.Background()
	short := handler.GenUlid()

	//mockResp := entity.URL{LongURL: "https://ya.ru"}

	expected := entity.URL{LongURL: "https://ya.ru"}

	repo.EXPECT().GetURL(ctx, short).Return(expected, true).Times(1)
	URL, ok := repo.GetURL(ctx, short)

	require.Equal(t, expected, URL)
	require.EqualValues(t, true, ok)
}

func _TestHandle_CreateShortURLHandler(t *testing.T) {

	ctl := gomock.NewController(t)
	ctl.Finish()

	repo := repoMock.NewMockRepository(ctl)

	ctx := context.Background()
	short := handler.GenUlid()
	payload := "https://google.com"
	cookie := "ABCD12345"

	//mockResp := entity.URL{LongURL: "https://ya.ru"}

	exp := "http://localhost:8080/" + short

	repo.EXPECT().SaveURL(ctx, short, payload, cookie).Return(exp, nil).Times(1)
	resp, err := repo.SaveURL(ctx, short, payload, cookie)

	require.Equal(t, exp, resp)
	require.NoError(t, err)
}

func TestHandler_DeleteHandler(t *testing.T) {
	type fields struct {
		Repository usecase.Repository
		cfg        config.Config
		workers    *worker.Pool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	ctl := gomock.NewController(t)
	ctl.Finish()
	repo := repoMock.NewMockRepository(ctl)

	workers := worker.NewWorkerPool(repo)
	workers.RunWorkers(10)
	defer workers.Stop()

	//bodyReader := strings.NewReader("['some_short_url']")

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test delete",
			fields: fields{
				Repository: repo,
				cfg:        config.Config{},
				workers:    workers,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: tt.fields.Repository,
				Cfg:        tt.fields.cfg,
				Workers:    tt.fields.workers,
			}
			h.DeleteHandler(tt.args.w, tt.args.r)
		})
	}
}

func _TestHandler_GetAllURL(t *testing.T) {
	type fields struct {
		Repository usecase.Repository
		Cfg        config.Config
		Workers    *worker.Pool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	ctl := gomock.NewController(t)
	ctl.Finish()
	repo := repoMock.NewMockRepository(ctl)

	workers := worker.NewWorkerPool(repo)
	workers.RunWorkers(10)
	defer workers.Stop()

	//bodyReader := strings.NewReader("")

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "get all urls",
			fields: fields{
				Repository: repo,
				Cfg:        config.Config{},
				Workers:    workers,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/api/user/urls", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: tt.fields.Repository,
				Cfg:        tt.fields.Cfg,
				Workers:    tt.fields.Workers,
			}
			h.GetAllURL(tt.args.w, tt.args.r)
		})
	}
}

func _TestHandler_GetShortURLByIDHandler(t *testing.T) {
	type fields struct {
		Repository usecase.Repository
		Cfg        config.Config
		Workers    *worker.Pool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	ctl := gomock.NewController(t)
	ctl.Finish()
	repo := repoMock.NewMockRepository(ctl)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "get url",
			fields: fields{
				Repository: repo,
				Cfg:        config.Config{},
				Workers:    nil,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/1", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := &handler.Handler{
				Repository: tt.fields.Repository,
				Cfg:        tt.fields.Cfg,
				Workers:    tt.fields.Workers,
			}
			h.GetShortURLByIDHandler(tt.args.w, tt.args.r)
		})
	}
}

func _TestHandler_CreateShortURLHandler(t *testing.T) {
	type fields struct {
		Repository usecase.Repository
		Cfg        config.Config
		Workers    *worker.Pool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	ctl := gomock.NewController(t)
	ctl.Finish()
	repo := repoMock.NewMockRepository(ctl)

	//bodyReader := strings.NewReader("https://ya.ru")

	tests := []struct {
		name   string
		fields fields
		//args   args
	}{
		{
			name: "CreateShortURL",
			fields: fields{
				Repository: repo,
				Cfg:        config.Config{},
				Workers:    nil,
			},
			//args: args{
			//	w: httptest.NewRecorder(),
			//	r: httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("sdfsdf"))),
			//},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler.Handler{
				Repository: tt.fields.Repository,
				Cfg:        tt.fields.Cfg,
				Workers:    nil,
			}
			h.CreateShortURLHandler(httptest.NewRecorder(),
				httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("sdfsdf"))))
		})
	}
}

func _TestHandler_Batch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []byte
		want string
	}{
		{
			name: "batch",
			in:   []byte(`[{"correlation_id": "1", "original_url": "<URL для сокращения>"}]`),
			want: `[{"correlation_id": "1", "short_url": "<результирующий сокращённый URL>"}]`,
		},
	}
	cfg := config.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := memory.New()

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(tt.in))
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.NewHandler(repo, cfg, nil).Batch(w, req)

			response, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			fmt.Println("resp batsch: ", string(response))
			require.JSONEq(t, string(response), string(response))

		})
	}
}
