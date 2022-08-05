package cookie

import (
	"context"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type ContextType string

const (
	CookieUserIDField                = "user_id"
	CookieMaxAge                     = 864000
	uuidStringLength                 = 36
	signSentencePosition             = 37
	UserIDCtxName        ContextType = "ctxUserId"
)

type CookieHandler struct {
	SecretKey string
	L         *logger.Logger
}

func (h *CookieHandler) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieUserID := util.GetCookieParam(CookieUserIDField, r)
		h.L.Debug(cookieUserID)
		if cookieUserID != "" {
			gotUUID := uuid.FromStringOrNil(cookieUserID[:uuidStringLength])
			requiredSign := util.GetSign(gotUUID.Bytes(), h.SecretKey)
			checkSign := cookieUserID[signSentencePosition:] == requiredSign
			h.L.Debug("Sign check status:", checkSign)
			if checkSign {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, gotUUID.String())))
				return
			}
		}
		userID := uuid.NewV4()
		h.L.Debug("User id: %s\n", userID)
		stringSign := util.GetSign(userID.Bytes(), h.SecretKey)
		cookieUserID = userID.String()
		cookie := &http.Cookie{
			Name:   CookieUserIDField,
			Value:  cookieUserID + "|" + stringSign,
			MaxAge: CookieMaxAge,
			Path:   "/",
		}
		http.SetCookie(w, cookie)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, cookieUserID)))
	})
}
