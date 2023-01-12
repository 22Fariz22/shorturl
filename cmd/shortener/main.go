package main

import (
	"github.com/22Fariz22/shorturl/repository/db"
	"github.com/22Fariz22/shorturl/repository/file"
	"github.com/22Fariz22/shorturl/repository/memory"
	"github.com/22Fariz22/shorturl/worker"
	"log"
	"net/http"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/handler"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.NewConfig()

	var repo repository.Repository

	if cfg.DatabaseDSN != "" {
		repo = db.New(cfg)
	} else if cfg.FileStoragePath != "" {
		repo = file.New(cfg)
	} else {
		repo = memory.New()
	}

	repo.Init()

	workers := worker.NewWorkerPool(repo)
	workers.RunWorkers(10)

	hd := handler.NewHandler(repo, cfg, workers)

	r := chi.NewRouter()

	r.Use(handler.DeCompress)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)
	r.Get("/api/user/urls", hd.GetAllURL)
	r.Get("/ping", hd.Ping)
	r.Post("/api/shorten/batch", hd.Batch)
	r.Delete("/api/user/urls", hd.DeleteHandler)

	/*
		// получение сигнала о прерывании
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)
		<-signalCh
	*/

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	workers.Stop()
}

//  go run cmd/shortener/main.go -d="postgres://postgres:55555@127.0.0.1:5432/dburl"
// go run cmd/shortener/main.go -f=" log.json"

// разобраться с куками
// сделать транзакции в инсерте дб
