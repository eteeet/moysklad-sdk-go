package organization

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Account — расчётный счёт юрлица (meta.type account). Currency — валюта счёта: от неё зависит
// проверка валюты платежей (accounts-in-currency).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-zaprosy-scheta-yurlica
type Account struct {
	Meta                 model.Meta       `json:"meta"`
	ID                   string           `json:"id"`
	AccountID            string           `json:"accountId"`
	Updated              model.DateTime   `json:"updated,omitzero"`
	IsDefault            bool             `json:"isDefault"`
	AccountNumber        string           `json:"accountNumber"`
	BankName             string           `json:"bankName"`
	BankLocation         string           `json:"bankLocation"`
	CorrespondentAccount string           `json:"correspondentAccount"`
	BIC                  string           `json:"bic"`
	Currency             *model.Reference `json:"currency,omitzero"`
}

// AccountCreateRequest: обязателен accountNumber; без currency сервер выбирает валюту сам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-sozdat-schet-yurlica
type AccountCreateRequest struct {
	AccountNumber        string                       `json:"accountNumber"`
	IsDefault            model.Field[bool]            `json:"isDefault,omitzero"`
	BankName             model.Field[string]          `json:"bankName,omitzero"`
	BankLocation         model.Field[string]          `json:"bankLocation,omitzero"`
	CorrespondentAccount model.Field[string]          `json:"correspondentAccount,omitzero"`
	BIC                  model.Field[string]          `json:"bic,omitzero"`
	Currency             model.Field[model.Reference] `json:"currency,omitzero"`
}

// AccountUpdateRequest: валюту существующего счёта сервер менять не даёт и возвращает ошибку;
// SDK поле передаёт как есть, не сравнивая с текущей валютой.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-izmenit-schet-yurlica
type AccountUpdateRequest struct {
	AccountNumber        model.Field[string]          `json:"accountNumber,omitzero"`
	IsDefault            model.Field[bool]            `json:"isDefault,omitzero"`
	BankName             model.Field[string]          `json:"bankName,omitzero"`
	BankLocation         model.Field[string]          `json:"bankLocation,omitzero"`
	CorrespondentAccount model.Field[string]          `json:"correspondentAccount,omitzero"`
	BIC                  model.Field[string]          `json:"bic,omitzero"`
	Currency             model.Field[model.Reference] `json:"currency,omitzero"`
}

func (s *Service) accounts(organizationID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "organization", ParentID: organizationID, Path: "accounts", Item: "account", MetaType: "account"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-poluchit-spisok-schetov-yurlica
func (s *Service) ListAccounts(ctx context.Context, organizationID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Account], error) {
	return nested.List[Account](ctx, s.accounts(organizationID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-poluchit-schet-yurlica
func (s *Service) GetAccount(ctx context.Context, organizationID, accountID string, opts ...transport.RequestOption) (*Account, error) {
	return nested.Get[Account](ctx, s.accounts(organizationID), accountID, opts)
}

// CreateAccount — POST entity/organization/{id}/accounts с одним счётом. Пример ответа
// источника — массив из одного счёта; принимается и объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-sozdat-schet-yurlica
func (s *Service) CreateAccount(ctx context.Context, organizationID string, input AccountCreateRequest, opts ...transport.RequestOption) (*Account, error) {
	if input.AccountNumber == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create organization %s account: accountNumber must not be empty", organizationID)
	}
	return nested.Create[Account](ctx, s.accounts(organizationID), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-izmenit-schet-yurlica
func (s *Service) UpdateAccount(ctx context.Context, organizationID, accountID string, input AccountUpdateRequest, opts ...transport.RequestOption) (*Account, error) {
	return nested.Update[Account](ctx, s.accounts(organizationID), accountID, input, opts)
}

// DeleteAccount — DELETE счёта. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-udalit-schet-yurlica
func (s *Service) DeleteAccount(ctx context.Context, organizationID, accountID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.accounts(organizationID), accountID, opts)
}

// CreateManyAccounts и UpdateManyAccounts — «Изменить счета Юрлица»: один POST
// entity/organization/{id}/accounts с массивом; элемент без meta создаёт счёт, с meta — изменяет.
// Непереданные поля не меняются, полного замещения списка счетов нет. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-izmenit-scheta-yurlica
func (s *Service) CreateManyAccounts(ctx context.Context, organizationID string, items []AccountCreateRequest, opts ...transport.RequestOption) ([]batch.Result[Account], error) {
	return nested.CreateMany[Account](ctx, s.accounts(organizationID), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-izmenit-scheta-yurlica
func (s *Service) UpdateManyAccounts(ctx context.Context, organizationID string, items []batch.Update[AccountUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Account], error) {
	return nested.UpdateMany[Account](ctx, s.accounts(organizationID), items, opts)
}

// DeleteManyAccounts — POST entity/organization/{id}/accounts/delete с meta счетов. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-massovoe-udalenie-schetov-yurlica
func (s *Service) DeleteManyAccounts(ctx context.Context, organizationID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.accounts(organizationID), ids, opts)
}
