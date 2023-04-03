// Package app модуль app запускает репозиторий  с учетом конфигурации, воркер и роутер
package app

import (
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/handler"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"github.com/22Fariz22/shorturl/internal/usecase/db"
	"github.com/22Fariz22/shorturl/internal/usecase/file"
	"github.com/22Fariz22/shorturl/internal/usecase/memory"
	"github.com/22Fariz22/shorturl/internal/worker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// Run запускает приложение с учетом конфигурации из main и роутеры
func Run(cfg *config.Config) {
	var repo usecase.Repository

	if cfg.DatabaseDSN != "" {
		repo = db.New(cfg)
	} else if cfg.FileStoragePath != "" {
		repo = file.New(cfg)
	} else {
		repo = memory.New()
	}

	repo.Init()

	workers := worker.NewWorkerPool(repo)
	workers.RunWorkers(10)
	defer workers.Stop()

	r := chi.NewRouter()

	hd := handler.NewHandler(repo, cfg, workers)

	r.Use(handler.DeCompress)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Mount("/debug", middleware.Profiler())

	r.Post("/", hd.CreateShortURLHandler)
	r.Get("/{id}", hd.GetShortURLByIDHandler)
	r.Post("/api/shorten", hd.CreateShortURLJSON)
	r.Get("/api/user/urls", hd.GetAllURL)
	r.Get("/ping", hd.Ping)
	r.Post("/api/shorten/batch", hd.Batch)
	r.Delete("/api/user/urls", hd.DeleteHandler)

	go func() {
		if cfg.PprofServerAddress == "" {
			log.Println("pprof server address is empty, skipping")
			return
		}
		err := http.ListenAndServe(cfg.PprofServerAddress, nil)
		if err != nil {
			log.Printf("pprof server error: %s\n", err)
		}
	}()

	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.BaseURL),
	}
	server := &http.Server{
		Addr:      cfg.ServerAddress,
		Handler:   r,
		TLSConfig: manager.TLSConfig(),
	}

	if cfg.EnableHTTPS {
		log.Println("start https server.")
		server.ListenAndServeTLS("", "")

	} else {
		if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
			log.Println("start http server.")
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}

	//for i := range os.Args {
	//	if strings.HasPrefix(os.Args[i], "-s") {
	//		log.Println("start http server.")
	//		server.ListenAndServeTLS("", "")
	//	} else {
	//		if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
	//			log.Println("start https server.")
	//			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	//		}
	//	}
	//}

	//for i := range os.Args {
	//	if strings.HasPrefix(os.Args[i], "-s") {
	//		fmt.Println("has -s")
	//	}
	//}

	//if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
	//	log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	//}
}
