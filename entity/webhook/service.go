package webhook

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
		return nil, fmt.Errorf("moysklad: webhook service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-poluchit-otdelnyj-vebhuk
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Webhook, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: webhook id must be a UUID")
	}
	return s.do(ctx, "get webhook "+id, http.MethodGet, "entity/webhook/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-sozdat-vebhuk
// Не повторять POST автоматически после неоднозначного сбоя: сочетание entityType, action, url
// уникально, повтор получит ошибку сервера или создаст лишний вебхук.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Webhook, error) {
	transport.ResetResponseInfo(opts)
	switch {
	case input.URL == "":
		return nil, fmt.Errorf("moysklad: create webhook: url must not be empty")
	case input.Action == "":
		return nil, fmt.Errorf("moysklad: create webhook: action must not be empty")
	case input.EntityType == "":
		return nil, fmt.Errorf("moysklad: create webhook: entityType must not be empty")
	}
	return s.do(ctx, "create webhook", http.MethodPost, "entity/webhook", input, opts)
}

// Update изменяет вебхук; UpdateRequest{Enabled: model.Value(false)} — «Отключить вебхук».
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-izmenit-vebhuk
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Webhook, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update webhook: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update webhook: no fields specified")
	}
	return s.do(ctx, "update webhook "+id, http.MethodPut, "entity/webhook/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Webhook, error) {
	var result Webhook
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
// Параметры фильтрации и пагинации источник для вебхуков не описывает: они уходят как заданы.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-poluchit-spisok-vebhukov
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Webhook], error) {
	return list.Entity[Webhook](ctx, s.client, "webhook", options, nil, opts...)
}

// Delete — DELETE entity/webhook/{id}, безвозвратно. Повторов нет; права проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-udalit-vebhuk
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete webhook: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete webhook "+id, "entity/webhook/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/webhook (до batch.MaxItems); обязательные поля
// проверяет сервер у каждого элемента, ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-massovoe-sozdanie-i-obnovlenie-vebhukov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Webhook], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Webhook](ctx, s.client, "webhook", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/webhook: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-massovoe-sozdanie-i-obnovlenie-vebhukov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Webhook], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Webhook](ctx, s.client, "webhook", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/webhook/delete, безвозвратно. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-massovoe-udalenie-vebhukov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "webhook", ids, opts)
}
