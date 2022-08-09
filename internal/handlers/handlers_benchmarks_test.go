package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

var testsURLs = [3]string{"https://practicum.yandex.ru/", "https://www.google.com/", "https://go.dev/"}

func Fibo(n int) int {
	if n <= 1 {
		return n
	}
	return Fibo(n-2) + Fibo(n-1)
}

func BenchmarkShortenerHandler_GetByShortURL(b *testing.B) {

	url := testsURLs[0]
	shortURL := util.ShortenLink(url)
	mockStorage := memorystorage.GenMockStorage()
	err := mockStorage.SaveURL(url, shortURL, mockUserID)
	if err != nil {
		b.Fatalf("Wrong status code: %v", err)
	}
	s := createMockServer(mockStorage, getMockWorkersPool())
	b.ResetTimer() // reset all timers
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		request := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)

		rr := httptest.NewRecorder()
		b.StartTimer()
		s.Router.ServeHTTP(rr, request)
		if rr.Code != http.StatusTemporaryRedirect {
			b.Fatalf("Wrong status code: %v", rr.Code)
		}
	}
}

func BenchmarkShortenerHandler_SaveJSON(b *testing.B) {
	b.ResetTimer() // reset all timers
	url := testsURLs[0]
	mockStorage := memorystorage.GenMockStorage()
	s := createMockServer(mockStorage, getMockWorkersPool())

	for i := 0; i < b.N; i++ {

		b.StopTimer()
		body := strings.NewReader(util.StuctEncode(struct {
			URL string `json:"url"`
		}{URL: url}))
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", body)

		rr := httptest.NewRecorder()
		b.StartTimer()
		s.Router.ServeHTTP(rr, request)

		if rr.Code != http.StatusConflict && rr.Code != http.StatusCreated {
			b.Fatalf("Wrong status code: %v", rr.Code)
		}

	}
}

func BenchmarkShortenerHandler_GetUserURLs(b *testing.B) {
	mockStorage := memorystorage.GenMockStorage()
	s := createMockServer(mockStorage, getMockWorkersPool())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		rr := httptest.NewRecorder()
		b.StartTimer()
		s.Router.ServeHTTP(rr, request)
		if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
			b.Fatalf("Wrong status code: %v", rr.Code)
		}
	}
}
