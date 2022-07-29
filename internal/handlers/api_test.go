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

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestShortenerHandler_GetByShortURL(t *testing.T) {
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
	nameTest := "GetByShortURLEdnpoint test "
	mockStorage := memorystorage.GenMockStorage()
	s := createMockServer(mockStorage, getMockWorkersPool())

	rStorage := reflect.ValueOf(mockStorage).Interface().(*memorystorage.URLStorage)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
			response := execRequest(request, s)
			// check status code
			assert.Equal(t, tt.want.code, response.Code, "wrong answer code")
			r := response.Result()
			defer r.Body.Close()
			// check Location in header
			assert.Equal(t, tt.want.locationHeader, r.Header.Get("Location"), "wrong answer code")

		})
	}
}

func TestShortenerHandler_PostShortURL(t *testing.T) {
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
	nameTest := "PostShortURL test "
	mockStorage := memorystorage.GenMockStorage()

	s := createMockServer(mockStorage, getMockWorkersPool())

	rStorage := reflect.ValueOf(mockStorage).Interface().(*memorystorage.URLStorage)
	for shortURL, webResourse := range rStorage.ShortURLs {
		tests = append(tests, test{
			name: nameTest + strconv.Itoa(testIndex),
			URL:  webResourse.URL,
			want: want{
				code:        http.StatusConflict,
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
			result:      "Wrong URL given -" + "wrong url\n",
			contentType: "text/plain; charset=utf-8",
		},
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.URL)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			response := execRequest(request, s)
			// check status code
			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			r := response.Result()
			defer r.Body.Close()
			assert.Equal(t, tt.want.code, r.StatusCode, "wrong answer code")

			assert.Equal(t, tt.want.result, string(resBody), "wrong result")

			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"), "wrong content-type")

		})
	}
}

func TestShortenerHandler_SaveJSON(t *testing.T) {
	contentType := "application/json; charset=utf-8"
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
	nameTest := "SaveJSON test "
	mockStorage := memorystorage.GenMockStorage()
	s := createMockServer(mockStorage, getMockWorkersPool())

	rStorage := reflect.ValueOf(mockStorage).Interface().(*memorystorage.URLStorage)
	for shortURL, webResourse := range rStorage.ShortURLs {
		tests = append(tests, test{
			name: nameTest + strconv.Itoa(testIndex),
			URL:  webResourse.URL,
			want: want{
				code: http.StatusConflict,
				result: util.StuctEncode(struct {
					Result string `json:"result"`
				}{Result: fmt.Sprintf("%s/%s", host, shortURL)}),
				contentType: contentType,
			},
		})
		testIndex += 1
	}
	tests = append(tests, test{
		name: nameTest + strconv.Itoa(testIndex) + " wrong data",
		URL:  "wrong url",
		want: want{
			code: http.StatusBadRequest,
			result: util.StuctEncode(struct {
				Result string `json:"api_error"`
			}{Result: "Wrong URL given -" + "wrong url"}),
			contentType: contentType,
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(util.StuctEncode(struct {
				URL string `json:"url"`
			}{URL: tt.URL}))

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
			response := execRequest(request, s)
			// check status code
			r := response.Result()
			defer r.Body.Close()
			assert.Equal(t, tt.want.code, r.StatusCode, "wrong answer code")

			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want.result, string(resBody), "wrong result")

			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"), "wrong content-type")

		})
	}
}
