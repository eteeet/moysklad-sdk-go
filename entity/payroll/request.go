package payroll

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, startPayrollPeriod и endPayrollPeriod; разница начала
// и окончания периода должна быть кратна 24 часам (секунды игнорируются) — проверяет сервер.
// Positions — не больше 1000 строк. Payments не пишется: источник помечает связанные платежи
// «Только для чтения».
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-sozdat-nachislenie-zarplaty
type CreateRequest struct {
	Organization       model.Reference                `json:"organization"`
	StartPayrollPeriod model.DateTime                 `json:"startPayrollPeriod"`
	EndPayrollPeriod   model.DateTime                 `json:"endPayrollPeriod"`
	Name               model.Field[string]            `json:"name,omitzero"`
	Description        model.Field[string]            `json:"description,omitzero"`
	ExternalCode       model.Field[string]            `json:"externalCode,omitzero"`
	SyncID             model.Field[string]            `json:"syncId,omitzero"`
	Moment             model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable         model.Field[bool]              `json:"applicable,omitzero"`
	Shared             model.Field[bool]              `json:"shared,omitzero"`
	Owner              model.Field[model.Reference]   `json:"owner,omitzero"`
	Group              model.Field[model.Reference]   `json:"group,omitzero"`
	State              model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes         model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Positions          model.Field[[]PositionRequest] `json:"positions,omitzero"`
}

// UpdateRequest. Positions заменяют все строки: строки без id — новые. syncId после заполнения
// не меняется (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-izmenit-nachislenie-zarplaty
type UpdateRequest struct {
	Organization       model.Field[model.Reference]   `json:"organization,omitzero"`
	StartPayrollPeriod model.Field[model.DateTime]    `json:"startPayrollPeriod,omitzero"`
	EndPayrollPeriod   model.Field[model.DateTime]    `json:"endPayrollPeriod,omitzero"`
	Name               model.Field[string]            `json:"name,omitzero"`
	Description        model.Field[string]            `json:"description,omitzero"`
	ExternalCode       model.Field[string]            `json:"externalCode,omitzero"`
	SyncID             model.Field[string]            `json:"syncId,omitzero"`
	Moment             model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable         model.Field[bool]              `json:"applicable,omitzero"`
	Shared             model.Field[bool]              `json:"shared,omitzero"`
	Owner              model.Field[model.Reference]   `json:"owner,omitzero"`
	Group              model.Field[model.Reference]   `json:"group,omitzero"`
	State              model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes         model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Positions          model.Field[[]PositionRequest] `json:"positions,omitzero"`
}

// PositionRequest — строка в Create/Update документа, CreatePositions и UpdatePosition. Для новой
// строки источник требует employee, baseSalary и pieceworkSalary (положительные, округляются до
// целого) — проверяет сервер; в UpdatePosition обязательных полей нет.
type PositionRequest struct {
	// ID сопоставляет существующую строку в Update документа; для новой строки его не задают.
	ID              model.Field[string]          `json:"id,omitzero"`
	Employee        model.Field[model.Reference] `json:"employee,omitzero"`
	BaseSalary      model.Field[float64]         `json:"baseSalary,omitzero"`
	PieceworkSalary model.Field[float64]         `json:"pieceworkSalary,omitzero"`
}
