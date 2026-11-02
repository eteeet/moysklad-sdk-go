package variant

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-poluchit-modifikaciyu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Variant, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: variant id must be a UUID")
	}
	return s.do(ctx, "get variant "+id, http.MethodGet, "entity/variant/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-sozdat-modifikaciyu
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Variant.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Variant, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create variant"
	if input.Product.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: product meta.href must not be empty", operation)
	}
	if len(input.Characteristics) == 0 {
		return nil, fmt.Errorf("moysklad: %s: characteristics must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/variant", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-izmenit-modifikaciyu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Variant, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update variant"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/variant/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Поиск — только по name; товар отбирается фильтром productid (= и !=), штрихкод — barcode
// в единственном числе. fields=minimumStock не принимается: остаток отложен, как у Product.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-poluchit-spisok-modifikacij
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Variant], error) {
	return list.Entity[Variant](ctx, s.client, "variant", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Variant, error) {
	var result Variant
	if err := s.send(ctx, operation, method, endpoint, body, &result, opts); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// Delete — «Удалить Модификацию»: DELETE entity/variant/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-udalit-modifikaciyu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete variant: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete variant "+id, "entity/variant/"+id, opts)
}
