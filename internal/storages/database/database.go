package database

type DataBaseStorage struct {
	DataSourceName string
}

func NewDataBaseStorage(DataSourceName string) *DataBaseStorage {
	return &DataBaseStorage{
		DataSourceName: DataSourceName,
	}
}

func (st *DataBaseStorage) Close() (err error) {
	return err
}

func (st *DataBaseStorage) GetUserURLs(user string) (result []struct {
	ShortURL string
	URL      string
}, err error) {
	return
}
func (st *DataBaseStorage) GetByShortURL(requiredURL string) (shortURL string, err error) {
	return
}

func (st *DataBaseStorage) SaveURL(URL string, shortURL string, userID string) {
	return
}
