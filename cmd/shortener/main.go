package main

import (
	"22Fariz22/shorturl/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	hd := &handler.Handler{}

	r.Get("/", hd.CreateShortURLHandler)
	r.Post("/{id}", hd.GetShortURLByIDHandler)

	http.ListenAndServe(":8080", r)
}
