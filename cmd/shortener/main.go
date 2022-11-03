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

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})
	http.ListenAndServe(":8080", r)
}
