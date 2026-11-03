package organizationbranch

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязательны name и organization. Organization после заполнения не меняется —
// его нет в UpdateRequest; syncId сервер принимает однократно, повторную смену отклоняет сам.
// Адрес — строкой или структурой; null адреса источник не поддерживает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-sozdat-podrazdelenie-yuridicheskogo-lica
type CreateRequest struct {
	Name              string                            `json:"name"`
	Organization      model.Reference                   `json:"organization"`
	Description       model.Field[string]               `json:"description,omitzero"`
	Code              model.Field[string]               `json:"code,omitzero"`
	ExternalCode      model.Field[string]               `json:"externalCode,omitzero"`
	Archived          model.Field[bool]                 `json:"archived,omitzero"`
	Email             model.Field[string]               `json:"email,omitzero"`
	Phone             model.Field[string]               `json:"phone,omitzero"`
	Fax               model.Field[string]               `json:"fax,omitzero"`
	ActualAddress     model.Field[string]               `json:"actualAddress,omitzero"`
	ActualAddressFull model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	RequisitesRU      model.Field[RequisitesRURequest]  `json:"mod__requisites__ru,omitzero"`
	SyncID            model.Field[string]               `json:"syncId,omitzero"`
	Shared            model.Field[bool]                 `json:"shared,omitzero"`
	Owner             model.Field[model.Reference]      `json:"owner,omitzero"`
	Group             model.Field[model.Reference]      `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch#3-izmenit-podrazdelenie-yuridicheskogo-lica
type UpdateRequest struct {
	Name              model.Field[string]               `json:"name,omitzero"`
	Description       model.Field[string]               `json:"description,omitzero"`
	Code              model.Field[string]               `json:"code,omitzero"`
	ExternalCode      model.Field[string]               `json:"externalCode,omitzero"`
	Archived          model.Field[bool]                 `json:"archived,omitzero"`
	Email             model.Field[string]               `json:"email,omitzero"`
	Phone             model.Field[string]               `json:"phone,omitzero"`
	Fax               model.Field[string]               `json:"fax,omitzero"`
	ActualAddress     model.Field[string]               `json:"actualAddress,omitzero"`
	ActualAddressFull model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	RequisitesRU      model.Field[RequisitesRURequest]  `json:"mod__requisites__ru,omitzero"`
	SyncID            model.Field[string]               `json:"syncId,omitzero"`
	Shared            model.Field[bool]                 `json:"shared,omitzero"`
	Owner             model.Field[model.Reference]      `json:"owner,omitzero"`
	Group             model.Field[model.Reference]      `json:"group,omitzero"`
}
