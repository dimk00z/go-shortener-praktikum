package server

type GRPCServer struct {
	port string
}

func NewGRPCServer(port string) *GRPCServer {
	return &GRPCServer{
		port: port,
	}
}
