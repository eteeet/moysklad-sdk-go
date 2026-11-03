package bundle

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// ComponentAddRequest — новый компонент комплекта: assortment обязателен (проверяется до
// запроса), quantity передаётся всегда, как в примере источника; лимит 50 компонентов проверяет
// сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-dobavit-komponent-komplekta
type ComponentAddRequest struct {
	Assortment model.Reference `json:"assortment"`
	Quantity   float64         `json:"quantity"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-izmenit-komponent
type ComponentUpdateRequest struct {
	Assortment model.Field[model.Reference] `json:"assortment,omitzero"`
	Quantity   model.Field[float64]         `json:"quantity,omitzero"`
}

func (s *Service) components(bundleID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "bundle", ParentID: bundleID, Path: "components", Item: "component", MetaType: "bundlecomponent"}
}

// AddComponent — POST entity/bundle/{id}/components. Не повторять после неоднозначного сбоя:
// возможен дубль компонента.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-dobavit-komponent-komplekta
func (s *Service) AddComponent(ctx context.Context, bundleID string, input ComponentAddRequest, opts ...transport.RequestOption) (*Component, error) {
	if input.Assortment.Meta.Href == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create bundle %s component: assortment meta.href must not be empty", bundleID)
	}
	return nested.Create[Component](ctx, s.components(bundleID), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-izmenit-komponent
func (s *Service) UpdateComponent(ctx context.Context, bundleID, componentID string, input ComponentUpdateRequest, opts ...transport.RequestOption) (*Component, error) {
	return nested.Update[Component](ctx, s.components(bundleID), componentID, input, opts)
}

// DeleteComponent — DELETE компонента. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-udalit-komponent
func (s *Service) DeleteComponent(ctx context.Context, bundleID, componentID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.components(bundleID), componentID, opts)
}

// DeleteManyComponents — POST entity/bundle/{id}/components/delete с meta компонентов
// (тип bundlecomponent). Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-massovoe-udalenie-komponentov-komplekta
func (s *Service) DeleteManyComponents(ctx context.Context, bundleID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.components(bundleID), ids, opts)
}
