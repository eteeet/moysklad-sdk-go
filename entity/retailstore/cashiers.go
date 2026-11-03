package retailstore

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Cashier — кассир точки продаж (meta.type cashier); только чтение: записи источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/cashier#2-kassir
type Cashier struct {
	Meta        model.Meta       `json:"meta"`
	ID          string           `json:"id"`
	AccountID   string           `json:"accountId"`
	Employee    *model.Reference `json:"employee,omitzero"`
	RetailStore *model.Reference `json:"retailStore,omitzero"`
}

func (s *Service) cashiers(retailStoreID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "retailstore", ParentID: retailStoreID, Path: "cashiers", Item: "cashier", MetaType: "cashier"}
}

// ListCashiers — GET entity/retailstore/{id}/cashiers, одна страница (limit/offset).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/cashier#3-poluchit-kassirov
func (s *Service) ListCashiers(ctx context.Context, retailStoreID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Cashier], error) {
	return nested.List[Cashier](ctx, s.cashiers(retailStoreID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/cashier#3-poluchit-kassira
func (s *Service) GetCashier(ctx context.Context, retailStoreID, cashierID string, opts ...transport.RequestOption) (*Cashier, error) {
	return nested.Get[Cashier](ctx, s.cashiers(retailStoreID), cashierID, opts)
}
