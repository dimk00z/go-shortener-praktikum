package storagedi

import (
	"sync"

	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/filestorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
)

var (
	st   storageinterface.Storage
	once sync.Once
)

func GetStorage(storageConfig settings.StorageConfig) storageinterface.Storage {
	once.Do(func() {
		if storageConfig.FileStorage.FilePath != "" {
			st = filestorage.NewFileStorage(settings.LoadConfig().Storage.FileStorage.FilePath)
		} else {
			st = memorystorage.NewStorage()
		}
	})
	return st
}
