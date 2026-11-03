package processingplanfolder

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

// Service — группы техкарт. Метаданные (GET entity/processingplanfolder/metadata) читает
// metadata.Service.Get(ctx, "processingplanfolder").
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-metadannye-grupp-tehkart
type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: processingplanfolder service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-poluchit-gruppu-tehkart
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*ProcessingPlanFolder, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: processingplanfolder id must be a UUID")
	}
	return s.do(ctx, "get processingplanfolder "+id, http.MethodGet, "entity/processingplanfolder/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-sozdat-novuyu-gruppu-tehkart
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*ProcessingPlanFolder, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create processingplanfolder: name must not be empty")
	}
	return s.do(ctx, "create processingplanfolder", http.MethodPost, "entity/processingplanfolder", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-izmenit-gruppu-tehkart
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*ProcessingPlanFolder, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update processingplanfolder: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update processingplanfolder: no fields specified")
	}
	return s.do(ctx, "update processingplanfolder "+id, http.MethodPut, "entity/processingplanfolder/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*ProcessingPlanFolder, error) {
	var result ProcessingPlanFolder
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-poluchit-spisok-grupp-tehkart
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[ProcessingPlanFolder], error) {
	return list.Entity[ProcessingPlanFolder](ctx, s.client, "processingplanfolder", options, nil, opts...)
}

// Delete — DELETE entity/processingplanfolder/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-udalit-gruppu-tehkart
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete processingplanfolder: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete processingplanfolder "+id, "entity/processingplanfolder/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/processingplanfolder (до batch.MaxItems); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-massovoe-sozdanie-i-obnovlenie-grupp-tehkart
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[ProcessingPlanFolder], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[ProcessingPlanFolder](ctx, s.client, "processingplanfolder", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/processingplanfolder: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-massovoe-sozdanie-i-obnovlenie-grupp-tehkart
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[ProcessingPlanFolder], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[ProcessingPlanFolder](ctx, s.client, "processingplanfolder", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/processingplanfolder/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-massovoe-udalenie-grupp-tehkart
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "processingplanfolder", ids, opts)
}
