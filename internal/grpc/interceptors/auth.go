package interceptors

import (
	"context"
	"errors"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Interceptor struct {
	l         *logger.Logger
	secretKey string
}

func SetInterceptorLogger(l *logger.Logger) func(*Interceptor) {
	return func(i *Interceptor) {
		i.l = l
	}
}
func SetSecretKey(secretKey string) func(*Interceptor) {
	return func(i *Interceptor) {
		i.secretKey = secretKey
	}
}

func (i *Interceptor) AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// https://stackoverflow.com/questions/71114401/grpc-how-to-pass-value-from-interceptor-to-service-function
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("could not load metadata")
	}
	if len(md[cookie.CookieUserIDField]) > 0 {
		userIDfromMD := md[cookie.CookieUserIDField][0]
		gotUUID := uuid.FromStringOrNil(userIDfromMD[:cookie.UIDStringLength])
		requiredSign := util.GetSign(gotUUID.Bytes(), i.secretKey)
		checkSign := userIDfromMD[cookie.SignSentencePosition:] == requiredSign
		if checkSign {
			md.Append(string(cookie.UserIDCtxName), gotUUID.String())
			md.Append(cookie.CookieUserIDField, requiredSign)
			return handler(metadata.NewIncomingContext(ctx, md), req)
		}
	}
	userUUID := uuid.NewV4()
	i.l.Info("User uuid is " + userUUID.String())
	stringSign := util.GetSign(userUUID.Bytes(), i.secretKey)
	i.l.Info("Signed uuid is  " + stringSign)
	md.Append(string(cookie.UserIDCtxName), userUUID.String())
	md.Append(cookie.CookieUserIDField, stringSign)
	return handler(metadata.NewIncomingContext(ctx, md), req)
}
