package main

import (
	"22Fariz22/shorturl/handler"
	"22Fariz22/shorturl/handler/config"
	"22Fariz22/shorturl/repo"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cfg := config.NewConnectorConfig()
	fileName := cfg.FileStoragePath

	if fileName == "" {
		fileName = "storage/events.json"
	}

	//запускаем открытие файла при новом запуске приложении
	producer, err := repo.NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}

	hd := handler.NewHandler(producer)

	hd.RecoverEvents()

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)

	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		producer.Close()
		log.Fatal("ListenAndServe: ", err)
	}
}
