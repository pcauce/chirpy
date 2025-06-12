package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "ApiKey ") || len(token) <= 7 {
		return "", errors.New("invalid key format. Should be 'ApiKey <key>'")
	}
	return strings.TrimPrefix(token, "ApiKey "), nil
}
