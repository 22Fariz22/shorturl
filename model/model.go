package model

import (
	"github.com/22Fariz22/shorturl/repository"
)

//type HandlerModel struct {
//	Mu       sync.Mutex
//	Urls     map[string]string `json:"urls"`
//	Count    int               `json:"count"`
//	Consumer *repo.Producer
//}

type HandlerModel struct {
	Repository repository.Repository
}
