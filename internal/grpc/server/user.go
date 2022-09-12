package server

import (
	"context"
	"net/http"

	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
)

var ()

func (s *GRPCServer) GetUsersURLs(ctx context.Context, request *pb.EmptyRequest) (*pb.UserURLsResponse, error) {
	var response pb.UserURLsResponse
	response.Code = http.StatusOK
	userIDCtx := getUserIDFromMD(ctx)
	var err error
	userURLs, err := s.Service.st.GetUserURLs(userIDCtx)
	if err != nil {
		response.Code = http.StatusNoContent
		return &response, err

	}
	for _, userURL := range userURLs {
		response.Urls = append(response.Urls, &pb.UserURLsResponse_URL{
			ShortUrl:    s.Service.host + "/" + userURL.ShortURL,
			OriginalUrl: userURL.URL,
		})
	}
	return &response, err
}

func (s *GRPCServer) DelBatch(ctx context.Context, request *pb.DelBatchRequest) (*pb.DelBatchResponse, error) {
	var response pb.DelBatchResponse
	var err error
	if len(request.Urls) == 0 {
		response.Code = http.StatusBadRequest
		err = errEmptyBatchGiven
		return &response, err
	}
	response.Code = http.StatusAccepted
	userIDCtx := getUserIDFromMD(ctx)
	var shortURLs models.BatchForDelete = request.Urls
	deleteBatchTask := func(ctx context.Context) error {
		return s.Service.st.DeleteBatch(ctx, shortURLs, userIDCtx)
	}
	s.Service.wp.Push(deleteBatchTask)
	return &response, err
}
