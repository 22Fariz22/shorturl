package main

import (
	"github.com/gorilla/mux"
	"log"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", CreateShortUrlHandler)

	log.Fatal(":8080", nil)
}
