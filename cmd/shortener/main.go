package main

import (
	"github.com/22Fariz22/shorturl/internal/app"
	"github.com/22Fariz22/shorturl/internal/config"
)

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}

//  go run cmd/shortener/main.go -d="postgres://postgres:55555@127.0.0.1:5432/dburl"
