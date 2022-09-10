package server

import (
	"context"
	"net/http"

	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
)

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
	// TODO:add logic
	return &response, err
}
