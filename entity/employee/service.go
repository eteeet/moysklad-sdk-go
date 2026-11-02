package employee

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: employee service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-poluchit-sotrudnika
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Employee, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: employee id must be a UUID")
	}
	return s.do(ctx, "get employee "+id, http.MethodGet, "entity/employee/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-sozdat-sotrudnika
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Employee.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Employee, error) {
	transport.ResetResponseInfo(opts)
	if input.LastName == "" {
		return nil, fmt.Errorf("moysklad: create employee: lastName must not be empty")
	}
	return s.do(ctx, "create employee", http.MethodPost, "entity/employee", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-izmenit-sotrudnika
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Employee, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update employee: id must be a UUID")
	}
	if input.LastName == "" {
		return nil, fmt.Errorf("moysklad: update employee: lastName must not be empty")
	}
	return s.do(ctx, "update employee "+id, http.MethodPut, "entity/employee/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Employee, error) {
	var result Employee
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, &result, opts...)
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

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Employee], error) {
	return list.Entity[Employee](ctx, s.client, "employee", options, nil, opts...)
}

// Delete — «Удалить Сотрудника»: DELETE entity/employee/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-udalit-sotrudnika
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete employee: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete employee "+id, "entity/employee/"+id, opts)
}
