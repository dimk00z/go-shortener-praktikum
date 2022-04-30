package storage

var mockURLs = []string{
	"http://ya.ru/", "https://yandex.ru/", "https://mail.ru/"}

func GenMockStorage() (*URLStorage, error) {
	storage := NewStorage()
	for _, url := range mockURLs {
		storage.SaveURL(url)
	}
	return storage, nil
}
