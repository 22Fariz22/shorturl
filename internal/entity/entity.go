// Package entity сущности для всего придложения
package entity

// URL сущность для урла
type URL struct {
	Cookies       string
	ID            string
	LongURL       string
	CorrelationID string
	Deleted       bool
}

//PackReq запрос для json
type PackReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}

//PackResponse ответ для json
type PackResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
