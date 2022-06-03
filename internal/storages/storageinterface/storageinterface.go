package storageinterface

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/internal/models"
)

type Storage interface {
	SaveURL(URL string, shortURL string, userID string) (err error)
	GetByShortURL(requiredURL string) (URL string, err error)
	GetUserURLs(user string) (result models.UserURLs, err error)
	Close() error
	CheckConnection(ctx context.Context) error
	SaveBatch(
		batch models.BatchURLs,
		user string) (result models.BatchShortURLs, err error)
	DeleteBatch(ctx context.Context, batch models.BatchForDelete, user string) (err error)
}
