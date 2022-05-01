package memory_storage

import "github.com/dimk00z/go-shortener-praktikum/internal/storages/storage_interface"

var mockStorage *URLStorage

func GenMockStorage() storage_interface.Storage {
	if mockStorage == nil {
		mockStorage = NewStorage()
		var mockURLs = []string{
			"http://ya.ru/", "https://yandex.ru/", "https://mail.ru/"}
		for _, url := range mockURLs {
			mockStorage.SaveURL(url)
		}
	}
	return mockStorage
}
