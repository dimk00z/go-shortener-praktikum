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
	shortUrls map[string]webResourse
}

func NewStorage() *UrlsStorage {
	return &UrlsStorage{
		shortUrls: make(map[string]webResourse),
	}
}

func (st UrlsStorage) SaveUrl(url string) (shortUrl string) {

	shortUrl = util.GetMD5Hash(url, 4)
	st.shortUrls[shortUrl] = webResourse{
		url:     url,
		counter: 0}
	log.Println(st.shortUrls[shortUrl])
	return
}

func (st UrlsStorage) GetByShortUrl(requiredUrl string) (shortUrl string, err error) {
	webResourse, ok := st.shortUrls[requiredUrl]
	if ok {
		webResourse.counter += 1
		st.shortUrls[requiredUrl] = webResourse

		log.Println(st.shortUrls[requiredUrl])

		return webResourse.url, nil
	} else {
		err = errors.New(requiredUrl + " does not exist")
		return
	}
}
