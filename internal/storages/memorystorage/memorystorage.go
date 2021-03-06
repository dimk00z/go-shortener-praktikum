package memorystorage

import (
	"context"
	"errors"
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
)

type webResourse struct {
	URL       string
	counter   int32
	isDeleted bool
}

type UserURL struct {
	ShortURL string
	URL      string
}

type URLStorage struct {
	ShortURLs map[string]webResourse
	UsersData map[string][]UserURL
}

func NewStorage() *URLStorage {
	return &URLStorage{
		ShortURLs: make(map[string]webResourse),
		UsersData: make(map[string][]UserURL),
	}
}

func (st *URLStorage) SaveURL(URL string, shortURL string, userID string) (err error) {
	if _, ok := st.ShortURLs[shortURL]; ok {
		log.Println(URL, " has been already saved")
		return storageerrors.ErrURLAlreadySave
	}
	st.ShortURLs[shortURL] = webResourse{
		URL:       URL,
		counter:   0,
		isDeleted: false}
	log.Println(shortURL, st.ShortURLs[shortURL])
	if _, ok := st.UsersData[userID]; !ok {
		st.UsersData[userID] = make([]UserURL, 0)
	}
	st.UsersData[userID] = append(st.UsersData[userID], UserURL{
		URL:      URL,
		ShortURL: shortURL,
	})
	return
}

func (st *URLStorage) GetByShortURL(requiredURL string) (URL string, err error) {
	webResourse, ok := st.ShortURLs[requiredURL]
	if ok {
		webResourse.counter += 1
		st.ShortURLs[requiredURL] = webResourse

		log.Println(st.ShortURLs[requiredURL])
		if webResourse.isDeleted {
			err = shortenererrors.ErrURLDeleted
		}
		return webResourse.URL, err
	}
	err = shortenererrors.ErrURLNotFound
	return

}

func (st *URLStorage) Close() error {
	log.Println("Memory storage closed")
	return nil
}

func (st *URLStorage) CheckConnection(ctx context.Context) error {
	return errors.New("wrong storage type")
}

func (st *URLStorage) GetUserURLs(user string) (result models.UserURLs, err error) {
	userURLS, ok := st.UsersData[user]
	result = make([]models.UserURL, len(userURLS))
	if !ok {
		return result, errors.New("no data fo user: " + user)
	}
	for index, userURL := range userURLS {
		result[index] = models.UserURL{ShortURL: userURL.ShortURL,
			URL: userURL.URL}
	}

	log.Println(user, result)
	return
}

func (st *URLStorage) SaveBatch(
	batch models.BatchURLs,
	user string) (result models.BatchShortURLs, err error) {
	result = make(models.BatchShortURLs, len(batch))
	for index, row := range batch {
		st.SaveURL(row.OriginalURL, row.ShortURL, user)
		result[index].CorrelationID = row.CorrelationID
		result[index].ShortURL = row.ShortURL
	}
	return result, err
}

func (st *URLStorage) DeleteBatch(ctx context.Context,
	batch models.BatchForDelete, user string) (err error) {
	for _, shortURL := range batch {
		w := st.ShortURLs[shortURL]
		w.isDeleted = false
		st.ShortURLs[shortURL] = w
	}
	return
}
