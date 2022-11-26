package model

import (
	"github.com/22Fariz22/shorturl/repo"
	"sync"
)

type HandlerModel struct {
	Mu       sync.Mutex
	Urls     map[string]string `json:"urls"`
	Count    int               `json:"count"`
	Producer *repo.Producer
}
