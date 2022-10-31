package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
Создаем глобальную переменную наш URL
вставляем в POST наш URL.
берем выданный нам id и вставляем в GET.
Смотрим на урл в заголовке Location и сравниваем с исходным
Тесты:
	1)посмотреть, что в посте что-то вернется
		*проверить cтатус-код 201
		*забрать из тела сокращенный урл
	2)Посмотреть, что, если заранее заготовить данные, в локейшне будет заранее заготовленный урл
		*переходим по сокращенному урлу,сравниваем статус-код 307 и сокращенный урл с урл из Location
	3)Совместить: пост + гет

*/
type Fields struct {
	urls map[string]string
}

func NewFields() *Fields {
	return &Fields{
		urls: make(map[string]string),
	}
}

func (f Fields) TestHandler_CreateShortURLHandler(t *testing.T) {
	/* POST
	вставляем свой урл url
	ссравниваю статусКод ,нужен 201
	беру урл и венрнувшийся шортурл,все это вставляю в map[string]string
	*/

	type want struct {
		code int
		url  string
		//response    string
		//contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test POST #1",
			want: want{
				code: 201,
				url:  "https://ya.ru",
			},
		},
	}
	for _, tt := range tests {
		for t.Run(tt.name, func(t *testing.T) {}) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)

			w := httptest.NewRecorder()

			hd := &Handler{}

			h := http.HandlerFunc(hd.CreateShortURLHandler)
			h.ServeHTTP(w, request)

			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Exepted status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal()
			}

			// записываем в глобальную мапу наш url и полученный id для сравнения в GET
			shortUrl := string(resBody)[len("http://localhost:8080/"):]
			f.urls[tt.want.url] = shortUrl
		}
	}
}

func (f Fields) TestHandler_GetShortURLByIDHandler(t *testing.T) {
	type want struct {
		code     int
		location string // f.
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test GET #1",
			want: want{
				code:     307,
				location: "https://ya.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.want.location, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			hd := &Handler{}
			h := http.HandlerFunc(hd.GetShortURLByIDHandler)

			// определяем хендлер
			//h := http.HandlerFunc(GetShortURLByIDHandler)

			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем location из заголовка
			res_location := res.Header.Get("Location")

			// сравниваем location
			if res_location != tt.want.location {
				t.Errorf("Expected location %s, got %s", tt.want.location, res_location)
			}
		})
	}
}
