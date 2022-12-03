package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"log"
)

const (
	DefaultServerAddress   = "127.0.0.1:8080"
	DefaultBaseURL         = "http://127.0.0.1:8080"
	DefaultFileStoragePath = "events.json"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"  envDefault:"events.json"`
}

func NewConfig() *Config {
	cfg := &Config{}

	pflag.StringVarP(&cfg.ServerAddress, "server", "a", DefaultServerAddress, "server address")
	pflag.StringVarP(&cfg.BaseURL, "baseurl", "b", DefaultBaseURL, "base URL")
	pflag.StringVarP(&cfg.FileStoragePath, "file", "f", DefaultFileStoragePath, "file storage path")

	if err := env.Parse(cfg); err != nil {
		log.Println(err)
	}

	pflag.Parse()

	return &Config{
		ServerAddress:   cfg.ServerAddress,
		BaseURL:         cfg.BaseURL,
		FileStoragePath: cfg.FileStoragePath,
	}
}
