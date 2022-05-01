package storagedi

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/file_storage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memory_storage"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storage_interface"
)

var st storage_interface.Storage

func GetStorage(storageConfig settings.StorageConfig) storage_interface.Storage {
	if st != nil {
		return st
	}
	if storageConfig.FileStorage.FilePath != "" {
		st = file_storage.NewFileStorage(settings.LoadConfig().Storage.FileStorage.FilePath)
	} else {
		st = memory_storage.NewStorage()
	}
	return st
}
