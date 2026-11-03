package processingplanfolder

import "github.com/eteeet/moysklad-sdk-go/model"

// ProcessingPlanFolder — группа техкарт: все 13 полей таблицы атрибутов источника.
// PathName — имя родительской группы, только чтение. Ссылки на родительскую группу в таблице
// нет (текст раздела «Изменить» упоминает поле processingplanfolder), поэтому модель её не читает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-gruppy-tehkart
type ProcessingPlanFolder struct {
	Meta         model.Meta       `json:"meta"`
	ID           string           `json:"id"`
	AccountID    string           `json:"accountId"`
	Updated      model.DateTime   `json:"updated,omitzero"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Code         string           `json:"code"`
	ExternalCode string           `json:"externalCode"`
	Archived     bool             `json:"archived"`
	PathName     string           `json:"pathName"`
	Shared       bool             `json:"shared"`
	Owner        *model.Reference `json:"owner,omitzero"`
	Group        *model.Reference `json:"group,omitzero"`
}
