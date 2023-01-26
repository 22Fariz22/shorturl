package entity

type URL struct {
	Cookies       string
	ID            string
	LongURL       string
	CorrelationID string
	Deleted       bool
}

type PackReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}

type PackResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
