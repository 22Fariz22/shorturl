// Package app модуль app запускает репозиторий  с учетом конфигурации, воркер и роутер
package app

import (
	"context"
	"fmt"
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/handler"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"github.com/22Fariz22/shorturl/internal/usecase/db"
	"github.com/22Fariz22/shorturl/internal/usecase/file"
	"github.com/22Fariz22/shorturl/internal/usecase/memory"
	"github.com/22Fariz22/shorturl/internal/worker"
	pb "github.com/22Fariz22/shorturl/pkg/proto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	r.Get("/api/internal/stats", hd.Stats)

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

	go func() {
		if cfg.EnableHTTPS {
			log.Println("start https server.")
			server.ListenAndServeTLS("", "")

		} else {
			if err := http.ListenAndServe(cfg.ServerAddress, r); err != http.ErrServerClosed {
				log.Println("start http server.")
				log.Fatalf("HTTP server ListenAndServe Error: %v", err)
			}
		}
	}()

	// определяем порт для сервера grpc
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	pb.RegisterServicesServer(s, handler.NewGRPCServer(*cfg, hd))

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	//gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	server.Shutdown(ctx)
}
