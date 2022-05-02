package settings

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Port string `env:"SERVER_ADDRESS" envDefault:":8080"`
	Host string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}
type FileStorageConfig struct {
	FilePath string `env:"FILE_STORAGE_PATH"`
}
type StorageConfig struct {
	FileStorage FileStorageConfig
}
type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

func flagGiven(fl string) bool {
	return flag.Lookup(fl) == nil
}

func (c *Config) checkFlags() {
	if flagGiven("a") {
		flag.StringVar(&c.Server.Port, "a", c.Server.Port, "SERVER_ADDRESS")
	}
	if flagGiven("b") {
		flag.StringVar(&c.Server.Host, "b", c.Server.Port, "BASE_URL")
	}
	if flagGiven("f") {
		flag.StringVar(&c.Storage.FileStorage.FilePath, "f", c.Storage.FileStorage.FilePath, "FILE_STORAGE_PATH")
	}
	flag.Parse()
}

var currentConfig *Config

func LoadConfig() *Config {
	if currentConfig != nil {
		return currentConfig
	}
	cfg := new(Config)
	if err := env.Parse(&cfg.Server); err != nil {
		log.Printf("%+v\n", err)
	}
	if err := env.Parse(&cfg.Storage.FileStorage); err != nil {
		log.Printf("%+v\n", err)
	}
	cfg.checkFlags()
	currentConfig = cfg
	return currentConfig

}
