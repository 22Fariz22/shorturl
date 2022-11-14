package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseUrl       string `env:"BASE_URL"`
}

func NewConnectorConfig() *Config {
	cfg := &Config{}
	opts := &env.Options{Environment: map[string]string{
		"SERVER_ADDRESS": ":8080",
	}}
	if err := env.Parse(cfg, *opts); err != nil {
		log.Fatal(err)
	}

	return &Config{
		ServerAddress: cfg.ServerAddress,
	}
}

func NewConnectorConfigURL(url string) *Config {
	cfg := &Config{}
	opts := &env.Options{Environment: map[string]string{
		"SERVER_ADDRESS": ":8080",
	}}
	if err := env.Parse(cfg, *opts); err != nil {
		log.Fatal(err)
	}

	return &Config{
		ServerAddress: cfg.ServerAddress,
		BaseUrl:       url,
	}

}
