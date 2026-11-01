package employee

import "github.com/eteeet/moysklad-sdk-go/model"

// Salary, Owner, Group, Archived и Email может менять только администратор; права проверяет сервер.
// Фото, активация и пароль сотрудника этими запросами не меняются.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-sozdat-sotrudnika
type CreateRequest struct {
	LastName     string              `json:"lastName"`
	FirstName    model.Field[string] `json:"firstName,omitzero"`
	MiddleName   model.Field[string] `json:"middleName,omitzero"`
	Description  model.Field[string] `json:"description,omitzero"`
	ExternalCode model.Field[string] `json:"externalCode,omitzero"`
	Archived     model.Field[bool]   `json:"archived,omitzero"`
	Shared       model.Field[bool]   `json:"shared,omitzero"`
	Email        model.Field[string] `json:"email,omitzero"`
	Phone        model.Field[string] `json:"phone,omitzero"`
	INN          model.Field[string] `json:"inn,omitzero"`
	Position     model.Field[string] `json:"position,omitzero"`

	Salary     model.Field[Salary]            `json:"salary,omitzero"`
	Owner      model.Field[model.Reference]   `json:"owner,omitzero"`
	Group      model.Field[model.Reference]   `json:"group,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
}

// LastName обязателен и в Update: так требует API.
// null для полей источник не описывает; Field позволяет его отправить, решает сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee#3-izmenit-sotrudnika
type UpdateRequest struct {
	LastName     string              `json:"lastName"`
	FirstName    model.Field[string] `json:"firstName,omitzero"`
	MiddleName   model.Field[string] `json:"middleName,omitzero"`
	Description  model.Field[string] `json:"description,omitzero"`
	ExternalCode model.Field[string] `json:"externalCode,omitzero"`
	Archived     model.Field[bool]   `json:"archived,omitzero"`
	Shared       model.Field[bool]   `json:"shared,omitzero"`
	Email        model.Field[string] `json:"email,omitzero"`
	Phone        model.Field[string] `json:"phone,omitzero"`
	INN          model.Field[string] `json:"inn,omitzero"`
	Position     model.Field[string] `json:"position,omitzero"`

	Salary     model.Field[Salary]            `json:"salary,omitzero"`
	Owner      model.Field[model.Reference]   `json:"owner,omitzero"`
	Group      model.Field[model.Reference]   `json:"group,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
}
