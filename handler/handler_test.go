package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type Counter struct{
	count int
}

func (c *Counter)CounterFunc()string{
	c.count++
	countStr := strconv.Itoa(c.count)
	return countStr
}

func TestHandler_CreateShortURLHandler(t *testing.T) {
	type want struct {
		statusCode        int
		shortURL string
	}
	tests := []struct {
		name string
		request string
		counter string
		want want
	}{
		{
			name: "simple test #1",
			counter: ,
			want: want{
				statusCode: 201,
				shortURL: "",
			},
		},
	}
	for _,tt := range tests{
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateShortURLHandler(w,request))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			shortURLResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.shortURL, shortURLResult)

		})
	}
}

