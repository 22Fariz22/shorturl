package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080/"`
}

func NewConnectorConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	return &Config{
		ServerAddress: cfg.ServerAddress,
		BaseURL:       cfg.BaseURL,
	}
}
