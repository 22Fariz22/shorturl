package repository

type Repository interface {
	SaveURL(shortID string, longURL string) error
	GetURL(shortID string) (string, bool)
	Init() error
}
