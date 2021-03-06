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
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/memorystorage"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/stretchr/testify/assert"
)

func TestShortenerAPIHandler_PostEndpoint(t *testing.T) {
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
	mockStorage := memorystorage.GenMockStorage()
	contentType := "application/json; charset=utf-8"

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
	wp := worker.GetWorkersPool(settings.WorkersConfig{WorkersNumber: 2, PoolLength: 10})
	defer wp.Close()
	server := server.NewServer(
		shortenerPort, wp)
	server.MountHandlers(host, mockStorage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(util.StuctEncode(struct {
				URL string `json:"url"`
			}{URL: tt.URL}))

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
			response := executeRequest(request, server)
			// check status code
			assert.Equal(t, tt.want.code, response.StatusCode, "wrong answer code")

			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want.result, string(resBody), "wrong result")

			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"), "wrong content-type")

			defer response.Body.Close()

		})
	}
}
