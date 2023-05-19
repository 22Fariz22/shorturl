package db

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/cookies"
	"github.com/22Fariz22/shorturl/internal/handler"
	mock_usecase "github.com/22Fariz22/shorturl/internal/usecase/mocks"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/22Fariz22/shorturl/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_inDBRepository_Stats(t *testing.T) {
	l := logger.New("debug")
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbMock := mock_usecase.NewMockRepository(ctrl)

	cfg := config.NewConfig()
	cfg.TrustedSubnet = "127.0.0.1/8"

	hd := handler.NewHandler(dbMock, cfg, nil, l)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", cfg.BaseURL+"/api/internal/stats", nil)

	dbMock.EXPECT().Stats(ctx, l).Return(0, 0, nil)
	hd.Stats(w, r)

	dbMock.Stats(ctx, l)

	require.Equal(t, http.StatusForbidden, w.Result().StatusCode)

	defer r.Response.Body.Close()
	defer r.Body.Close()
	defer w.Result().Body.Close()
}

//func Test_inDBRepository_SaveURL(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	dbMock := mock_usecase.NewMockRepository(ctrl)
//
//	cfg := config.NewConfig()
//
//	hd := handler.NewHandler(dbMock, cfg, nil)
//
//	//req := cfg.BaseURL + "short"
//
//	w := httptest.NewRecorder()
//	r := httptest.NewRequest("POST", cfg.BaseURL, nil)
//
//	cookies.SetCookieHandler(w, r, cfg.SecretKey)
//	cookies := r.Cookies()[0].Value
//
//	dbMock.EXPECT().SaveURL(gomock.Any(), "myshort", "mylong", cookies).Return("", nil)
//
//	hd.CreateShortURLHandler(w, r)
//	hd.Repository.SaveURL(context.Background(), "myshort", "mylong", cookies)
//
//	fmt.Println("status: ", w.Result().StatusCode)
//	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
//
//}

//func Test_inDBRepository_GetURL(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	dbMock := mock_usecase.NewMockRepository(ctrl)
//
//	cfg := config.NewConfig()
//
//	hd := handler.NewHandler(dbMock, cfg, nil)
//
//	w := httptest.NewRecorder()
//	r := httptest.NewRequest("GET", cfg.BaseURL+"/myshort", bytes.NewBuffer([]byte("")))
//
//	dbMock.EXPECT().GetURL(gomock.Any(), "myshort").Return(entity.URL{}, false)
//	hd.GetShortURLByIDHandler(w, r)
//
//	hd.Repository.GetURL(context.Background(), "myshort")
//
//	fmt.Println("StatusCode ", w.Result().StatusCode)
//
//}

func Test_inDBRepository_Ping(t *testing.T) {
	l := logger.New("debug")
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbMock := mock_usecase.NewMockRepository(ctrl)

	//cfg := config.NewConfig()
	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	hd := handler.NewHandler(dbMock, &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		SecretKey:     secretKey,
	}, nil, l)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", hd.Cfg.BaseURL+"/ping", nil)
	defer r.Response.Body.Close()
	defer r.Body.Close()
	defer w.Result().Body.Close()

	dbMock.EXPECT().Ping(ctx, l).Return(nil)
	hd.Ping(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func Test_inDBRepository_Delete(t *testing.T) {
	l := logger.New("debug")
	//ctx := context.Background()

	ctrl := gomock.NewController(t)
	dbMock := mock_usecase.NewMockRepository(ctrl)

	//cfg := config.NewConfig()

	workers := worker.NewWorkerPool(l, dbMock)

	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	hd := handler.NewHandler(dbMock, &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		SecretKey:     secretKey,
	}, workers, l)

	query := []string{"1", "2", "3"}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(query)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", hd.Cfg.BaseURL+"/api/user/urls", buf)
	defer r.Response.Body.Close()
	defer r.Body.Close()
	defer w.Result().Body.Close()

	cookies.SetCookieHandler(w, r, secretKey)
	cookies := r.Cookies()[0].Value

	dbMock.EXPECT().Delete(l, query, cookies).Return(nil)
	hd.DeleteHandler(w, r)

	hd.Repository.Delete(l, query, cookies)

	require.Equal(t, http.StatusAccepted, w.Result().StatusCode)
}
