package processingplan

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
		return nil, fmt.Errorf("moysklad: processingplan service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-tehkartu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*ProcessingPlan, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get processingplan: id must be a UUID")
	}
	return s.do(ctx, "get processingplan "+id, http.MethodGet, "entity/processingplan/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль техкарты.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-sozdat-tehkartu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*ProcessingPlan, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create processingplan"
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	if len(input.Products) == 0 {
		return nil, fmt.Errorf("moysklad: %s: products must not be empty", operation)
	}
	if err := checkRows(model.Value(input.Products), input.Stages, input.Materials, input.ParametricMaterials); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/processingplan", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-izmenit-tehkartu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*ProcessingPlan, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update processingplan"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	if err := checkRows(input.Products, input.Stages, input.Materials, input.ParametricMaterials); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/processingplan/"+id, input, opts)
}

// Лимит строк, сохраняемых вместе с техкартой, — 1000 на каждый набор.
func checkRows(products model.Field[[]ProductRequest], stages model.Field[[]StageRequest], materials model.Field[[]MaterialRequest], parametric model.Field[[]ParametricMaterialRequest]) error {
	for i, n := range []int{count(products), count(stages), count(materials), count(parametric)} {
		if n > 1000 {
			return fmt.Errorf("%s must not exceed 1000", []string{"products", "stages", "materials", "parametricMaterials"}[i])
		}
	}
	return nil
}

func count[T any](rows model.Field[[]T]) int {
	v, _ := rows.Get()
	return len(v)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-spisok-tehkart
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[ProcessingPlan], error) {
	return list.Entity[ProcessingPlan](ctx, s.client, "processingplan", options, nil, opts...)
}

// Delete — DELETE entity/processingplan/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-udalit-tehkartu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete processingplan: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete processingplan "+id, "entity/processingplan/"+id, opts)
}

// CreateMany создаёт техкарты одним POST entity/processingplan (до batch.MaxItems). Обязательные
// поля проверяет сервер у каждого элемента; ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-massovoe-sozdanie-i-obnovlenie-tehkart
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[ProcessingPlan], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[ProcessingPlan](ctx, s.client, "processingplan", items, opts)
}

// UpdateMany изменяет техкарты одним POST entity/processingplan: meta по ID и поля Request.
// Ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-massovoe-sozdanie-i-obnovlenie-tehkart
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[ProcessingPlan], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[ProcessingPlan](ctx, s.client, "processingplan", items, opts)
}

// DeleteMany удаляет техкарты по ids одним POST entity/processingplan/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-massovoe-udalenie-tehkart
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "processingplan", ids, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*ProcessingPlan, error) {
	var result ProcessingPlan
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
