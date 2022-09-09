package server

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/dimk00z/go-shortener-praktikum/internal/grpc/interceptors"
	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"google.golang.org/grpc/metadata"
)

type GRPCServer struct {
	pb.UnimplementedShortenerServer
	Service *Service
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{Service: &Service{Interceptor: &interceptors.Interceptor{}}}
}
func getUserIDFromMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if len(md[strings.ToLower(string(cookie.UserIDCtxName))]) > 0 {
		return md[strings.ToLower(string(cookie.UserIDCtxName))][0]
	}
	return ""
}

// Ping implement method for gRPC for check DB connection
func (s *GRPCServer) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
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

	// s.Service.l.Debug("Ping userID " + getUserIDFromMD(ctx))
	return &response, err
}

// Get IP from GRPC context
func GetIP(ctx context.Context) string {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("x-real-ip")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIp := ips[0]
				return clientIp
			}
		}
	}
	return ""
}

func (s *GRPCServer) GetStats(ctx context.Context, request *pb.StatsRequest) (*pb.StatsResponse, error) {
	var response pb.StatsResponse
	var err error
	clientRealIP := GetIP(ctx)
	response.Code = http.StatusForbidden
	if GetIP(ctx) == "" {
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
