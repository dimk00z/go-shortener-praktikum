package storagedi

import (
	"sync"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/database"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/filestorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

var (
	st   storageinterface.Storage
	once sync.Once
)

func GetStorage(l *logger.Logger, storageConfig config.Storage) storageinterface.Storage {
	once.Do(func() {

		if storageConfig.DataSourceName != "" {
			st = database.NewDataBaseStorage(l, storageConfig)
			return
		}
		if storageConfig.FilePath != "" {
			st = filestorage.NewFileStorage(l, storageConfig.FilePath)
			return
		}
		st = memorystorage.NewStorage(l)
	})
	return st
}
