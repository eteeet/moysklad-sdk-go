package counterparty

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// ContactPerson — контактное лицо контрагента (meta.type contactperson). Agent — ссылка на
// контрагента: маркера «только для чтения» нет, но родителя задаёт путь запроса, поэтому SDK
// agent не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-kontaktnye-lica-kontragenta
type ContactPerson struct {
	Meta         model.Meta       `json:"meta"`
	ID           string           `json:"id"`
	AccountID    string           `json:"accountId"`
	Updated      model.DateTime   `json:"updated,omitzero"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	ExternalCode string           `json:"externalCode"`
	Email        string           `json:"email"`
	Phone        string           `json:"phone"`
	Position     string           `json:"position"`
	Agent        *model.Reference `json:"agent,omitzero"`
}

// ContactPersonCreateRequest: обязательно name (ФИО).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-kontaktnoe-lico
type ContactPersonCreateRequest struct {
	Name         string              `json:"name"`
	Description  model.Field[string] `json:"description,omitzero"`
	ExternalCode model.Field[string] `json:"externalCode,omitzero"`
	Email        model.Field[string] `json:"email,omitzero"`
	Phone        model.Field[string] `json:"phone,omitzero"`
	Position     model.Field[string] `json:"position,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-kontaktnoe-lico
type ContactPersonUpdateRequest struct {
	Name         model.Field[string] `json:"name,omitzero"`
	Description  model.Field[string] `json:"description,omitzero"`
	ExternalCode model.Field[string] `json:"externalCode,omitzero"`
	Email        model.Field[string] `json:"email,omitzero"`
	Phone        model.Field[string] `json:"phone,omitzero"`
	Position     model.Field[string] `json:"position,omitzero"`
}

// Удаления контактного лица источник не описывает — метода нет.
func (s *Service) contactPersons(counterpartyID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "counterparty", ParentID: counterpartyID, Path: "contactpersons", Item: "contactperson", MetaType: "contactperson"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-spisok-kontaktnyh-lic
func (s *Service) ListContactPersons(ctx context.Context, counterpartyID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[ContactPerson], error) {
	return nested.List[ContactPerson](ctx, s.contactPersons(counterpartyID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-kontaktnoe-lico
func (s *Service) GetContactPerson(ctx context.Context, counterpartyID, personID string, opts ...transport.RequestOption) (*ContactPerson, error) {
	return nested.Get[ContactPerson](ctx, s.contactPersons(counterpartyID), personID, opts)
}

// CreateContactPerson — пример ответа источника — массив из одного объекта; принимается и объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-kontaktnoe-lico
func (s *Service) CreateContactPerson(ctx context.Context, counterpartyID string, input ContactPersonCreateRequest, opts ...transport.RequestOption) (*ContactPerson, error) {
	if input.Name == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create counterparty %s contactperson: name must not be empty", counterpartyID)
	}
	return nested.Create[ContactPerson](ctx, s.contactPersons(counterpartyID), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-kontaktnoe-lico
func (s *Service) UpdateContactPerson(ctx context.Context, counterpartyID, personID string, input ContactPersonUpdateRequest, opts ...transport.RequestOption) (*ContactPerson, error) {
	return nested.Update[ContactPerson](ctx, s.contactPersons(counterpartyID), personID, input, opts)
}
