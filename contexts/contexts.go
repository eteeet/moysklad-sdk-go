// Package contexts читает контекст запроса: сотрудника, от имени которого авторизован клиент.
// Это не model.ResponseContext из коллекций.
package contexts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: contexts service requires a client")
	}
	return &Service{client: client}, nil
}

// EmployeeContext — сотрудник, от имени которого выполняется запрос, и его права.
// Все поля только для чтения. Поля сотрудника — из employee.Employee; Salary в таблице
// контекста нет и остаётся nil. Права — в общем формате model.Permissions, без потерь.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-kontekst-zaprosa-sotrudnika
type EmployeeContext struct {
	employee.Employee
	Code        string            `json:"code"`
	Permissions model.Permissions `json:"permissions"`
}

// Employee выполняет один GET context/employee без ID: сервер определяет сотрудника по авторизации.
// SDK не выводит права из роли и не делает других запросов.
func (s *Service) Employee(ctx context.Context, opts ...transport.RequestOption) (*EmployeeContext, error) {
	transport.ResetResponseInfo(opts)
	return get(ctx, s, "get employee context", "context/employee", func(c *EmployeeContext) string { return c.ID }, opts)
}

// get выполняет один GET контекста. Критерий успеха контекста — HTTP 200 и id, который
// источник помечает обязательным в ответе; это не общий ID-gate сущностей.
func get[T any](ctx context.Context, s *Service, operation, endpoint string, id func(*T) string, opts []transport.RequestOption) (*T, error) {
	var result T
	resp, err := s.client.Send(ctx, http.MethodGet, endpoint, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if id(&result) == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}
