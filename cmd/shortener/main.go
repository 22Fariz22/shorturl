// Package main точка запуска всего приложения и конфигурации
package main

import (
	"github.com/22Fariz22/shorturl/internal/app"
	"github.com/22Fariz22/shorturl/internal/config"
)

//main точка запуска приложения
func main() {
	//cfg запускает конфигурацию
	cfg := config.NewConfig()

	// запускает приложение с учетом конфигурации
	app.Run(cfg)
}

//  go run cmd/shortener/main.go -d="postgres://postgres:55555@127.0.0.1:5432/dburl"
