package storageinterface

import "context"

type Storage interface {
	SaveURL(URL string, shortURL string, userID string)
	GetByShortURL(requiredURL string) (shortURL string, err error)
	GetUserURLs(user string) (userURLS []struct {
		ShortURL string
		URL      string
	}, err error)
	Close() error
	CheckConnection(ctx context.Context) error
}
