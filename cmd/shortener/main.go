// Package main точка запуска всего приложения и конфигурации
package main

import (
	"github.com/22Fariz22/shorturl/internal/app"
	"github.com/22Fariz22/shorturl/internal/config"
	"log"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// main точка запуска приложения
func main() {

	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	//cfg запускает конфигурацию
	cfg := config.NewConfig()

	// запускает приложение с учетом конфигурации
	app.Run(cfg)
	//os.Exit(1)

}

//  go run cmd/shortener/main.go -d="postgres://postgres:55555@127.0.0.1:5432/dburl"
