package repository

import "context"

type Repository interface {
	SaveURL(ctx context.Context, shortID string, longURL string, cook string) error
	GetURL(ctx context.Context, shortID string) (string, bool)
	GetAll(context.Context, string) []map[string]string
	Init() error
	Ping(context.Context) error
}
