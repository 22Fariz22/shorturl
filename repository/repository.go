package repository

import (
	"context"
	"github.com/22Fariz22/shorturl/model"
)

type Repository interface {
	SaveURL(context.Context, string, string, string) error
	GetURL(context.Context, string, string) (string, bool)
	GetAll(context.Context, string) ([]map[string]string, error)
	Init() error
	Ping(context.Context) error
	RepoBatch(context.Context, string, []model.PackReq) error
}
