package model

import "sync"

type HandlerModel struct {
	Mu    sync.Mutex
	Urls  map[string]string `json:"urls"` //map[0:http://ya.ru]
	Count int               `json:"count"`
}
