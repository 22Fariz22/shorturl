package repository

import (
	"context"
	"github.com/22Fariz22/shorturl/model"
)

type Repository interface {
	SaveURL(ctx context.Context, shortID string, longURL string, cook string) error
	GetURL(ctx context.Context, shortID string) (string, bool)
	GetAll(context.Context, string) ([]map[string]string, error)
	Init() error
	Ping(context.Context) error
	RepoBatch(context.Context, string, []model.PackReq) error
}
