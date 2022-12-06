package config

import (
	"encoding/hex"
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
	SecretKey       []byte
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

	//var err error
	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		ServerAddress:   cfg.ServerAddress,
		BaseURL:         cfg.BaseURL,
		FileStoragePath: cfg.FileStoragePath,
		SecretKey:       secretKey,
	}
}
