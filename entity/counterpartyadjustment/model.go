package counterpartyadjustment

import "github.com/eteeet/moysklad-sdk-go/model"

// CounterpartyAdjustment — корректировка взаиморасчётов. Модель чтения без expand и без позиций.
// Sum — копейки, readonly по таблице атрибутов; знак суммы источник не описывает (все примеры
// положительные), SDK значение не меняет и не проверяет. Agent — контрагент или сотрудник.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment
type CounterpartyAdjustment struct {
	Meta         model.Meta                         `json:"meta"`
	ID           string                             `json:"id"`
	AccountID    string                             `json:"accountId"`
	Name         string                             `json:"name"`
	Description  string                             `json:"description"`
	ExternalCode string                             `json:"externalCode"`
	Created      model.DateTime                     `json:"created,omitzero"`
	Updated      model.DateTime                     `json:"updated,omitzero"`
	Deleted      model.DateTime                     `json:"deleted,omitzero"`
	Moment       model.DateTime                     `json:"moment,omitzero"`
	Applicable   bool                               `json:"applicable"`
	Shared       bool                               `json:"shared"`
	Printed      bool                               `json:"printed"`
	Published    bool                               `json:"published"`
	Organization *model.Reference                   `json:"organization,omitzero"`
	Agent        *model.Reference                   `json:"agent,omitzero"`
	Owner        *model.Reference                   `json:"owner,omitzero"`
	Group        *model.Reference                   `json:"group,omitzero"`
	Attributes   []model.Attribute                  `json:"attributes"`
	Sum          float64                            `json:"sum"`
	Files        *model.Collection[model.Reference] `json:"files,omitzero"`
}
