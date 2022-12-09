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
	Correlation_id string `json:"correlation_id"`
	Original_url   string `json:"original_url"`
	Short_url      string `json:"short_url"`
}

type PackResponse struct {
	Correlation_id string `json:"correlation_id"`
	//Original_url   string `json:"original_url"`
	Short_url string `json:"short_url"`
}
