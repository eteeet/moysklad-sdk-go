package processing

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Два набора строк одной формы: материалы — entity/processing/{id}/materials, продукты —
// …/{id}/products. Отдельные ресурсы нужны для документов больше 1000 строк. Повторов нет.

func (s *Service) materials(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processing", ParentID: id, Path: "materials", Item: "material", MetaType: "processingpositionmaterial"}
}

func (s *Service) products(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processing", ParentID: id, Path: "products", Item: "product", MetaType: "processingpositionresult"}
}

// page — одна страница набора: только limit и offset, без скрытого Get и автообхода.
func page(ctx context.Context, c nested.Collection, p positions.PageOptions, opts []transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	operation := "get processing " + c.ParentID + " " + c.Path
	if !validID(c.ParentID) {
		return nil, fmt.Errorf("moysklad: %s: document id must be a UUID", operation)
	}
	if p.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, p.Fields)
	}
	params, err := p.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Position](ctx, c.Client, operation, "entity/processing/"+c.ParentID+"/"+c.Path, params, opts...)
}

// create — POST одной строки объектом (как в примере источника); ответ — объект или массив из
// одного объекта. assortment и quantity источник называет обязательными.
func create(ctx context.Context, c nested.Collection, input PositionRequest, opts []transport.RequestOption) (*Position, error) {
	transport.ResetResponseInfo(opts)
	if ref, ok := input.Assortment.Get(); !ok || ref.Meta.Href == "" || input.Quantity.IsZero() || input.Quantity.IsNull() {
		return nil, fmt.Errorf("moysklad: create processing %s %s: assortment meta.href and quantity must be set", c.ParentID, c.Item)
	}
	return nested.Create[Position](ctx, c, input, opts)
}

// GetMaterials читает одну страницу материалов: только limit и offset (fields отклоняется до запроса).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-poluchit-materialy-tehoperacii
func (s *Service) GetMaterials(ctx context.Context, id string, p positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return page(ctx, s.materials(id), p, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-poluchit-material
func (s *Service) GetMaterial(ctx context.Context, id, materialID string, opts ...transport.RequestOption) (*Position, error) {
	return nested.Get[Position](ctx, s.materials(id), materialID, opts)
}

// CreateMaterial добавляет один материал одним POST; Assortment и Quantity обязательны.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-dobavit-material
func (s *Service) CreateMaterial(ctx context.Context, id string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return create(ctx, s.materials(id), input, opts)
}

// UpdateMaterial: обязательных полей нет, пустой input отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-izmenit-material
func (s *Service) UpdateMaterial(ctx context.Context, id, materialID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.materials(id), materialID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-udalit-material
func (s *Service) DeleteMaterial(ctx context.Context, id, materialID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.materials(id), materialID, opts)
}

// DeleteManyMaterials — POST …/materials/delete с meta строк, 1–1000 id. Источник отвечает 200 с
// пустым телом, поэтому результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-massovoe-udalenie-materialov
func (s *Service) DeleteManyMaterials(ctx context.Context, id string, materialIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.materials(id), "delete processing "+id+" materials", materialIDs, opts)
}

// GetProducts читает одну страницу продуктов: только limit и offset (fields отклоняется до запроса).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-poluchit-produkty-tehoperacii
func (s *Service) GetProducts(ctx context.Context, id string, p positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return page(ctx, s.products(id), p, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-poluchit-produkt
func (s *Service) GetProduct(ctx context.Context, id, productID string, opts ...transport.RequestOption) (*Position, error) {
	return nested.Get[Position](ctx, s.products(id), productID, opts)
}

// CreateProduct добавляет один продукт одним POST; Assortment и Quantity обязательны.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-dobavit-produkt
func (s *Service) CreateProduct(ctx context.Context, id string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return create(ctx, s.products(id), input, opts)
}

// UpdateProduct: обязательных полей нет, пустой input отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-izmenit-produkt
func (s *Service) UpdateProduct(ctx context.Context, id, productID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.products(id), productID, input, opts)
}

// DeleteProduct: техоперация должна сохранить хотя бы один продукт — это проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-udalit-produkt
func (s *Service) DeleteProduct(ctx context.Context, id, productID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.products(id), productID, opts)
}

// DeleteManyProducts — POST …/products/delete с meta строк, 1–1000 id; ответ без тела.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-massovoe-udalenie-produktov
func (s *Service) DeleteManyProducts(ctx context.Context, id string, productIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.products(id), "delete processing "+id+" products", productIDs, opts)
}
