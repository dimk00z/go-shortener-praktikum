package server

import (
	"context"
	"strings"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"google.golang.org/grpc/metadata"
)

// Get IP from GRPC context
func getIP(ctx context.Context) string {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("x-real-ip")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIP := ips[0]
				return clientIP
			}
		}
	}
	return ""
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
