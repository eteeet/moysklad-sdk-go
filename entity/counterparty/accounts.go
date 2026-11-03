package counterparty

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Account — расчётный счёт контрагента (meta.type account). Валюты у счёта контрагента
// источник не описывает, в отличие от счёта юрлица.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-scheta-kontragenta
type Account struct {
	Meta                 model.Meta     `json:"meta"`
	ID                   string         `json:"id"`
	AccountID            string         `json:"accountId"`
	Updated              model.DateTime `json:"updated,omitzero"`
	IsDefault            bool           `json:"isDefault"`
	AccountNumber        string         `json:"accountNumber"`
	BankName             string         `json:"bankName"`
	BankLocation         string         `json:"bankLocation"`
	CorrespondentAccount string         `json:"correspondentAccount"`
	BIC                  string         `json:"bic"`
}

// AccountCreateRequest: обязателен только accountNumber.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-schet-kontragenta
type AccountCreateRequest struct {
	AccountNumber        string              `json:"accountNumber"`
	IsDefault            model.Field[bool]   `json:"isDefault,omitzero"`
	BankName             model.Field[string] `json:"bankName,omitzero"`
	BankLocation         model.Field[string] `json:"bankLocation,omitzero"`
	CorrespondentAccount model.Field[string] `json:"correspondentAccount,omitzero"`
	BIC                  model.Field[string] `json:"bic,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-schet-kontragenta
type AccountUpdateRequest struct {
	AccountNumber        model.Field[string] `json:"accountNumber,omitzero"`
	IsDefault            model.Field[bool]   `json:"isDefault,omitzero"`
	BankName             model.Field[string] `json:"bankName,omitzero"`
	BankLocation         model.Field[string] `json:"bankLocation,omitzero"`
	CorrespondentAccount model.Field[string] `json:"correspondentAccount,omitzero"`
	BIC                  model.Field[string] `json:"bic,omitzero"`
}

func (s *Service) accounts(counterpartyID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "counterparty", ParentID: counterpartyID, Path: "accounts", Item: "account", MetaType: "account"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-scheta-kontragenta
func (s *Service) ListAccounts(ctx context.Context, counterpartyID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Account], error) {
	return nested.List[Account](ctx, s.accounts(counterpartyID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-schet-kontragenta
func (s *Service) GetAccount(ctx context.Context, counterpartyID, accountID string, opts ...transport.RequestOption) (*Account, error) {
	return nested.Get[Account](ctx, s.accounts(counterpartyID), accountID, opts)
}

// CreateAccount — POST entity/counterparty/{id}/accounts. Пример ответа источника — массив
// из одного счёта; принимается и объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-schet-kontragenta
func (s *Service) CreateAccount(ctx context.Context, counterpartyID string, input AccountCreateRequest, opts ...transport.RequestOption) (*Account, error) {
	if input.AccountNumber == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create counterparty %s account: accountNumber must not be empty", counterpartyID)
	}
	return nested.Create[Account](ctx, s.accounts(counterpartyID), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-schet-kontragenta
func (s *Service) UpdateAccount(ctx context.Context, counterpartyID, accountID string, input AccountUpdateRequest, opts ...transport.RequestOption) (*Account, error) {
	return nested.Update[Account](ctx, s.accounts(counterpartyID), accountID, input, opts)
}

// DeleteAccount — DELETE счёта. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-udalit-schet-kontragenta
func (s *Service) DeleteAccount(ctx context.Context, counterpartyID, accountID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.accounts(counterpartyID), accountID, opts)
}

// DeleteManyAccounts — POST entity/counterparty/{id}/accounts/delete с meta счетов. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-massovoe-udalenie-schetov-kontragentov
func (s *Service) DeleteManyAccounts(ctx context.Context, counterpartyID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.accounts(counterpartyID), ids, opts)
}
