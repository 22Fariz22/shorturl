package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	ServerAddress   = "127.0.0.1:8080"
	BaseURL         = "http://127.0.0.1:8080"
	FileStoragePath = "storage/events.json"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage/events.json"`
}

//envDefault:"storage/events.json"
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
