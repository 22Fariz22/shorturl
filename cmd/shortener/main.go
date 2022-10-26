package main

import (
	"22Fariz22/shorturl/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", handler.CreateShortUrlHandler)
	r.HandleFunc("/{id:[0-9]+}", handler.GetShortUrlByIdHandler)
	//http.Handle("/", r)

	http.ListenAndServe("localhost:8080", r)

}
