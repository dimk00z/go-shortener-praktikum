package storage

var st *URLStorage

func GetStorage() (*URLStorage, error) {
	if st != nil {
		return st, nil
	}
	var err error
	st = NewStorage()
	return st, err
}
