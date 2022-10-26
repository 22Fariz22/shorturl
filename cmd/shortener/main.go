package main

import (
	"22Fariz22/shorturl/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	//m := map[string]string{}
	h := handler.NewHandler()
	r := mux.NewRouter()

	r.HandleFunc("/", h.CreateShortUrlHandler)
	r.HandleFunc("/{id:[0-9]+}", h.GetShortUrlByIdHandler)
	//http.Handle("/", r)

	http.ListenAndServe("localhost:8080", r)

}
