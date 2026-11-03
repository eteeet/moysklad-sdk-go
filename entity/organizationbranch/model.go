package organizationbranch

import "github.com/eteeet/moysklad-sdk-go/model"

// OrganizationBranch — подразделение юрлица. KPP — реквизит РФ (+RU): валиден только у подразделений
// юрлица РФ, проверяет сервер. Реквизиты Organization по аналогии не копируются.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organizationbranch
type OrganizationBranch struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	Updated           model.DateTime   `json:"updated,omitzero"`
	Created           model.DateTime   `json:"created,omitzero"`
	Name              string           `json:"name"`
	Organization      *model.Reference `json:"organization,omitzero"`
	Description       string           `json:"description"`
	Code              string           `json:"code"`
	ExternalCode      string           `json:"externalCode"`
	Archived          bool             `json:"archived"`
	Email             string           `json:"email"`
	Phone             string           `json:"phone"`
	Fax               string           `json:"fax"`
	ActualAddress     string           `json:"actualAddress"`
	ActualAddressFull *model.Address   `json:"actualAddressFull,omitzero"`
	// Реквизиты РФ: так в 8 из 9 примеров источника; таблица и один пример ответа показывают kpp на верхнем уровне.
	RequisitesRU *RequisitesRU    `json:"mod__requisites__ru,omitzero"`
	SyncID       string           `json:"syncId"`
	Shared       bool             `json:"shared"`
	Owner        *model.Reference `json:"owner,omitzero"`
	Group        *model.Reference `json:"group,omitzero"`
}

// RequisitesRU — реквизиты подразделения юрлица РФ (+RU).
type RequisitesRU struct {
	KPP string `json:"kpp"`
}

type RequisitesRURequest struct {
	KPP model.Field[string] `json:"kpp,omitzero"`
}
