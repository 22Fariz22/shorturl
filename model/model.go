package model

import "net/http"

type URL struct {
	Cookies *http.Cookie
	ID      string //`json:"short_url"`
	LongURL string //`json:"original_url"`
}

type OwnerURL struct {
	ownerURL map[string][]URL
}

type OwnerID struct {
	ownerID map[string][]string
}
