package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	DefaultServerAddress   = "127.0.0.1:8080"
	DefaultBaseURL         = "http://127.0.0.1:8080"
	DefaultFileStoragePath = "misc/events.json"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"misc/events.json"`
}

func NewConnectorConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	return &Config{
		ServerAddress:   cfg.ServerAddress,
		BaseURL:         cfg.BaseURL,
		FileStoragePath: cfg.FileStoragePath,
	}
}
