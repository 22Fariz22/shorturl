// Package usecase интерфейс
package usecase

import (
	"context"
	"github.com/22Fariz22/shorturl/pkg/logger"

	"github.com/22Fariz22/shorturl/internal/entity"
)

// Repository интерфейс для всех методов
type Repository interface {
	SaveURL(context.Context, logger.Interface, string, string, string) (string, error)
	GetURL(context.Context, logger.Interface, string) (entity.URL, bool)
	GetAll(context.Context, logger.Interface, string) ([]map[string]string, error)
	Init(logger.Interface) error
	Ping(context.Context, logger.Interface) error
	RepoBatch(context.Context, logger.Interface, string, []entity.PackReq) error
	Delete(logger.Interface, []string, string) error
	Stats(ctx context.Context, l logger.Interface) (int, int, error)
}
