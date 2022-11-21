package main

import (
	"22Fariz22/shorturl/handler"
	"22Fariz22/shorturl/handler/config"
	"22Fariz22/shorturl/repo"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConnectorConfig()

	flag.StringVar(&cfg.ServerAddress, "s", config.ServerAddress, "-s to set server address")
	flag.StringVar(&cfg.BaseURL, "b", config.BaseURL, "-b to set base url")
	flag.StringVar(&cfg.FileStoragePath, "f", config.FileStoragePath, "-f to set location storage files")

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = config.FileStoragePath
	}

	//запускаем открытие файла при новом запуске приложении
	producer, err := repo.NewProducer(cfg.FileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	hd := handler.NewHandler(producer)

	hd.RecoverEvents(cfg.FileStoragePath)

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)

	if err = http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		producer.Close()
		log.Fatal("ListenAndServe: ", err)
	}
}
