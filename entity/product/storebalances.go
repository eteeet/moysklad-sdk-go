package product

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Неснижаемые остатки по складам (тип WAREHOUSE_VARIED): entity/product/{id}/storebalances.
// Источник перечисляет только маршруты (список, объект, создание, изменение, удаление, массовое
// удаление); поля — по таблице позиции storeBalances: store и quantity. Create принимает
// ответ-объект и массив из одного объекта.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-nesnizhaemyj-ostatok
func (s *Service) storeBalances(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "product", ParentID: id, Path: "storebalances", Item: "storebalance", MetaType: "minimumstock"}
}

func (s *Service) ListStoreBalances(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[model.StoreBalance], error) {
	return nested.List[model.StoreBalance](ctx, s.storeBalances(id), page, opts)
}

func (s *Service) GetStoreBalance(ctx context.Context, id, balanceID string, opts ...transport.RequestOption) (*model.StoreBalance, error) {
	return nested.Get[model.StoreBalance](ctx, s.storeBalances(id), balanceID, opts)
}

// CreateStoreBalance — POST с store и quantity; повторов нет.
func (s *Service) CreateStoreBalance(ctx context.Context, id string, input model.StoreBalanceRequest, opts ...transport.RequestOption) (*model.StoreBalance, error) {
	return nested.Create[model.StoreBalance](ctx, s.storeBalances(id), input, opts)
}

func (s *Service) UpdateStoreBalance(ctx context.Context, id, balanceID string, input model.StoreBalanceRequest, opts ...transport.RequestOption) (*model.StoreBalance, error) {
	return nested.Update[model.StoreBalance](ctx, s.storeBalances(id), balanceID, input, opts)
}

// DeleteStoreBalance — DELETE позиции. Повторов нет.
func (s *Service) DeleteStoreBalance(ctx context.Context, id, balanceID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.storeBalances(id), balanceID, opts)
}

// DeleteManyStoreBalances — POST …/storebalances/delete с meta позиций (тип minimumstock).
func (s *Service) DeleteManyStoreBalances(ctx context.Context, id string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.storeBalances(id), ids, opts)
}
