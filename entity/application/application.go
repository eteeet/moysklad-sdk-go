// Package application читает установленное на аккаунте серверное решение по id:
// GET entity/application/{id}. Контекст решения текущего токена — contexts.Service.Application.
// Список, установку и изменение решений источник для JSON API не описывает; это не Vendor API.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#4-poluchenie-sushnosti-ustanovlennogo-resheniya
package application

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/contexts"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Application — установленное решение: meta, id установки и appUid
// (<Application-Alias>.<Vendor-Alias>@<Account-Name>, строка). Форма ответа совпадает с
// context/application, поэтому тип тот же. Все поля только для чтения.
type Application = contexts.ApplicationContext

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: application service requires a client")
	}
	return &Service{client: client}, nil
}

// Get читает решение по id установки на аккаунте. Один GET; доступ проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#4-poluchenie-sushnosti-ustanovlennogo-resheniya
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Application, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get application: id %q must be a UUID", id)
	}
	operation := "get application " + id
	var result Application
	resp, err := s.client.Send(ctx, http.MethodGet, "entity/application/"+id, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
