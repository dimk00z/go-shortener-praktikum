package handlers_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/stretchr/testify/assert"
)

func TestShortenerHandler_GetUserURLs(t *testing.T) {
	mockStorage := memorystorage.GenMockStorage()
	defer mockStorage.Close()
	wp := getMockWorkersPool()
	s := createMockServer(mockStorage, wp)

	req, err := http.NewRequest("GET", "/api/user/urls", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := execRequest(req, s).Result()
	defer response.Body.Close()
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNoContent, response.StatusCode, "wrong answer code")
	assert.Equal(t, "[]\n", string(resBody))
	assert.Equal(t, "application/json; charset=utf-8",
		response.Header.Get("Content-Type"))

}
