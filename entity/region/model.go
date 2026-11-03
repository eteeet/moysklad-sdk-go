package region

import "github.com/eteeet/moysklad-sdk-go/model"

// Region — регион; справочник только для чтения. Version — версия сущности.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/region
type Region struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	AccountID    string         `json:"accountId"`
	Updated      model.DateTime `json:"updated,omitzero"`
	Name         string         `json:"name"`
	Code         string         `json:"code"`
	ExternalCode string         `json:"externalCode"`
	Version      int            `json:"version"`
}
