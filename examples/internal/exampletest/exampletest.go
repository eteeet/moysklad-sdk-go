// Package exampletest поднимает httptest-сервер для тестов примеров.
package exampletest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
)

// Client — клиент примеров с токеном "token" на httptest-сервер с handler h. Адрес сервера
// для href в ответах — "http://" + r.Host; сервер закрывается в t.Cleanup.
func Client(t *testing.T, h http.HandlerFunc) *moysklad.Client {
	t.Helper()
	server := httptest.NewServer(h)
	t.Cleanup(server.Close)
	c, err := env.ClientFor(server.URL+"/api/remap/1.2", "token")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Unauthorized отвечает 401 с ошибкой API 1056.
func Unauthorized(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"errors":[{"error":"auth","code":1056}]}`))
}
