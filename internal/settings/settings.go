package settings

import (
	"flag"
	"fmt"
	"log"
	"sync"

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

func (c *Config) checkFlags() {
	flagPort := flag.String("a", "", "SERVER_ADDRESS")
	flagHost := flag.String("b", "", "BASE_URL")
	flagFileStorage := flag.String("f", "", "FILE_STORAGE_PATH")
	flag.Parse()
	if *flagPort != "" {
		c.Server.Port = *flagPort
	}
	if *flagHost != "" {
		c.Server.Host = *flagHost
	}
	if *flagFileStorage != "" {
		c.Storage.FileStorage.FilePath = *flagFileStorage
	}
}

var (
	currentConfig Config
	once          sync.Once
)

func LoadConfig() Config {
	once.Do(func() {
		if err := env.Parse(&currentConfig.Server); err != nil {
			log.Printf("%+v\n", err)
		}
		if err := env.Parse(&currentConfig.Storage.FileStorage); err != nil {
			log.Printf("%+v\n", err)
		}
		currentConfig.checkFlags()
		fmt.Println("here")
	})
	return currentConfig

}
