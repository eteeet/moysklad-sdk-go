package webhookstock

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
		return nil, fmt.Errorf("moysklad: webhookstock service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-poluchit-otdelnyj-vebhuk-na-izmenenie-ostatkov
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*WebhookStock, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: webhookstock id must be a UUID")
	}
	return s.do(ctx, "get webhookstock "+id, http.MethodGet, "entity/webhookstock/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-sozdat-vebhuk-na-izmenenie-ostatkov
// Не повторять POST автоматически после неоднозначного сбоя: сочетание stockType, reportType, url
// уникально, повтор получит ошибку сервера или создаст лишний вебхук.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*WebhookStock, error) {
	transport.ResetResponseInfo(opts)
	switch {
	case input.URL == "":
		return nil, fmt.Errorf("moysklad: create webhookstock: url must not be empty")
	case input.StockType == "":
		return nil, fmt.Errorf("moysklad: create webhookstock: stockType must not be empty")
	case input.ReportType == "":
		return nil, fmt.Errorf("moysklad: create webhookstock: reportType must not be empty")
	}
	return s.do(ctx, "create webhookstock", http.MethodPost, "entity/webhookstock", input, opts)
}

// Update изменяет вебхук; UpdateRequest{Enabled: model.Value(false)} — «Отключить вебхук на изменение остатков».
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-izmenit-vebhuk-na-izmenenie-ostatkov
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*WebhookStock, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update webhookstock: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update webhookstock: no fields specified")
	}
	return s.do(ctx, "update webhookstock "+id, http.MethodPut, "entity/webhookstock/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*WebhookStock, error) {
	var result WebhookStock
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
// Параметры фильтрации и пагинации источник для этого ресурса не описывает: они уходят как заданы.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-poluchit-spisok-vebhukov-na-izmenenie-ostatkov
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[WebhookStock], error) {
	return list.Entity[WebhookStock](ctx, s.client, "webhookstock", options, nil, opts...)
}

// Delete — DELETE entity/webhookstock/{id}, безвозвратно. Повторов нет; права проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-udalit-vebhuk-na-izmenenie-ostatkov
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete webhookstock: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete webhookstock "+id, "entity/webhookstock/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/webhookstock (до batch.MaxItems); обязательные поля
// проверяет сервер у каждого элемента, ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-massovoe-sozdanie-i-obnovlenie-vebhukov-na-izmenenie-ostatkov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[WebhookStock], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[WebhookStock](ctx, s.client, "webhookstock", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/webhookstock: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-massovoe-sozdanie-i-obnovlenie-vebhukov-na-izmenenie-ostatkov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[WebhookStock], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[WebhookStock](ctx, s.client, "webhookstock", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/webhookstock/delete, безвозвратно. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-massovoe-udalenie-vebhukov-na-izmenenie-ostatkov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "webhookstock", ids, opts)
}
