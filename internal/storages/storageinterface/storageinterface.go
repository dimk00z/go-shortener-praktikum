package storageinterface

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/internal/models"
)

type Storage interface {
	SaveURL(URL string, shortURL string, userID string)
	GetByShortURL(requiredURL string) (URL string, err error)
	GetUserURLs(user string) (userURLS []struct {
		ShortURL string
		URL      string
	}, err error)
	Close() error
	CheckConnection(ctx context.Context) error
	SaveBatch(
		batch models.BatchURLs,
		user string) (result models.BatchShortURLs, err error)
}
