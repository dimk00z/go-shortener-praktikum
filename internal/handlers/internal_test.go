package handlers_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/stretchr/testify/assert"
)

func TestShortenerHandler_GetStats(t *testing.T) {
	mockStorage := memorystorage.GenMockStorage()
	defer mockStorage.Close()
	wp := getMockWorkersPool()
	s := createMockServer(mockStorage, wp)
	endpointURL := "/api/internal/stats"
	req, _ := http.NewRequest("GET", endpointURL, nil)
	response := execRequest(req, s).Result()
	defer response.Body.Close()
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusForbidden, response.StatusCode, "wrong answer code")
	assert.Equal(t, "{\"api_error\":\"Trusted network should be given\"}\n", string(resBody))
	assert.Equal(t, "application/json; charset=utf-8",
		response.Header.Get("Content-Type"))
}
