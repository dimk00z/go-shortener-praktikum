package memorystorage

import (
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

var (
	mockStorage *URLStorage
)

const (
	userUUID string = "7f69f562-e035-41cf-a07e-14e5606f4fbf"
)

func GenMockStorage() storageinterface.Storage {

	if mockStorage == nil {
		mockStorage = NewStorage(logger.New("debug"))
		var mockURLs = []string{
			"http://ya.ru/", "https://yandex.ru/", "https://mail.ru/"}
		for _, url := range mockURLs {
			err := mockStorage.SaveURL(url, util.ShortenLink(url), userUUID)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return mockStorage
}
