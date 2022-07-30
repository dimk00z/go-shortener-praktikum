package cookie

import (
	"context"
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	uuid "github.com/satori/go.uuid"
)

type ContextType string

const (
	cookieUserIDField                = "user_id"
	cookieMaxAge                     = 864000
	uuidStringLength                 = 36
	signSentencePosition             = 37
	UserIDCtxName        ContextType = "ctxUserId"
)

type CookieHandler struct {
	SecretKey string
}

func (h *CookieHandler) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieUserID := util.GetCookieParam(cookieUserIDField, r)
		log.Println(cookieUserID)
		if cookieUserID != "" {
			gotUUID := uuid.FromStringOrNil(cookieUserID[:uuidStringLength])
			requiredSign := util.GetSign(gotUUID.Bytes(), h.SecretKey)
			checkSign := cookieUserID[signSentencePosition:] == requiredSign
			log.Println("Sign check status:", checkSign)
			if checkSign {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, gotUUID.String())))
				return
			}
		}
		userID := uuid.NewV4()
		log.Printf("User id: %s\n", userID)
		stringSign := util.GetSign(userID.Bytes(), h.SecretKey)
		cookieUserID = userID.String()
		cookie := &http.Cookie{
			Name:   cookieUserIDField,
			Value:  cookieUserID + "|" + stringSign,
			MaxAge: cookieMaxAge,
			Path:   "/",
		}
		http.SetCookie(w, cookie)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIDCtxName, cookieUserID)))
	})
}
