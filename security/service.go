// Package security выпускает токены доступа JSON API.
package security

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: security service requires a client")
	}
	return &Service{client: client}, nil
}

type Token struct {
	AccessToken string `json:"access_token"`
}

// Токен — секрет: fmt и slog его не раскрывают. json.Marshal раскрывает: не сериализуйте Token.
func (t Token) String() string       { return "security.Token{AccessToken:<redacted>}" }
func (t Token) GoString() string     { return t.String() }
func (t Token) LogValue() slog.Value { return slog.StringValue(t.String()) }

// CreateToken выпускает токен по логину и паролю сотрудника (Basic, даже если клиент настроен на Bearer).
// Выпуск отзывает все прежние токены сотрудника. Запрос не повторяется: при потере ответа
// токен мог быть выпущен. Статус, отличный от 200, — ошибка, даже если тело содержит токен.
// Auth клиента не меняется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#4-poluchenie-novogo-tokena
func (s *Service) CreateToken(ctx context.Context, username, password string, opts ...transport.RequestOption) (*Token, error) {
	auth, err := transport.BasicAuth(username, password)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}
	var result Token
	// Basic и запрет повтора идут последними, чтобы опции caller их не заменили;
	// cap защищает массив caller от append.
	resp, err := s.client.Send(ctx, http.MethodPost, "security/token", nil, nil, &result,
		append(opts[:len(opts):len(opts)], transport.WithAuth(auth), transport.WithoutRetry())...)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: create token: unexpected HTTP status %d", resp.StatusCode)
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("moysklad: create token: response has no access_token")
	}
	return &result, nil
}
