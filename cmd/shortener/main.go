package main

import (
	"22Fariz22/shorturl/handler"
	"22Fariz22/shorturl/handler/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	cfg := config.NewConnectorConfig()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	hd := handler.NewHandler()

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)

	http.ListenAndServe(cfg.ServerAddress, r)

}
