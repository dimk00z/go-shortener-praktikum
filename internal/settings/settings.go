package settings

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type ServerConfig struct {
	Port string
}
type Config struct {
	Server ServerConfig
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: util.GetEnv("SHORTENER_PORT", "8080"),
		},
	}
}
