package project

import "github.com/eteeet/moysklad-sdk-go/model"

// Project — проект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/project
type Project struct {
	Meta         model.Meta        `json:"meta"`
	ID           string            `json:"id"`
	AccountID    string            `json:"accountId"`
	Updated      model.DateTime    `json:"updated,omitzero"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Code         string            `json:"code"`
	ExternalCode string            `json:"externalCode"`
	Archived     bool              `json:"archived"`
	Attributes   []model.Attribute `json:"attributes"`
	Shared       bool              `json:"shared"`
	Owner        *model.Reference  `json:"owner,omitzero"`
	Group        *model.Reference  `json:"group,omitzero"`
}
