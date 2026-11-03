package processingplan

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Продукты — entity/processingplan/{id}/products. Повторов нет.

func (s *Service) products(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processingplan", ParentID: id, Path: "products", Item: "product", MetaType: "processingplanresult"}
}

// GetProducts читает одну страницу продуктов: только limit и offset, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-produkty-tehkarty
func (s *Service) GetProducts(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Product], error) {
	return nested.List[Product](ctx, s.products(id), page, opts)
}

// GetProduct: пример ответа источника — массив из одного объекта (у остальных наборов — объект);
// принимаются обе формы, массив другого размера — ошибка.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-produkt
func (s *Service) GetProduct(ctx context.Context, id, productID string, opts ...transport.RequestOption) (*Product, error) {
	transport.ResetResponseInfo(opts)
	operation := "get processingplan " + id + " product " + productID
	if !validID(id) || !validID(productID) {
		return nil, fmt.Errorf("moysklad: %s: processingplan and product ids must be UUIDs", operation)
	}
	var raw jsontext.Value
	resp, err := s.client.Send(ctx, http.MethodGet, "entity/processingplan/"+id+"/products/"+productID, nil, nil, &raw, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if raw.Kind() == '[' {
		var items []jsontext.Value
		if err := json.Unmarshal(raw, &items); err != nil || len(items) != 1 {
			return nil, fmt.Errorf("moysklad: %s: response array must contain exactly one object", operation)
		}
		raw = items[0]
	}
	var result Product
	if err := json.Unmarshal(raw, &result); err != nil || result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// CreateProducts добавляет продукты к существующим одним POST массивом; ответ — по правилам batch.
// В параметрической техкарте продукт один — добавление ограничивает сервер. Не повторять после
// неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-sozdat-produkt
func (s *Service) CreateProducts(ctx context.Context, id string, items []ProductRequest, opts ...transport.RequestOption) ([]batch.Result[Product], error) {
	return nested.CreateMany[Product](ctx, s.products(id), items, opts)
}

// UpdateProduct: обязательных полей нет; пустой input отклоняется до запроса. Пример ответа
// источника — массив из одного объекта, принимается и объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-izmenit-produkt
func (s *Service) UpdateProduct(ctx context.Context, id, productID string, input ProductRequest, opts ...transport.RequestOption) (*Product, error) {
	return nested.Update[Product](ctx, s.products(id), productID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-udalit-produkt
func (s *Service) DeleteProduct(ctx context.Context, id, productID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.products(id), productID, opts)
}

// DeleteManyProducts — POST …/products/delete с meta продуктов (1–1000 id); ответ без тела,
// результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-massovoe-udalenie-produkta
func (s *Service) DeleteManyProducts(ctx context.Context, id string, productIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.products(id), "delete processingplan "+id+" products", productIDs, opts)
}
