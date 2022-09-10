package server

import (
	"context"
	"net/http"

	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
)

// Ping implement method for gRPC for check DB connection
func (s *GRPCServer) Ping(ctx context.Context, request *pb.EmptyRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse
	message := "connection established"
	status := http.StatusOK
	var err error
	if err = s.Service.st.CheckConnection(ctx); err != nil {
		message = err.Error()
		status = http.StatusInternalServerError
	}
	response.Message = message
	response.Code = int32(status)

	return &response, err
}
