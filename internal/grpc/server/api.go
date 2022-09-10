package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

var (
	errEmptyURLGiven   = errors.New("empty url was given")
	errWrongURLGiven   = errors.New("wrong url was given")
	errEmptyBatchGiven = errors.New("empty batch was given")
)

func (s *GRPCServer) GetByShortURL(ctx context.Context, request *pb.ShortURLRequest) (*pb.URLResponse, error) {
	var response pb.URLResponse
	response.Code = http.StatusTemporaryRedirect
	s.Service.l.Debug("GetByShortURL - got - " + request.ShortUrl)
	var err error
	response.Url, err = s.Service.st.GetByShortURL(request.ShortUrl)
	if err != nil {

		if errors.Is(err, shortenererrors.ErrURLNotFound) {
			response.Code = http.StatusNotFound
		}
		if errors.Is(err, shortenererrors.ErrURLDeleted) {
			response.Code = http.StatusGone
		}
	}
	return &response, err
}
func (s *GRPCServer) SaveURLFromText(ctx context.Context, request *pb.URLRequest) (*pb.URLFromTextResponse, error) {
	var response pb.URLFromTextResponse
	s.Service.l.Debug("SaveURLFromText - got " + request.Url)
	response.Code = http.StatusOK
	if request.Url == "" {
		response.Code = http.StatusBadRequest
		return &response, errEmptyURLGiven
	}
	if !util.IsURL(request.Url) {
		response.Code = http.StatusBadRequest
		return &response, errWrongURLGiven
	}
	shortURL := util.ShortenLink(request.Url)
	userIDCtx := getUserIDFromMD(ctx)
	err := s.Service.st.SaveURL(request.Url, shortURL, userIDCtx)
	if errors.Is(err, storageerrors.ErrURLAlreadySave) {
		response.Code = http.StatusConflict
	}
	response.ShortUrl = fmt.Sprintf("%s/%s", s.Service.host, shortURL)
	return &response, nil
}

func (s *GRPCServer) SaveBatch(ctx context.Context, request *pb.BatchSaveRequest) (*pb.BatchSaveResponse, error) {
	var response pb.BatchSaveResponse
	if len(request.Urls) == 0 {
		response.Code = http.StatusBadRequest
		return &response, errEmptyBatchGiven
	}
	userIDCtx := getUserIDFromMD(ctx)
	requestData := make(models.BatchURLs, len(request.Urls))
	for index, row := range request.Urls {
		requestData[index].CorrelationID = row.CorrelationId
		requestData[index].OriginalURL = row.OriginalUrl
	}
	resultURLs, err := s.Service.st.SaveBatch(requestData, userIDCtx)
	if err != nil {
		response.Code = http.StatusInternalServerError
		return &response, err
	}
	for index := range resultURLs {
		response.Urls = append(response.Urls,
			&pb.BatchSaveResponse_URL{CorrelationId: resultURLs[index].CorrelationID,
				ShortUrl: fmt.Sprintf("%s/%s", s.Service.host, resultURLs[index].ShortURL)})
	}
	response.Code = http.StatusAccepted
	return &response, nil
}
