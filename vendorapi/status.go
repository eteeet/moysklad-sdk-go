package vendorapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

// AppStatus — статус установки решения на аккаунте. Access заполнен только в Activated и
// SettingsRequired, токена в нём нет.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/moysklad-endpoints#3-poluchenie-statusa-resheniya-na-akkaunte
type AppStatus struct {
	Status       Status        `json:"status"`
	Cause        Cause         `json:"cause,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`
	Access       []Access      `json:"access,omitempty"`
}

// Status читает статус установки: один GET, RetryPolicy его повторяет. 404 (2004) — не установлено.
func (s *Service) Status(ctx context.Context, appID, accountID string, opts ...transport.RequestOption) (*AppStatus, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get vendor app status"
	if !validID(appID) || !validID(accountID) {
		return nil, fmt.Errorf("moysklad: %s: appId and accountId must be UUIDs", operation)
	}
	var result AppStatus
	resp, err := s.client.Send(ctx, http.MethodGet, statusPath(appID, accountID), nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Status == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no status", operation)
	}
	return &result, nil
}

// SetStatus сообщает МоемуСкладу итог активации после ответа Activating или SettingsRequired.
// Решение выставляет только Activating, SettingsRequired, Activated; 409 — перехода нет.
// https://dev.moysklad.ru/doc/api/vendor/1.0/#/vendor-api/moysklad-endpoints#3-izmenenie-statusa-resheniya-na-akkaunte
func (s *Service) SetStatus(ctx context.Context, appID, accountID string, status Status, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	const operation = "set vendor app status"
	if !validID(appID) || !validID(accountID) {
		return fmt.Errorf("moysklad: %s: appId and accountId must be UUIDs", operation)
	}
	if status != StatusActivating && status != StatusSettingsRequired && status != StatusActivated {
		return fmt.Errorf("moysklad: %s: status %q is set only by MoySklad", operation, status)
	}
	body := struct {
		Status Status `json:"status"`
	}{status}
	resp, err := s.client.Send(ctx, http.MethodPut, statusPath(appID, accountID), nil, body, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}

func statusPath(appID, accountID string) string {
	return "apps/" + appID + "/" + accountID + "/status"
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
