package memorystorage

import "github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"

var mockStorage *URLStorage

func GenMockStorage() storageinterface.Storage {
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