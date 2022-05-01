package filestorage

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type webResourse struct {
	URL     string `json:"url"`
	Counter int32  `json:"counter"`
}

type FileStorage struct {
	fileName  string `json:"-"`
	ShortURLs map[string]webResourse
}

func NewFileStorage(filename string) (st *FileStorage) {
	storage := &FileStorage{
		ShortURLs: make(map[string]webResourse),
		fileName:  filename,
	}
	storage.load()
	return storage
}

func (st *FileStorage) load() (err error) {
	file, err := os.OpenFile(st.fileName, os.O_RDONLY|os.O_CREATE, 0777)
	defer file.Close()
	if err != nil {
		log.Panicln(err)
	}
	loadedData := make(map[string]webResourse)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedData); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("%+v\n", loadedData)
	st.ShortURLs = loadedData
	log.Println("Loaded from", st.fileName)
	return nil
}

func (st *FileStorage) SaveURL(URL string) (shortURL string) {

	shortURL = util.GetMD5Hash(URL)
	if _, ok := st.ShortURLs[shortURL]; ok {
		return shortURL
	}
	wb := webResourse{
		URL:     URL,
		Counter: 0,
	}
	st.ShortURLs[shortURL] = wb
	log.Println(shortURL, st.ShortURLs[shortURL])
	err := st.updateFile()
	if err != nil {
		log.Println(err)
	}
	return shortURL

}
func (st *FileStorage) GetByShortURL(requiredURL string) (shortURL string, err error) {
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
	file, err := os.OpenFile(st.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.Encode(&st.ShortURLs)
	return nil
}
func (st *FileStorage) Close() (err error) {
	err = st.updateFile()
	log.Println("Filestorage closed correctly")

	return err
}
