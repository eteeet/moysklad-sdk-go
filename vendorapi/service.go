// Package vendorapi — Vendor API 1.0 серверного решения: запросы к МоемуСкладу и проверка входящих.
// Состояние (jti, токены, статусы) хранит решение, не SDK.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/moysklad-endpoints#2-rest-endpointy-na-storone-moegosklada
package vendorapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/contexts"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// DefaultBaseURL — BASE-ENDPOINT Vendor API; задаётся в transport.Config.BaseURL.
const DefaultBaseURL = "https://apps-api.moysklad.ru/api/vendor/1.0"

// Service обращается к Vendor API через client с BaseURL Vendor API и Auth из AppAuth.
type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: vendor service requires a client")
	}
	return &Service{client: client}, nil
}

// Context возвращает сотрудника iframe или виджета в форме context/employee JSON API.
// Один POST без тела: RetryPolicy POST не повторяет.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/moysklad-endpoints#3-poluchenie-konteksta-polzovatelya
func (s *Service) Context(ctx context.Context, contextKey string, opts ...transport.RequestOption) (*contexts.EmployeeContext, error) {
	transport.ResetResponseInfo(opts)
	// "." и ".." JoinPath схлопнул бы в context или корень API; "user" — маршрут обмена токена.
	if contextKey == "" || contextKey == "." || contextKey == ".." || contextKey == "user" || strings.ContainsAny(contextKey, "/?#") {
		return nil, fmt.Errorf("moysklad: get vendor context: invalid context key")
	}
	var result contexts.EmployeeContext
	resp, err := s.client.Send(ctx, http.MethodPost, "context/"+contextKey, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("get vendor context: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: get vendor context: unexpected HTTP status %d", resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: get vendor context: response has no id")
	}
	return &result, nil
}
