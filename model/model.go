package model

type URL struct {
	Cookies       string
	ID            string
	LongURL       string
	CorrelationID string
}

type OwnerID struct {
	ownerID map[string][]string
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
