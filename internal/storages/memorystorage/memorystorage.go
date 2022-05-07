package memorystorage

import (
	"errors"
	"log"
)

type webResourse struct {
	URL     string
	counter int32
}

type URLStorage struct {
	ShortURLs map[string]webResourse
}

func NewStorage() *URLStorage {
	return &URLStorage{
		ShortURLs: make(map[string]webResourse),
	}
}

func (st *URLStorage) SaveURL(URL string, shortURL string) {

	st.ShortURLs[shortURL] = webResourse{
		URL:     URL,
		counter: 0}
	log.Println(shortURL, st.ShortURLs[shortURL])

}

func (st *URLStorage) GetByShortURL(requiredURL string) (shortURL string, err error) {
	webResourse, ok := st.ShortURLs[requiredURL]
	if ok {
		webResourse.counter += 1
		st.ShortURLs[requiredURL] = webResourse

		log.Println(st.ShortURLs[requiredURL])

		return webResourse.URL, nil
	} else {
		err = errors.New(requiredURL + " does not exist")
		return
	}
}

func (st *URLStorage) Close() error {
	log.Println("Memory storage closed")
	return nil
}
