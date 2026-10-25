package transport

// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-autentifikaciya

import (
	"fmt"
	"net/http"
	"strings"
)

// nil Auth оставляет авторизацию пользовательскому HTTP-клиенту.
type Auth func(*http.Request) error

func BearerToken(token string) (Auth, error) {
	if token == "" || strings.ContainsAny(token, " \t\r\n") {
		return nil, fmt.Errorf("moysklad: invalid bearer token")
	}

	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}, nil
}

func BasicAuth(username, password string) (Auth, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("moysklad: username and password are required")
	}
	if strings.Contains(username, ":") {
		return nil, fmt.Errorf("moysklad: username must not contain a colon")
	}

	return func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	}, nil
}
