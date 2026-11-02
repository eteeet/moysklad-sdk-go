// Package env собирает клиент примеров из переменных окружения.
package env

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Client читает MOYSKLAD_TOKEN: персональный токен сотрудника и токен решения передаются одинаково,
// Bearer. Выпуск токена (security/token) в примеры не входит. MOYSKLAD_BASE_URL — другой адрес API.
func Client() (*moysklad.Client, error) {
	token, err := Require("MOYSKLAD_TOKEN")
	if err != nil {
		return nil, err
	}
	return ClientFor(os.Getenv("MOYSKLAD_BASE_URL"), token)
}

// ClientFor — клиент с Bearer token; пустой baseURL — адрес API по умолчанию.
func ClientFor(baseURL, token string) (*moysklad.Client, error) {
	auth, err := transport.BearerToken(token)
	if err != nil {
		return nil, err
	}
	return moysklad.NewClient(moysklad.Config{BaseURL: baseURL, Auth: auth})
}

// Require возвращает непустое значение переменной окружения или ошибку с её именем.
func Require(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

// Context отменяется по Ctrl+C и через timeout.
func Context(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	return ctx, func() { cancel(); stop() }
}

// Main печатает ошибку run и завершает процесс с кодом 1.
func Main(run func() error) {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
