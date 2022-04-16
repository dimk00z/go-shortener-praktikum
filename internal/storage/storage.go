package storage

import (
	"errors"
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type webResourse struct {
	url     string
	counter int32
}

type UrlsStorage struct {
	shortURLs map[string]webResourse
}

func NewStorage() *UrlsStorage {
	return &UrlsStorage{
		shortURLs: make(map[string]webResourse),
	}
}

func (st UrlsStorage) SaveURL(url string) (shortURL string) {

	shortURL = util.GetMD5Hash(url, 4)
	st.shortURLs[shortURL] = webResourse{
		url:     url,
		counter: 0}
	log.Println(st.shortURLs[shortURL])
	return
}

func (st UrlsStorage) GetByShortUrl(requiredUrl string) (shortUrl string, err error) {
	webResourse, ok := st.shortURLs[requiredUrl]
	if ok {
		webResourse.counter += 1
		st.shortURLs[requiredUrl] = webResourse

		log.Println(st.shortURLs[requiredUrl])

		return webResourse.url, nil
	} else {
		err = errors.New(requiredUrl + " does not exist")
		return
	}
}
