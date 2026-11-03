package productiontask

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Product — продукт производственного задания, meta.type productiontaskresult. PlanQuantity —
// количество, не деньги; ProductionRow — позиция задания, только для чтения.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-produkty-proizvodstvennogo-zadaniya
type Product struct {
	Meta          model.Meta       `json:"meta"`
	ID            string           `json:"id"`
	AccountID     string           `json:"accountId"`
	Assortment    *model.Reference `json:"assortment,omitzero"`
	PlanQuantity  float64          `json:"planQuantity"`
	ProductionRow *model.Reference `json:"productionRow,omitzero"`
}

// ProductRequest — продукт в Create/Update документа, CreateProducts и UpdateProduct. Раздел
// «Создать продукт» требует assortment, положительный planQuantity и productionRow, но таблица
// помечает productionRow «Только для чтения» — SDK его не пишет; обязательность проверяет сервер.
type ProductRequest struct {
	// ID сопоставляет существующий продукт в Update документа; для нового продукта его не задают.
	ID           model.Field[string]          `json:"id,omitzero"`
	Assortment   model.Field[model.Reference] `json:"assortment,omitzero"`
	PlanQuantity model.Field[float64]         `json:"planQuantity,omitzero"`
}

func (s *Service) products(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "productiontask", ParentID: id, Path: "products", Item: "product", MetaType: "productiontaskresult"}
}

// GetProducts читает одну страницу продуктов: только limit и offset, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-produkty-proizvodstvennogo-zadaniya
func (s *Service) GetProducts(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Product], error) {
	return getPage[Product](ctx, s.client, "get productiontask "+id+" products", id, "entity/productiontask/"+id+"/products", page, "", opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-produkt-proizvodstvennogo-zadaniya
func (s *Service) GetProduct(ctx context.Context, id, productID string, opts ...transport.RequestOption) (*Product, error) {
	return nested.Get[Product](ctx, s.products(id), productID, opts)
}

// CreateProducts добавляет продукты к существующим одним POST массивом (один продукт — массив из
// одного элемента); ответ — по правилам batch. Пример источника отправляет POST на …/{id}/product,
// а получение, изменение, удаление и meta.href коллекции используют …/{id}/products — SDK
// отправляет на …/products. Не повторять после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-sozdat-produkt
func (s *Service) CreateProducts(ctx context.Context, id string, items []ProductRequest, opts ...transport.RequestOption) ([]batch.Result[Product], error) {
	return nested.CreateMany[Product](ctx, s.products(id), items, opts)
}

// UpdateProduct: обязательных полей нет; пустой input отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-produkt
func (s *Service) UpdateProduct(ctx context.Context, id, productID string, input ProductRequest, opts ...transport.RequestOption) (*Product, error) {
	return nested.Update[Product](ctx, s.products(id), productID, input, opts)
}

// DeleteProduct: последний продукт позиции и продукты после начала финальной позиции сервер не
// удаляет. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-udalit-produkt
func (s *Service) DeleteProduct(ctx context.Context, id, productID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.products(id), productID, opts)
}

// DeleteManyProducts — POST …/products/delete с meta продуктов (1–1000 id). Источник отвечает 200
// с пустым телом, поэтому результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-massovoe-udalenie-produktov-proizvodstvennogo-zadaniya
func (s *Service) DeleteManyProducts(ctx context.Context, id string, productIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.products(id), "delete productiontask "+id+" products", productIDs, opts)
}
