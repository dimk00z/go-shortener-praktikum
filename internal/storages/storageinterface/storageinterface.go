package storageinterface

type Storage interface {
	SaveURL(URL string, shortURL string)
	GetByShortURL(requiredURL string) (shortURL string, err error)
	Close() error
}
