package main

import (
	"22Fariz22/shorturl/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	h := handler.NewHandler()
	r := mux.NewRouter()

	r.HandleFunc("/", h.CreateShortURLHandler)
	r.HandleFunc("/{id:[0-9]+}", h.GetShortURLByIDHandler)

	http.ListenAndServe("localhost:8080", r)

}
