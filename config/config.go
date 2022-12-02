package config

import (
	"github.com/spf13/pflag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	DefaultServerAddress   = "127.0.0.1:8080"
	DefaultBaseURL         = "http://127.0.0.1:8080"
	DefaultFileStoragePath = "events.json"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"  envDefault:"events.json"` //
}

func NewConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	var (
		servAddr string
		bURL     string
		filest   string
	)

	pflag.StringVarP(&servAddr, "server", "a", DefaultServerAddress, "server address")
	pflag.StringVarP(&bURL, "baseurl", "b", DefaultBaseURL, "base URL")
	pflag.StringVarP(&filest, "file", "f", DefaultFileStoragePath, "file storage path")

	pflag.Parse()

	//fmt.Println("servAddr", servAddr)
	//fmt.Println("bURL", bURL)
	//fmt.Println("filest", filest)
	//fmt.Printf("конф %+v\n", cfg)
	//
	//foo, ok := os.LookupEnv("SERVER_ADDRESS")
	//fmt.Println(foo, ok)

	return &Config{
		ServerAddress:   servAddr,
		BaseURL:         bURL,
		FileStoragePath: filest,
	}
}
