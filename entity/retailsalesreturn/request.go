package retailsalesreturn

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы name, organization, store, retailStore, retailShift и agent (agent метод
// называет обязательным, хотя ни один пример записи его не передаёт — SDK следует методу). Возврат
// по основанию задаёт Demand — ссылку на розничную продажу; без основания Demand не передают, но
// источник требует cashSum и noCashSum (Value(0) отправляется). Смена должна быть активной, moment —
// позже moment смены; контрагент, юрлицо и валюта совпадают с основанием, строки — его строки с
// количеством не больше исходного и без новых. Всё это проверяет сервер. Positions — не больше 1000.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retail-sales-return#3-sozdat-roznichnyj-vozvrat
type CreateRequest struct {
	Name                string                                 `json:"name"`
	Organization        model.Reference                        `json:"organization"`
	Store               model.Reference                        `json:"store"`
	RetailStore         model.Reference                        `json:"retailStore"`
	RetailShift         model.Reference                        `json:"retailShift"`
	Agent               model.Reference                        `json:"agent"`
	AgentAccount        model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Demand              model.Field[model.Reference]           `json:"demand,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID              model.Field[string]                    `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable          model.Field[bool]                      `json:"applicable,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	Contract            model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	TaxSystem           model.Field[string]                    `json:"taxSystem,omitzero"`
	VATEnabled          model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded         model.Field[bool]                      `json:"vatIncluded,omitzero"`
	CashSum             model.Field[float64]                   `json:"cashSum,omitzero"`
	NoCashSum           model.Field[float64]                   `json:"noCashSum,omitzero"`
	QRSum               model.Field[float64]                   `json:"qrSum,omitzero"`
	Positions           model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// UpdateRequest. Источник запрещает менять agent, agentAccount и demand — их здесь нет. При смене
// organization нужен и organizationAccount; валюту, отличную от основания, и строки не из основания
// сервер отклоняет. organizationBranch только для чтения; syncId после заполнения не меняется
// (проверяет сервер). Positions заменяют все строки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retail-sales-return#3-izmenit-roznichnyj-vozvrat
type UpdateRequest struct {
	Name                model.Field[string]                    `json:"name,omitzero"`
	Organization        model.Field[model.Reference]           `json:"organization,omitzero"`
	Store               model.Field[model.Reference]           `json:"store,omitzero"`
	RetailStore         model.Field[model.Reference]           `json:"retailStore,omitzero"`
	RetailShift         model.Field[model.Reference]           `json:"retailShift,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID              model.Field[string]                    `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable          model.Field[bool]                      `json:"applicable,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	Contract            model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	TaxSystem           model.Field[string]                    `json:"taxSystem,omitzero"`
	VATEnabled          model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded         model.Field[bool]                      `json:"vatIncluded,omitzero"`
	CashSum             model.Field[float64]                   `json:"cashSum,omitzero"`
	NoCashSum           model.Field[float64]                   `json:"noCashSum,omitzero"`
	QRSum               model.Field[float64]                   `json:"qrSum,omitzero"`
	Positions           model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}
