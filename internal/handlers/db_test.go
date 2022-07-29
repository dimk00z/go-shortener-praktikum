package handlers_test

import (
	"net/http"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/stretchr/testify/assert"
)

func TestShortenerHandler_PingDB(t *testing.T) {
	mockStorage := memorystorage.GenMockStorage()
	defer mockStorage.Close()
	wp := getMockWorkersPool()
	defer wp.Close()
	s := createMockServer(mockStorage, wp)
	req, _ := http.NewRequest("GET", "/ping", nil)
	response := execRequest(req, s)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "wrong answer code")

	assert.Equal(t, "{\"message\":\"wrong storage type\"}\n", response.Body.String())
	assert.Equal(t, "application/json; charset=utf-8",
		response.Result().Header.Get("Content-Type"))

}
