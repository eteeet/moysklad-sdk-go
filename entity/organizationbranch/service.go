package organizationbranch

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
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
		return nil, fmt.Errorf("moysklad: organizationbranch service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-poluchit-podrazdelenie-yuridicheskogo-lica
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*OrganizationBranch, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: organizationbranch id must be a UUID")
	}
	return s.do(ctx, "get organizationbranch "+id, http.MethodGet, "entity/organizationbranch/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-sozdat-podrazdelenie-yuridicheskogo-lica
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*OrganizationBranch, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create organizationbranch: name must not be empty")
	}
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create organizationbranch: organization meta.href must not be empty")
	}
	return s.do(ctx, "create organizationbranch", http.MethodPost, "entity/organizationbranch", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-izmenit-podrazdelenie-yuridicheskogo-lica
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*OrganizationBranch, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update organizationbranch: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update organizationbranch: no fields specified")
	}
	return s.do(ctx, "update organizationbranch "+id, http.MethodPut, "entity/organizationbranch/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*OrganizationBranch, error) {
	var result OrganizationBranch
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-poluchit-spisok-podrazdelenij-yuridicheskogo-lica
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[OrganizationBranch], error) {
	return list.Entity[OrganizationBranch](ctx, s.client, "organizationbranch", options, nil, opts...)
}

// Delete — DELETE entity/organizationbranch/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-udalit-podrazdelenie-yuridicheskogo-lica
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete organizationbranch: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete organizationbranch "+id, "entity/organizationbranch/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/organizationbranch (до batch.MaxItems); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-massovoe-sozdanie-i-obnovlenie-podrazdelenij-yuridicheskogo-lica
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[OrganizationBranch], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[OrganizationBranch](ctx, s.client, "organizationbranch", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/organizationbranch: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-massovoe-sozdanie-i-obnovlenie-podrazdelenij-yuridicheskogo-lica
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[OrganizationBranch], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[OrganizationBranch](ctx, s.client, "organizationbranch", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/organizationbranch/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-massovoe-udalenie-podrazdelenij-yuridicheskogo-lica
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "organizationbranch", ids, opts)
}
