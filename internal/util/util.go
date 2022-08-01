package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const hashLen = 4

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
func ShortenLink(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:hashLen])
}

func RequestBodyCheck(w http.ResponseWriter, r *http.Request) error {
	if r.Body == http.NoBody {
		err := errors.New("request should have body")
		http.Error(w, string(err.Error()), http.StatusBadRequest)
		return err
	}
	return nil
}

func StuctEncode(s interface{}) string {
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(s)
	if err != nil {
		fmt.Println(err)
	}
	return buf.String()
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	JSONResponse(w, struct {
		ErrorMessage interface{} `json:"api_error"`
	}{
		ErrorMessage: err,
	}, code)
}

func JSONResponse(w http.ResponseWriter, message interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(message)
	if err != nil {
		log.Println(err)
	}
}

func GetCookieParam(paramName string, r *http.Request) (paramValue string) {
	cookieParam, err := r.Cookie(paramName)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Cookie '%s':%s\n", paramName, cookieParam.Value)
	return cookieParam.Value
}

func GetSign(msg []byte, secretKey string) (stringSign string) {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(msg)
	sign := h.Sum(nil)
	stringSign = hex.EncodeToString(sign)
	return
}
