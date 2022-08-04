package storagedi

import (
	"sync"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/database"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/filestorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
)

var (
	st   storageinterface.Storage
	once sync.Once
)

func GetStorage(storageConfig config.Storage) storageinterface.Storage {
	once.Do(func() {

		if storageConfig.DataSourceName != "" {
			st = database.NewDataBaseStorage(storageConfig)
			return
		}
		if storageConfig.FilePath != "" {
			st = filestorage.NewFileStorage(storageConfig.FilePath)
			return
		}
		st = memorystorage.NewStorage()
	})
	return st
}
