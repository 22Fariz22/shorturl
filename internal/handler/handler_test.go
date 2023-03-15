// тестирование
package handler_test

import (
	"encoding/hex"
	"github.com/22Fariz22/shorturl/internal/handler"
	"log"
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

//func _TestHandlerCreateShortURLJSON(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name string
//		in   []byte
//		want string
//	}{
//		{
//			name: "create json",
//			in:   []byte(`{"url":"https://google.ru"}`),
//			want: `{"result":"http://localhost:8080/"}`,
//		},
//	}
//
//	for i := range tests {
//		test := tests[i]
//		t.Run(test.name, func(t *testing.T) {
//			t.Parallel()
//
//			cfg := config.NewConfig()
//
//			repo := memory.New()
//
//			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(test.in))
//			req.Header.Add("Content-Type", "application/json")
//			w := httptest.NewRecorder()
//
//			handler.NewHandler(repo, cfg, nil).CreateShortURLJSON(w, req)
//
//			_, err := io.ReadAll(w.Result().Body)
//
//			require.NoError(t, err)
//			//require.JSONEq(t, string(response), string(response)) // как вставить  want геренрирующий
//		})
//	}
//}

//TestHandler_CreateShortURLHandler тест хэндлера из эндпойнта r.Post("/", hd.CreateShortURLHandler)
//func TestHandler_GetShortURLByIDHandler(t *testing.T) {
//	ctl := gomock.NewController(t)
//	ctl.Finish()
//
//	repo := repoMock.NewMockRepository(ctl)
//
//	ctx := context.Background()
//	short := handler.GenUlid()
//
//	//mockResp := entity.URL{LongURL: "https://ya.ru"}
//
//	expected := entity.URL{LongURL: "https://ya.ru"}
//
//	repo.EXPECT().GetURL(ctx, short).Return(expected, true).Times(1)
//	URL, ok := repo.GetURL(ctx, short)
//
//	require.Equal(t, expected, URL)
//	require.EqualValues(t, true, ok)
//}

//func TestHandler_CreateShortURLHandler(t *testing.T) {
//
//	ctl := gomock.NewController(t)
//	ctl.Finish()
//
//	repo := repoMock.NewMockRepository(ctl)
//
//	ctx := context.Background()
//	short := handler.GenUlid()
//	payload := "https://google.com"
//	cookie := "ABCD12345"
//
//	//mockResp := entity.URL{LongURL: "https://ya.ru"}
//
//	exp := "http://localhost:8080/" + short
//
//	repo.EXPECT().SaveURL(ctx, short, payload, cookie).Return(exp, nil).Times(1)
//	resp, err := repo.SaveURL(ctx, short, payload, cookie)
//
//	require.Equal(t, exp, resp)
//	require.NoError(t, err)
//}
