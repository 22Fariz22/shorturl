package model

type URL struct {
	Cookies string
	ID      string
	LongURL string
}

type OwnerID struct {
	ownerID map[string][]string
}
