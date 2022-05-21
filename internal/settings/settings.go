package settings

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Port string `env:"SERVER_ADDRESS" envDefault:":8080"`
	Host string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}
type DBStorageConfig struct {
	DataSourceName string `env:"DATABASE_DSN"`
	SQLScriptsFile string `env:"SERVER_ADDRESS" envDefault:"/scripts/sql/scripts.sql"`
}
type FileStorageConfig struct {
	FilePath string `env:"FILE_STORAGE_PATH"`
}
type StorageConfig struct {
	DBStorage   DBStorageConfig
	FileStorage FileStorageConfig
}

type SecurityConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"SECRET_KEY"`
}
type Config struct {
	Server   ServerConfig
	Storage  StorageConfig
	Security SecurityConfig
}

func (c *Config) checkFlags() {
	flagPort := flag.String("a", "", "SERVER_ADDRESS")
	flagHost := flag.String("b", "", "BASE_URL")
	flagFileStorage := flag.String("f", "", "FILE_STORAGE_PATH")
	flagDBStorage := flag.String("d", "", "DATABASE_DSN")
	flag.Parse()
	if *flagPort != "" {
		c.Server.Port = *flagPort
	}
	if *flagHost != "" {
		c.Server.Host = *flagHost
	}
	if *flagDBStorage != "" {
		c.Storage.DBStorage.DataSourceName = *flagDBStorage
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
		if err := env.Parse(&currentConfig.Storage.DBStorage); err != nil {

			log.Printf("%+v\n", err)
		}
		if err := env.Parse(&currentConfig.Security); err != nil {
			log.Printf("%+v\n", err)
		}
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		//TODO не забыть, что под виндой не оживет из-за разницы в путях
		projectDir := filepath.Dir(filepath.Dir(path))
		currentConfig.Storage.DBStorage.SQLScriptsFile = filepath.Join(projectDir, currentConfig.Storage.DBStorage.SQLScriptsFile)
		currentConfig.checkFlags()
	})
	return currentConfig

}
