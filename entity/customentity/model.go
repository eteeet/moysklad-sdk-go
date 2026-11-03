package customentity

import "github.com/eteeet/moysklad-sdk-go/model"

// CustomEntity — пользовательский справочник. Список справочников и их доп. поля читает
// metadata.Service (GetCompanySettings, GetCustomEntity); ID — он же id метаданных справочника.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#2-polzovatelskij-spravochnik
type CustomEntity struct {
	Meta model.Meta `json:"meta"`
	ID   string     `json:"id"`
	Name string     `json:"name"`
}

// Element — элемент пользовательского справочника (meta.type customentity). Доп. поля
// описаны в метаданных справочника.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-elementy-polzovatelskogo-spravochnika
type Element struct {
	Meta         model.Meta        `json:"meta"`
	ID           string            `json:"id"`
	AccountID    string            `json:"accountId"`
	Updated      model.DateTime    `json:"updated,omitzero"`
	Name         string            `json:"name"`
	Code         string            `json:"code"`
	Description  string            `json:"description"`
	ExternalCode string            `json:"externalCode"`
	Attributes   []model.Attribute `json:"attributes"`
	Shared       bool              `json:"shared"`
	Owner        *model.Reference  `json:"owner,omitzero"`
	Group        *model.Reference  `json:"group,omitzero"`
}
