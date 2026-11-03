// Package single — одиночный ресурс без id (настройки справочника, компании, пользователя):
// один GET или PUT, ответ — объект с непустым meta.href.
package single

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Get читает ресурс endpoint.
func Get[T any](ctx context.Context, client *transport.Client, operation, endpoint string, opts []transport.RequestOption) (*T, error) {
	transport.ResetResponseInfo(opts)
	return send[T](ctx, client, operation, http.MethodGet, endpoint, nil, opts)
}

// Put изменяет ресурс частично: передаются только заданные поля, ответ — ресурс целиком. Пустой
// body (все Field нулевые) — ошибка до запроса. PUT не повторяется.
func Put[T any](ctx context.Context, client *transport.Client, operation, endpoint string, body any, opts []transport.RequestOption) (*T, error) {
	transport.ResetResponseInfo(opts)
	if reflect.ValueOf(body).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return send[T](ctx, client, operation, http.MethodPut, endpoint, body, opts)
}

func send[T any](ctx context.Context, client *transport.Client, operation, method, endpoint string, body any, opts []transport.RequestOption) (*T, error) {
	var raw jsontext.Value
	resp, err := client.Send(ctx, method, endpoint, nil, body, &raw, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	var head struct {
		Meta struct {
			Href string `json:"href"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(raw, &head); err != nil || head.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("moysklad: %s: decode response: %w", operation, err)
	}
	return &result, nil
}
