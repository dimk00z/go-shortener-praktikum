package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

func Fibo(n int) int {
	if n <= 1 {
		return n
	}
	return Fibo(n-2) + Fibo(n-1)
}
func BenchmarkShortenerHandler_GetByShortURL(b *testing.B) {
	//
	url := "https://practicum.yandex.ru/"
	shortURL := util.ShortenLink(url)
	mockUser := "7f69f562-e035-41cf-a07e-14e5606f4fbf"
	mockStorage := memorystorage.GenMockStorage()
	mockStorage.SaveURL(url, shortURL, mockUser)
	s := createMockServer(mockStorage, getMockWorkersPool())
	b.ResetTimer() // reset all timers
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		b.StartTimer()
		request := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
		execRequest(request, s)
	}
}

func BenchmarkShortenerHandler_SaveJSON(b *testing.B) {
	b.ResetTimer() // reset all timers

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// TODO add logic
		b.StartTimer()

		// TODO add logic
		Fibo(20)
	}
}

func BenchmarkShortenerHandler_SaveBatch(b *testing.B) {
	b.ResetTimer() // reset all timers

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// TODO add logic
		b.StartTimer()

		// TODO add logic
		Fibo(20)
	}
}
