package repository

import "github.com/22Fariz22/shorturl/model"

type Repository interface {
	SaveURL(shortID string, longURL string) error
	GetURL(shortID string) (string, bool)
	GetAll() []model.URL
	Init() error
}
