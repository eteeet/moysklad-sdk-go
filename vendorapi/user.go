package vendorapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

// UserContext — сотрудник, для которого хост-окно выдало одноразовый токен.
type UserContext struct {
	AccountID string `json:"accountId"`
	UserID    string `json:"userId"`
	UserUID   string `json:"userUid"`
	Role      Role   `json:"role"`
}

// Role — роль сотрудника, как в обработке кастомной кнопки. Неизвестные значения сохраняются.
type Role string

const (
	RoleAdmin      Role = "admin"
	RoleCashier    Role = "cashier"
	RoleWorker     Role = "worker"
	RoleIndividual Role = "individual"
)

// UserContext обменивает одноразовый токен хост-окна (60 с) на сотрудника. RetryPolicy POST не
// повторяет: после потерянного ответа токен уже израсходован. 404 (3007) — использован или истёк.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/moysklad-endpoints#3-poluchenie-konteksta-polzovatelya-po-odnorazovomu-tokenu
func (s *Service) UserContext(ctx context.Context, token string, opts ...transport.RequestOption) (*UserContext, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get vendor user context"
	if token == "" {
		return nil, fmt.Errorf("moysklad: %s: token is required", operation)
	}
	body := struct {
		Token string `json:"token"`
	}{token}
	var result UserContext
	resp, err := s.client.Send(ctx, http.MethodPost, "context/user", nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	// Все четыре поля источник помечает обязательными.
	if result.AccountID == "" || result.UserID == "" || result.UserUID == "" || result.Role == "" {
		return nil, fmt.Errorf("moysklad: %s: response lacks a required field", operation)
	}
	return &result, nil
}
