package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("authorization header does not exist")
	}
	split := strings.Fields(authorization)
	if len(split) != 2 || split[0] != "ApiKey" {
		return "", fmt.Errorf("invalid authorization header, valid syntax ApiKey <token>")
	}
	return split[1], nil
}
