package storagedi

import (
	"sync"

	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/database"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/filestorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
)

var (
	st   storageinterface.Storage
	once sync.Once
)

func GetStorage(storageConfig settings.StorageConfig) (st storageinterface.Storage) {
	once.Do(func() {
		if storageConfig.DBStorage.DataSourceName != "" {
			st = database.NewDataBaseStorage(storageConfig.DBStorage)
			return
		}
		if storageConfig.FileStorage.FilePath != "" {
			st = filestorage.NewFileStorage(storageConfig.FileStorage.FilePath)
			return
		}
		st = memorystorage.NewStorage()
	})
	return st
}
