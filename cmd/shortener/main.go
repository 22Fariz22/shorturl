package main

import (
	"22Fariz22/shorturl/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//h := handler.NewHandler()
	//r := mux.NewRouter()
	//
	//r.HandleFunc("/", h.CreateShortURLHandler)
	//r.HandleFunc("/{id:[0-9]+}", h.GetShortURLByIDHandler)
	//http.ListenAndServe("localhost:8080", r)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handler.CreateShortURLHandler)
	r.Post("/{id}", handler.GetShortURLByIDHandler)

}
