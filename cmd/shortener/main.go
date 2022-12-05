package main

import (
	"log"
	"net/http"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/handler"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/repository/file"
	"github.com/22Fariz22/shorturl/repository/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.NewConfig()

	var fileRepo repository.Repository

	if cfg.FileStoragePath != "" {
		fileRepo = file.New(cfg)
	} else {
		fileRepo = memory.New()
	}
	fileRepo.Init()

	r := chi.NewRouter()
	r.Use(handler.DeCompress)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	hd := handler.NewHandler(fileRepo, cfg)
	//r.Use(hd.SetCookieMiddleware)

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)
	r.Get("/api/user/urls", hd.GetAllURL)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

}
