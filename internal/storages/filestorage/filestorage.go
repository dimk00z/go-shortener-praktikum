package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
)

type webResourse struct {
	URL       string `json:"url"`
	Counter   int32  `json:"counter"`
	IsDeleted bool   `json:"is_deleted"`
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

func (st *FileStorage) SaveURL(URL string, shortURL string, userID string) (err error) {

	if _, ok := st.ShortURLs[shortURL]; ok {
		return storageerrors.ErrURLAlreadySave
	}
	wb := webResourse{
		URL:       URL,
		Counter:   0,
		IsDeleted: false,
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

	err = st.updateFile()
	if err != nil {
		log.Println(err)
	}
	return
}

func (st *FileStorage) SaveBatch(
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

func (st *FileStorage) GetByShortURL(requiredURL string) (URL string, err error) {
	webResourse, ok := st.ShortURLs[requiredURL]
	if ok {
		webResourse.Counter += 1
		st.ShortURLs[requiredURL] = webResourse

		log.Println(st.ShortURLs[requiredURL])
		if webResourse.IsDeleted {
			err = shortenererrors.ErrURLDeleted
		}
		return webResourse.URL, err
	}
	err = errors.New(requiredURL + " does not exist")
	return

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

func (st *FileStorage) GetUserURLs(user string) (result models.UserURLs, err error) {
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

func (st *FileStorage) DeleteBatch(ctx context.Context, batch models.BatchForDelete, user string) (err error) {
	for _, shortURL := range batch {
		w := st.ShortURLs[shortURL]
		w.IsDeleted = true
		st.ShortURLs[shortURL] = w
	}
	err = st.updateFile()
	if err != nil {
		log.Println(err)
	}
	return
}
