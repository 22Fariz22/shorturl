package repository

import (
	"github.com/22Fariz22/shorturl/model"
	"net/http"
)

type Repository interface {
	SaveURL(shortID string, longURL string, cook *http.Cookie) error
	GetURL(shortID string) (string, bool)
	GetAll() []model.URL
	Init() error
}
