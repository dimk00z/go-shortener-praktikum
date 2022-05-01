package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memory_storage"
	"github.com/stretchr/testify/assert"
)

func executeRequest(req *http.Request, s *server.ShortenerServer) *http.Response {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr.Result()
}

func TestRootHandler_GetEndpoint(t *testing.T) {
	shortenerPort := ":8080"
	host := "http://localhost" + shortenerPort
	mockStorage := memory_storage.GenMockStorage()
	type want struct {
		code           int
		locationHeader string
	}
	type test struct {
		name     string
		shortURL string
		want     want
	}
	tests := []test{}
	testIndex := 1
	nameTest := "GetEndpoint test "
	rStorage := reflect.ValueOf(mockStorage).Interface().(*memory_storage.URLStorage)
	//add correct mock data
	for shortURL, webResourse := range rStorage.ShortURLs {
		tests = append(tests, test{
			name:     nameTest + strconv.Itoa(testIndex),
			shortURL: shortURL,
			want: want{
				code:           http.StatusTemporaryRedirect,
				locationHeader: webResourse.URL,
			},
		})
		testIndex += 1
	}
	// add wrong URL
	tests = append(tests, test{
		name:     nameTest + strconv.Itoa(testIndex) + " wrong data",
		shortURL: "WrongShortUrl",
		want: want{
			code:           http.StatusNotFound,
			locationHeader: "",
		},
	})

	server := server.NewServer(
		shortenerPort)
	server.MountHandlers(host, mockStorage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
			response := executeRequest(request, server)
			// check status code
			assert.Equal(t, tt.want.code, response.StatusCode, "wrong answer code")

			// check Location in header
			assert.Equal(t, tt.want.locationHeader, response.Header.Get("Location"), "wrong answer code")
			defer response.Body.Close()
		})
	}
}

func TestRootHandler_PostEndpoint(t *testing.T) {
	shortenerPort := ":8080"
	host := "http://localhost" + shortenerPort
	type want struct {
		code        int
		result      string
		contentType string
	}
	type test struct {
		name string
		URL  string
		want want
	}
	tests := []test{}
	testIndex := 1
	nameTest := "PostEndpoint test "
	mockStorage := memory_storage.GenMockStorage()
	rStorage := reflect.ValueOf(mockStorage).Interface().(*memory_storage.URLStorage)
	for shortURL, webResourse := range rStorage.ShortURLs {
		tests = append(tests, test{
			name: nameTest + strconv.Itoa(testIndex),
			URL:  webResourse.URL,
			want: want{
				code:        http.StatusCreated,
				result:      fmt.Sprintf("%s/%s", host, shortURL),
				contentType: "text/plain; charset=utf-8",
			},
		})
		testIndex += 1
	}
	tests = append(tests, test{
		name: nameTest + strconv.Itoa(testIndex) + " wrong data",
		URL:  "wrong url",
		want: want{
			code:        http.StatusBadRequest,
			result:      "Wrong URL given\n",
			contentType: "text/plain; charset=utf-8",
		},
	})

	server := server.NewServer(
		shortenerPort)
	server.MountHandlers(host, mockStorage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.URL)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			response := executeRequest(request, server)
			// check status code
			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.code, response.StatusCode, "wrong answer code")

			assert.Equal(t, tt.want.result, string(resBody), "wrong result")

			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "wrong content-type")

			defer response.Body.Close()

		})
	}
}
