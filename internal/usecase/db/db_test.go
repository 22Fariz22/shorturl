package db

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/cookies"
	"github.com/22Fariz22/shorturl/internal/handler"
	mock_usecase "github.com/22Fariz22/shorturl/internal/usecase/mocks"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_inDBRepository_Stats(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_usecase.NewMockRepository(ctrl)

	cfg := config.NewConfig()
	cfg.TrustedSubnet = "127.0.0.1/8"

	hd := handler.NewHandler(dbMock, cfg, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", cfg.BaseURL+"/api/internal/stats", nil)

	dbMock.EXPECT().Stats(gomock.Any()).Return(0, 0, nil)
	hd.Stats(w, r)

	dbMock.Stats(context.Background())

	require.Equal(t, http.StatusForbidden, w.Result().StatusCode)

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
	ctrl := gomock.NewController(t)
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
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", hd.Cfg.BaseURL+"/ping", nil)

	dbMock.EXPECT().Ping(gomock.Any()).Return(nil)
	hd.Ping(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func Test_inDBRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock_usecase.NewMockRepository(ctrl)

	//cfg := config.NewConfig()

	workers := worker.NewWorkerPool(dbMock)

	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	hd := handler.NewHandler(dbMock, &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		SecretKey:     secretKey,
	}, workers)

	query := []string{"1", "2", "3"}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(query)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", hd.Cfg.BaseURL+"/api/user/urls", buf)

	cookies.SetCookieHandler(w, r, secretKey)
	cookies := r.Cookies()[0].Value

	dbMock.EXPECT().Delete(query, cookies).Return(nil)
	hd.DeleteHandler(w, r)

	hd.Repository.Delete(query, cookies)

	require.Equal(t, http.StatusAccepted, w.Result().StatusCode)
}
