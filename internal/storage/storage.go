package storage

import (
	"errors"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type WebResourse struct {
	url     string
	counter int32
}

type UrlsStorage struct {
	shortUrls map[string]WebResourse
}

func NewStorage() *UrlsStorage {
	return &UrlsStorage{
		shortUrls: make(map[string]WebResourse),
	}
}

func (st UrlsStorage) SaveUrl(url string) (shortUrl string) {

	shortUrl = util.GetMD5Hash(url, 4)
	st.shortUrls[shortUrl] = WebResourse{
		url:     url,
		counter: 0}
	return
}

func (st UrlsStorage) GetByShortUrl(requiredUrl string) (shortUrl string, err error) {
	webResourse, ok := st.shortUrls[requiredUrl]
	if ok {
		return webResourse.url, nil
	} else {
		err = errors.New(requiredUrl + " does not exist")
		return
	}
}
