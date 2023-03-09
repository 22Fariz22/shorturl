// Package config конфигуратор переменных окружения
package config

import (
	"encoding/hex"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

// переменные окружения по умолчанию
const (
	DefaultServerAddress      = "localhost:8080"        // адрес сервера
	DefaultBaseURL            = "http://localhost:8080" // базовый адрес
	DefaultPprofServerAddress = "http://localhost:8081" // адрес сервера для профилирования

	DefaultDatabaseDSN = "" //"postgres://postgres:55555@127.0.0.1:5432/dburl"
)

// Config структура конфига
type Config struct {
	ServerAddress      string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL            string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath    string `env:"FILE_STORAGE_PATH"  `
	PprofServerAddress string `env:"PPROF_SERVER_ADDRESS" envDefault:"localhost:8081"`
	SecretKey          []byte

	DatabaseDSN string `env:"DATABASE_DSN" ` //envDefault:"postgres://postgres:55555@127.0.0.1:5432/dburl"
}

//NewConfig создание конфига
func NewConfig() *Config {
	cfg := &Config{}
	pflag.StringVarP(&cfg.ServerAddress, "server", "a", DefaultServerAddress, "server address")
	pflag.StringVarP(&cfg.BaseURL, "baseurl", "b", DefaultBaseURL, "base URL")
	pflag.StringVarP(&cfg.FileStoragePath, "file", "f", "", "file storage path")
	pflag.StringVarP(&cfg.DatabaseDSN, "databasedsn", "d", "", "databaseDSN")
	pflag.StringVarP(&cfg.PprofServerAddress, "pprof server", "p", DefaultPprofServerAddress, "pprof server address")

	if err := env.Parse(cfg); err != nil {
		log.Println(err)
	}

	pflag.Parse()

	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		ServerAddress:      cfg.ServerAddress,
		BaseURL:            cfg.BaseURL,
		FileStoragePath:    cfg.FileStoragePath,
		SecretKey:          secretKey,
		DatabaseDSN:        cfg.DatabaseDSN,
		PprofServerAddress: cfg.PprofServerAddress,
	}
}
