package memorystorage

import (
	"errors"
	"fmt"
	"log"
)

type webResourse struct {
	URL     string
	counter int32
}

type UserURL struct {
	Short_URL string
	URL       string
}

type URLStorage struct {
	ShortURLs map[string]webResourse
	UsersData map[string][]UserURL
}

func NewStorage() *URLStorage {
	return &URLStorage{
		ShortURLs: make(map[string]webResourse),
		UsersData: make(map[string][]UserURL),
	}
}

func (st *URLStorage) SaveURL(URL string, shortURL string, userId string) {
	if _, ok := st.ShortURLs[shortURL]; ok {
		log.Println(URL, " has been already saved")
		return
	}
	st.ShortURLs[shortURL] = webResourse{
		URL:     URL,
		counter: 0}
	log.Println(shortURL, st.ShortURLs[shortURL])
	if _, ok := st.UsersData[userId]; !ok {
		st.UsersData[userId] = make([]UserURL, 0)
	}
	st.UsersData[userId] = append(st.UsersData[userId], UserURL{
		URL:       URL,
		Short_URL: shortURL,
	})

}

func (st *URLStorage) GetByShortURL(requiredURL string) (shortURL string, err error) {
	webResourse, ok := st.ShortURLs[requiredURL]
	if ok {
		webResourse.counter += 1
		st.ShortURLs[requiredURL] = webResourse

		log.Println(st.ShortURLs[requiredURL])

		return webResourse.URL, nil
	} else {
		err = errors.New(requiredURL + " does not exist")
		return
	}
}

func (st *URLStorage) Close() error {
	log.Println("Memory storage closed")
	return nil
}

func (st *URLStorage) GetUserURLs(user string) (result []struct {
	Short_URL string
	URL       string
}, err error) {
	userURLS, ok := st.UsersData[user]
	result = make([]struct {
		Short_URL string
		URL       string
	}, len(userURLS))

	for index, userURL := range userURLS {
		result[index] = struct {
			Short_URL string
			URL       string
		}{Short_URL: userURL.Short_URL,
			URL: userURL.URL}
	}

	result = append(result,
		UserURL{Short_URL: "Short_URL_test1", URL: "URL_test1"},
		UserURL{Short_URL: "Short_URL_test2", URL: "URL_test2"})
	fmt.Println(result)

	if !ok {
		return result, errors.New("no data fo user: " + user)
	}

	return
}
