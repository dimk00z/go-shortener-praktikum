package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
)

type webResourse struct {
	URL     string `json:"url"`
	Counter int32  `json:"counter"`
}

type UserURL struct {
	ShortURL string
	URL      string
}

type FileStorage struct {
	fileName  string                 `json:"-"`
	ShortURLs map[string]webResourse `json:"short_urls"`
	UsersData map[string][]UserURL   `json:"users_data"`
}

func NewFileStorage(filename string) (st *FileStorage) {
	storage := &FileStorage{
		ShortURLs: make(map[string]webResourse),
		UsersData: make(map[string][]UserURL),
		fileName:  filename,
	}
	storage.load()
	return storage
}

func (st *FileStorage) load() {
	var err error
	file, err := os.OpenFile(st.fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&st); err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", st)
	log.Println("Loaded from", st.fileName)
}

func (st *FileStorage) SaveURL(URL string, shortURL string, userID string) {

	if _, ok := st.ShortURLs[shortURL]; ok {
		return
	}
	wb := webResourse{
		URL:     URL,
		Counter: 0,
	}
	st.ShortURLs[shortURL] = wb
	log.Println(shortURL, st.ShortURLs[shortURL])

	if _, ok := st.UsersData[userID]; !ok {
		st.UsersData[userID] = make([]UserURL, 0)
	}
	st.UsersData[userID] = append(st.UsersData[userID], UserURL{
		URL:      URL,
		ShortURL: shortURL,
	})

	err := st.updateFile()
	if err != nil {
		log.Println(err)
	}

}
func (st *FileStorage) GetByShortURL(requiredURL string) (URL string, err error) {
	webResourse, ok := st.ShortURLs[requiredURL]
	if ok {
		webResourse.Counter += 1
		st.ShortURLs[requiredURL] = webResourse

		log.Println(st.ShortURLs[requiredURL])

		return webResourse.URL, nil
	} else {
		err = errors.New(requiredURL + " does not exist")
		return
	}
}
func (st *FileStorage) updateFile() error {
	file, err := os.OpenFile(st.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.Encode(&st)
	err = file.Sync()
	if err != nil {
		log.Println(err)
	}
	return nil
}
func (st *FileStorage) Close() (err error) {
	err = st.updateFile()
	log.Println("Filestorage closed correctly")

	return err
}

func (st *FileStorage) CheckConnection(ctx context.Context) error {
	return errors.New("wrong storage type")
}

func (st *FileStorage) GetUserURLs(user string) (result []struct {
	ShortURL string
	URL      string
}, err error) {
	userURLS, ok := st.UsersData[user]
	result = make([]struct {
		ShortURL string
		URL      string
	}, len(userURLS))
	if !ok {
		return result, errors.New("no data fo user: " + user)
	}
	for index, userURL := range userURLS {
		result[index] = struct {
			ShortURL string
			URL      string
		}{ShortURL: userURL.ShortURL,
			URL: userURL.URL}
	}

	log.Println(user, result)

	return
}
