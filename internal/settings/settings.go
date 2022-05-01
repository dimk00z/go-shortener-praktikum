package settings

import (
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

var currentConfig *Config

func LoadConfig() *Config {
	if currentConfig != nil {
		return currentConfig
	}
	cfg := Config{}
	serverConfig := ServerConfig{}
	if err := env.Parse(&serverConfig); err != nil {
		log.Printf("%+v\n", err)
	}
	cfg.Server = serverConfig
	var fileStorage FileStorageConfig
	if err := env.Parse(&fileStorage); err != nil {
		log.Printf("%+v\n", err)
	}
	cfg.Storage = StorageConfig{
		FileStorage: fileStorage,
	}
	return &cfg

}
