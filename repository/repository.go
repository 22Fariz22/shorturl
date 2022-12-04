package repository

type Repository interface {
	SaveURL(shortID string, longURL string, cook string) error
	GetURL(shortID string) (string, bool)
	GetAll(string) []map[string]string
	Init() error
}
