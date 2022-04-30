package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func GetMD5Hash(text string, len int) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:len])
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
	encoder.Encode(s)
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
	encoder.Encode(message)
}
