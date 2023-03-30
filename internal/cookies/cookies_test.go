package cookies

import (
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_writeEncrypted(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test writeEncrypted",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			if err := writeEncrypted(w,
				http.Cookie{
					Name:     "exampleCookie",
					Value:    "Hello Zoë!",
					Path:     "/",
					MaxAge:   3600,
					HttpOnly: true,
					Secure:   false,
					SameSite: http.SameSiteLaxMode,
				},
				secretKey, req); (err != nil) != tt.wantErr {
				t.Errorf("writeEncrypted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCookieHandler(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			"set cookie",
		},
	}
	for _, tt := range tests {
		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {

			SetCookieHandler(w, req, secretKey)
		})
	}
}

func TestGetCookieHandler(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "get cookie",
		},
	}
	for _, tt := range tests {
		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			GetCookieHandler(w, req, secretKey)
		})
	}
}

func Test_read(t *testing.T) {

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "read",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		SetCookieHandler(w, req, secretKey)

		name := req.Cookies()[0].Value

		t.Run(tt.name, func(t *testing.T) {
			got, err := read(req, name)
			if (err != nil) != tt.wantErr {
				t.Errorf("read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readEncrypted(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "readEncr",
			args:    args{name: "123"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
		if err != nil {
			log.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		SetCookieHandler(w, req, secretKey)
		name := req.Cookies()[0].Value

		t.Run(tt.name, func(t *testing.T) {
			got, err := readEncrypted(req, name, secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("readEncrypted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readEncrypted() got = %v, want %v", got, tt.want)
			}
		})
	}
}
