package server

import (
	"context"
	"net"
	"net/http"

	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
)

func (s *GRPCServer) GetStats(ctx context.Context, request *pb.EmptyRequest) (*pb.StatsResponse, error) {
	var response pb.StatsResponse
	var err error
	clientRealIP := getIP(ctx)
	response.Code = http.StatusForbidden
	if getIP(ctx) == "" {
		return &response, err
	}
	_, ipnet, err := net.ParseCIDR(s.Service.trustedSubnet)
	if err != nil {
		response.Code = http.StatusInternalServerError
		return &response, err
	}
	if !ipnet.Contains(net.ParseIP(clientRealIP)) {
		return &response, err

	}
	response.Code = http.StatusOK
	stat, err := s.Service.st.GetStat()
	response.Urls = int32(stat.URLs)
	response.Users = int32(stat.Users)
	return &response, err
}
