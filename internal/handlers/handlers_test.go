package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler_GetEndpoint(t *testing.T) {
	shortenerPort := ":8080"
	host := "http://localhost" + shortenerPort
	mockStorage := *storage.GenMockData()
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
	test_index := 1
	//add correct mock data
	for shortURL, webResourse := range mockStorage.ShortURLs {
		tests = append(tests, test{
			name:     "simple test " + strconv.Itoa(test_index),
			shortURL: shortURL,
			want: want{
				code:           http.StatusTemporaryRedirect,
				locationHeader: webResourse.URL,
			},
		})
		test_index += 1
	}
	// add wrong URL
	tests = append(tests, test{
		name:     "simple test " + strconv.Itoa(test_index+1),
		shortURL: "WrongShortUrl",
		want: want{
			code:           http.StatusNotFound,
			locationHeader: "",
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
			w := httptest.NewRecorder()

			h := NewRootHandler(host)
			h.storage = mockStorage

			h.ServeHTTP(w, request)
			res := w.Result()
			// check status code
			assert.Equal(t, tt.want.code, res.StatusCode, "wrong answer code")

			// check Location in header
			assert.Equal(t, tt.want.locationHeader, res.Header.Get("Location"), "wrong answer code")
		})
	}
}
