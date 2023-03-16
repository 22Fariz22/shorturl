package config

import (
	"encoding/hex"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress      string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	BaseURL            string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath    string `env:"FILE_STORAGE_PATH"`
	PprofServerAddress string `env:"PPROF_SERVER_ADDRESS" envDefault:"localhost:8081"`
	SecretKey          []byte
	DatabaseDSN        string `env:"DATABASE_DSN" `
}

func NewConfig() *Config {
	cfg := Config{}

	secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", "", "server address")
	flag.StringVar(&cfg.BaseURL, "b", "", "database address")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "db")
	flag.StringVar(&cfg.PprofServerAddress, "p", "", "pprofserver")

	//flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Parse()

	cfg.SecretKey = secretKey

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
