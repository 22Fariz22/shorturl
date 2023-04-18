// Package usecase интерфейс
package usecase

import (
	"context"

	"github.com/22Fariz22/shorturl/internal/entity"
)

// Repository интерфейс для всех методов
type Repository interface {
	SaveURL(context.Context, string, string, string) (string, error)
	GetURL(context.Context, string) (entity.URL, bool)
	GetAll(context.Context, string) ([]map[string]string, error)
	Init() error
	Ping(context.Context) error
	RepoBatch(context.Context, string, []entity.PackReq) error
	Delete([]string, string) error
	Stats(ctx context.Context) (int, int, error)
}
