package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/22Fariz22/shorturl/repo"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/repository/memory"

	"github.com/22Fariz22/shorturl/handler"
	"github.com/22Fariz22/shorturl/handler/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
)

func main() {
	var fileRepo repository.Repository
	fileRepo = memory.New()

	cfg := config.NewConnectorConfig()

	flag.StringVar(&ServerAddress, "s", "", "-s to set server address")           //cfg.ServerAddress
	flag.StringVar(&BaseURL, "b", "", "-b to set base url")                       //cfg.BaseURL
	flag.StringVar(&FileStoragePath, "f", "", "-f to set location storage files") //cfg.FileStoragePath

	flag.Parse()

	r := chi.NewRouter()
	r.Use(handler.DeCompress)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if ServerAddress != "" {
		cfg.ServerAddress = ServerAddress
	}
	if BaseURL != "" {
		cfg.BaseURL = BaseURL
	}
	if FileStoragePath != "" {
		cfg.FileStoragePath = FileStoragePath
	}

	//запускаем открытие файла при новом запуске приложении
	producer, err := repo.NewProducer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	hd := handler.NewHandler(fileRepo)

	if BaseURL != "" {
		hd.Producer.Cfg.BaseURL = BaseURL
	}
	if ServerAddress != "" {
		hd.Producer.Cfg.ServerAddress = ServerAddress
	}
	if FileStoragePath != "" {
		hd.Producer.Cfg.FileStoragePath = FileStoragePath
	}

	hd.RecoverEvents(hd.Producer.Cfg.FileStoragePath)

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
		producer.Close()
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}
}
