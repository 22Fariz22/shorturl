package main

import (
	"github.com/22Fariz22/shorturl/handler"
	"github.com/22Fariz22/shorturl/handler/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/22Fariz22/shorturl/repository/file"
	"github.com/22Fariz22/shorturl/repository/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

var (
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
)

func flagParse() {

}

func main() {
	cfg := config.NewConnectorConfig()

	var fileRepo repository.Repository

	// При отсутствии переменной окружения или
	//при её пустом значении вернитесь к хранению сокращённых URL в памяти.
	if cfg.FileStoragePath != "" {
		fileRepo = file.New()
		fileRepo.Init()
	} else {
		fileRepo = memory.New()
	}

	//flag.StringVar(&ServerAddress, "s", "", "-s to set server address")           //cfg.ServerAddress
	//flag.StringVar(&BaseURL, "b", "", "-b to set base url")                       //cfg.BaseURL
	//flag.StringVar(&FileStoragePath, "f", "", "-f to set location storage files") //cfg.FileStoragePath
	//
	//flag.Parse()

	//if ServerAddress != "" {
	//	cfg.ServerAddress = ServerAddress
	//}
	//if BaseURL != "" {
	//	cfg.BaseURL = BaseURL
	//}
	//if FileStoragePath != "" {
	//	cfg.FileStoragePath = FileStoragePath
	//}

	r := chi.NewRouter()
	r.Use(handler.DeCompress)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	hd := handler.NewHandler(fileRepo)

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}
}
