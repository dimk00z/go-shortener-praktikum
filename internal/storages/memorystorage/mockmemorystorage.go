package memorystorage

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

var (
	mockStorage *URLStorage
)

const (
	userUUID string = "7f69f562-e035-41cf-a07e-14e5606f4fbf"
)

func GenMockStorage() storageinterface.Storage {

	if mockStorage == nil {
		mockStorage = NewStorage()
		var mockURLs = []string{
			"http://ya.ru/", "https://yandex.ru/", "https://mail.ru/"}
		for _, url := range mockURLs {
			mockStorage.SaveURL(url, util.ShortenLink(url), userUUID)
		}
	}

	return mockStorage
}
