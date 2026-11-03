package productionstagecompletion

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Материалы — entity/productionstagecompletion/{id}/materials (список, добавление, изменение),
// продукты — …/{id}/products (список, изменение). Получения одной строки, удаления и добавления
// продуктов источник не описывает. Повторов нет.

func (s *Service) materials(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: entity, ParentID: id, Path: "materials", Item: "material", MetaType: "productionstagecompletionmaterial"}
}

func (s *Service) products(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: entity, ParentID: id, Path: "products", Item: "product", MetaType: "productionstagecompletionresult"}
}

// GetMaterials читает одну страницу материалов: только limit и offset, без скрытого Get и
// автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-poluchit-materialy-vypolneniya-etapa-proizvodstva
func (s *Service) GetMaterials(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Material], error) {
	return nested.List[Material](ctx, s.materials(id), page, opts)
}

// AddMaterial — POST …/{id}/materials одним объектом. Не повторять после неоднозначного сбоя:
// возможен дубль материала. Серийные номера при создании выполнения этапа не поддерживаются,
// но добавление материала их принимает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-dobavit-material-vypolneniya-etapa-proizvodstva
func (s *Service) AddMaterial(ctx context.Context, id string, input MaterialAddRequest, opts ...transport.RequestOption) (*Material, error) {
	if input.Assortment.Meta.Href == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create productionstagecompletion %s material: assortment meta.href must not be empty", id)
	}
	return nested.Create[Material](ctx, s.materials(id), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-izmenit-material-vypolneniya-etapa-proizvodstva
func (s *Service) UpdateMaterial(ctx context.Context, id, materialID string, input MaterialUpdateRequest, opts ...transport.RequestOption) (*Material, error) {
	return nested.Update[Material](ctx, s.materials(id), materialID, input, opts)
}

// GetProducts читает одну страницу продуктов: только limit и offset, без скрытого Get и
// автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-poluchit-produkty-vypolneniya-etapa-proizvodstva
func (s *Service) GetProducts(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Product], error) {
	return nested.List[Product](ctx, s.products(id), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-izmenit-produkt-vypolneniya-etapa-proizvodstva
func (s *Service) UpdateProduct(ctx context.Context, id, productID string, input ProductUpdateRequest, opts ...transport.RequestOption) (*Product, error) {
	return nested.Update[Product](ctx, s.products(id), productID, input, opts)
}
