package settings

import (
	"log"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Port string `env:"SERVER_ADDRESS" envDefault:"8080"`
	Host string `env:"BASE_URL" envDefault:"localhost"`
}
type Config struct {
	Server ServerConfig
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
	return &cfg

}
