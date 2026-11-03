package retireorder

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, retireOrderType и trackingType (пустые отклоняются до
// запроса); значения enum не проверяются. Допустимость комбинаций trackingType/retireOrderType/
// supportingTransaction/destinationCountry и условно обязательные поля (supportingTransaction*,
// stateContractId, reasonDescription, primaryDocumentName) проверяет сервер. Sum в таблице не
// помечен «Только для чтения», поэтому пишется. Positions — не больше 1000 строк.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-sozdat-vyvod-iz-oborota
type CreateRequest struct {
	Organization                model.Reference                        `json:"organization"`
	RetireOrderType             string                                 `json:"retireOrderType"`
	TrackingType                string                                 `json:"trackingType"`
	Name                        model.Field[string]                    `json:"name,omitzero"`
	Description                 model.Field[string]                    `json:"description,omitzero"`
	ExternalCode                model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID                      model.Field[string]                    `json:"syncId,omitzero"`
	Moment                      model.Field[model.DateTime]            `json:"moment,omitzero"`
	Shared                      model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationBranch          model.Field[model.Reference]           `json:"organizationBranch,omitzero"`
	Agent                       model.Field[model.Reference]           `json:"agent,omitzero"`
	DestinationCountry          model.Field[model.Reference]           `json:"destinationCountry,omitzero"`
	Owner                       model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                       model.Field[model.Reference]           `json:"group,omitzero"`
	State                       model.Field[model.Reference]           `json:"state,omitzero"`
	Rate                        model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Sum                         model.Field[float64]                   `json:"sum,omitzero"`
	VATEnabled                  model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded                 model.Field[bool]                      `json:"vatIncluded,omitzero"`
	ReasonDescription           model.Field[string]                    `json:"reasonDescription,omitzero"`
	PrimaryDocumentName         model.Field[string]                    `json:"primaryDocumentName,omitzero"`
	StateContractID             model.Field[string]                    `json:"stateContractId,omitzero"`
	SupportingTransaction       model.Field[string]                    `json:"supportingTransaction,omitzero"`
	SupportingTransactionDate   model.Field[model.DateTime]            `json:"supportingTransactionDate,omitzero"`
	SupportingTransactionNumber model.Field[string]                    `json:"supportingTransactionNumber,omitzero"`
	Positions                   model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// UpdateRequest. Изменение допустимо только при documentState CREATED, CHECKED_NOT_OK или
// PROCESSING_ERROR — проверяет сервер. Positions заменяют все строки: строка с ID —
// существующая, без — новая. syncId после заполнения не меняется (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-izmenenie-vyvoda-iz-oborota
type UpdateRequest struct {
	Organization                model.Field[model.Reference]           `json:"organization,omitzero"`
	RetireOrderType             model.Field[string]                    `json:"retireOrderType,omitzero"`
	TrackingType                model.Field[string]                    `json:"trackingType,omitzero"`
	Name                        model.Field[string]                    `json:"name,omitzero"`
	Description                 model.Field[string]                    `json:"description,omitzero"`
	ExternalCode                model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID                      model.Field[string]                    `json:"syncId,omitzero"`
	Moment                      model.Field[model.DateTime]            `json:"moment,omitzero"`
	Shared                      model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationBranch          model.Field[model.Reference]           `json:"organizationBranch,omitzero"`
	Agent                       model.Field[model.Reference]           `json:"agent,omitzero"`
	DestinationCountry          model.Field[model.Reference]           `json:"destinationCountry,omitzero"`
	Owner                       model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                       model.Field[model.Reference]           `json:"group,omitzero"`
	State                       model.Field[model.Reference]           `json:"state,omitzero"`
	Rate                        model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Sum                         model.Field[float64]                   `json:"sum,omitzero"`
	VATEnabled                  model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded                 model.Field[bool]                      `json:"vatIncluded,omitzero"`
	ReasonDescription           model.Field[string]                    `json:"reasonDescription,omitzero"`
	PrimaryDocumentName         model.Field[string]                    `json:"primaryDocumentName,omitzero"`
	StateContractID             model.Field[string]                    `json:"stateContractId,omitzero"`
	SupportingTransaction       model.Field[string]                    `json:"supportingTransaction,omitzero"`
	SupportingTransactionDate   model.Field[model.DateTime]            `json:"supportingTransactionDate,omitzero"`
	SupportingTransactionNumber model.Field[string]                    `json:"supportingTransactionNumber,omitzero"`
	Positions                   model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}
