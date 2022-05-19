package storageinterface

type Storage interface {
	SaveURL(URL string, shortURL string)
	GetByShortURL(requiredURL string) (shortURL string, err error)
	GetUserURLs(user string) (userURLS []struct {
		Short_URL string
		URL       string
	}, err error)
	Close() error
}
