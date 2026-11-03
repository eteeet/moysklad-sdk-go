package payroll

import "github.com/eteeet/moysklad-sdk-go/model"

// Payroll — начисление зарплаты. Модель чтения без expand: Sum — копейки, только для чтения.
// Payments — связанные исходящие платежи и расходные ордера (paymentout, cashout), только чтение.
// Строки — GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll
type Payroll struct {
	Meta               model.Meta                         `json:"meta"`
	ID                 string                             `json:"id"`
	AccountID          string                             `json:"accountId"`
	Name               string                             `json:"name"`
	Description        string                             `json:"description"`
	ExternalCode       string                             `json:"externalCode"`
	SyncID             string                             `json:"syncId"`
	Created            model.DateTime                     `json:"created,omitzero"`
	Updated            model.DateTime                     `json:"updated,omitzero"`
	Moment             model.DateTime                     `json:"moment,omitzero"`
	StartPayrollPeriod model.DateTime                     `json:"startPayrollPeriod,omitzero"`
	EndPayrollPeriod   model.DateTime                     `json:"endPayrollPeriod,omitzero"`
	Applicable         bool                               `json:"applicable"`
	Shared             bool                               `json:"shared"`
	Printed            bool                               `json:"printed"`
	Published          bool                               `json:"published"`
	Organization       *model.Reference                   `json:"organization,omitzero"`
	Owner              *model.Reference                   `json:"owner,omitzero"`
	Group              *model.Reference                   `json:"group,omitzero"`
	State              *model.Reference                   `json:"state,omitzero"`
	Attributes         []model.Attribute                  `json:"attributes"`
	Sum                float64                            `json:"sum"`
	Positions          *model.Collection[Position]        `json:"positions,omitzero"`
	Files              *model.Collection[model.Reference] `json:"files,omitzero"`
	Payments           []model.Reference                  `json:"payments"`
}

// Position — строка начисления (meta.type payrollposition): сотрудник, начисление по окладу и
// сдельная оплата. Собственная форма, не positions.Position. BaseSalary и PieceworkSalary — в
// единицах sum документа (копейки): sum равна их сумме по строкам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#4-pozicii-nachisleniya-zarplat
type Position struct {
	Meta            model.Meta       `json:"meta"`
	ID              string           `json:"id"`
	AccountID       string           `json:"accountId"`
	Employee        *model.Reference `json:"employee,omitzero"`
	BaseSalary      float64          `json:"baseSalary"`
	PieceworkSalary float64          `json:"pieceworkSalary"`
}
