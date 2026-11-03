package demand

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/demand
type CreateRequest struct {
	Name                    model.Field[string]                        `json:"name,omitzero"`
	Description             model.Field[string]                        `json:"description,omitzero"`
	Code                    model.Field[string]                        `json:"code,omitzero"`
	ExternalCode            model.Field[string]                        `json:"externalCode,omitzero"`
	SyncID                  model.Field[string]                        `json:"syncId,omitzero"`
	Moment                  model.Field[model.DateTime]                `json:"moment,omitzero"`
	Applicable              model.Field[bool]                          `json:"applicable,omitzero"`
	Shared                  model.Field[bool]                          `json:"shared,omitzero"`
	Organization            model.Reference                            `json:"organization"`
	Owner                   model.Field[model.Reference]               `json:"owner,omitzero"`
	Group                   model.Field[model.Reference]               `json:"group,omitzero"`
	Project                 model.Field[model.Reference]               `json:"project,omitzero"`
	State                   model.Field[model.Reference]               `json:"state,omitzero"`
	Attributes              model.Field[[]model.Attribute]             `json:"attributes,omitzero"`
	Rate                    model.Field[model.DocumentRateRequest]     `json:"rate,omitzero"`
	Agent                   model.Reference                            `json:"agent"`
	OrganizationAccount     model.Field[model.Reference]               `json:"organizationAccount,omitzero"`
	AgentAccount            model.Field[model.Reference]               `json:"agentAccount,omitzero"`
	Contract                model.Field[model.Reference]               `json:"contract,omitzero"`
	Positions               model.Field[[]PositionRequest]             `json:"positions,omitzero"`
	VATEnabled              model.Field[bool]                          `json:"vatEnabled,omitzero"`
	VATIncluded             model.Field[bool]                          `json:"vatIncluded,omitzero"`
	Store                   model.Field[model.Reference]               `json:"store,omitzero"`
	SalesChannel            model.Field[model.Reference]               `json:"salesChannel,omitzero"`
	CustomerOrder           model.Field[model.Reference]               `json:"customerOrder,omitzero"`
	ShipmentAddress         model.Field[string]                        `json:"shipmentAddress,omitzero"`
	ShipmentAddressFull     model.Field[model.AddressRequest]          `json:"shipmentAddressFull,omitzero"`
	Overhead                model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
	FactureOut              model.Field[model.Reference]               `json:"factureOut,omitzero"`
	Returns                 model.Field[[]model.Reference]             `json:"returns,omitzero"`
	Payments                model.Field[[]model.Reference]             `json:"payments,omitzero"`
	ProductionTasks         model.Field[[]model.Reference]             `json:"productionTasks,omitzero"`
	InvoicesOut             model.Field[[]model.Reference]             `json:"invoicesOut,omitzero"`
	CargoName               model.Field[string]                        `json:"cargoName,omitzero"`
	Carrier                 model.Field[model.Reference]               `json:"carrier,omitzero"`
	Consignee               model.Field[model.Reference]               `json:"consignee,omitzero"`
	ConsigneeBranch         model.Field[model.Reference]               `json:"consigneeBranch,omitzero"`
	Consignor               model.Field[model.Reference]               `json:"consignor,omitzero"`
	ConsignorBranch         model.Field[model.Reference]               `json:"consignorBranch,omitzero"`
	GoodPackQuantity        model.Field[int]                           `json:"goodPackQuantity,omitzero"`
	ShippingInstructions    model.Field[string]                        `json:"shippingInstructions,omitzero"`
	StateContractID         model.Field[string]                        `json:"stateContractId,omitzero"`
	TransportFacility       model.Field[string]                        `json:"transportFacility,omitzero"`
	TransportFacilityNumber model.Field[string]                        `json:"transportFacilityNumber,omitzero"`
	ShippingDocumentNumber  model.Field[string]                        `json:"shippingDocumentNumber,omitzero"`
	ShippingDocumentDate    model.Field[model.DateTime]                `json:"shippingDocumentDate,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/demand
type UpdateRequest struct {
	Name                    model.Field[string]                        `json:"name,omitzero"`
	Description             model.Field[string]                        `json:"description,omitzero"`
	Code                    model.Field[string]                        `json:"code,omitzero"`
	ExternalCode            model.Field[string]                        `json:"externalCode,omitzero"`
	SyncID                  model.Field[string]                        `json:"syncId,omitzero"`
	Moment                  model.Field[model.DateTime]                `json:"moment,omitzero"`
	Applicable              model.Field[bool]                          `json:"applicable,omitzero"`
	Shared                  model.Field[bool]                          `json:"shared,omitzero"`
	Organization            model.Field[model.Reference]               `json:"organization,omitzero"`
	Owner                   model.Field[model.Reference]               `json:"owner,omitzero"`
	Group                   model.Field[model.Reference]               `json:"group,omitzero"`
	Project                 model.Field[model.Reference]               `json:"project,omitzero"`
	State                   model.Field[model.Reference]               `json:"state,omitzero"`
	Attributes              model.Field[[]model.Attribute]             `json:"attributes,omitzero"`
	Rate                    model.Field[model.DocumentRateRequest]     `json:"rate,omitzero"`
	Agent                   model.Field[model.Reference]               `json:"agent,omitzero"`
	OrganizationAccount     model.Field[model.Reference]               `json:"organizationAccount,omitzero"`
	AgentAccount            model.Field[model.Reference]               `json:"agentAccount,omitzero"`
	Contract                model.Field[model.Reference]               `json:"contract,omitzero"`
	Positions               model.Field[[]PositionRequest]             `json:"positions,omitzero"`
	VATEnabled              model.Field[bool]                          `json:"vatEnabled,omitzero"`
	VATIncluded             model.Field[bool]                          `json:"vatIncluded,omitzero"`
	Store                   model.Field[model.Reference]               `json:"store,omitzero"`
	SalesChannel            model.Field[model.Reference]               `json:"salesChannel,omitzero"`
	CustomerOrder           model.Field[model.Reference]               `json:"customerOrder,omitzero"`
	ShipmentAddress         model.Field[string]                        `json:"shipmentAddress,omitzero"`
	ShipmentAddressFull     model.Field[model.AddressRequest]          `json:"shipmentAddressFull,omitzero"`
	Overhead                model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
	FactureOut              model.Field[model.Reference]               `json:"factureOut,omitzero"`
	Returns                 model.Field[[]model.Reference]             `json:"returns,omitzero"`
	Payments                model.Field[[]model.Reference]             `json:"payments,omitzero"`
	ProductionTasks         model.Field[[]model.Reference]             `json:"productionTasks,omitzero"`
	InvoicesOut             model.Field[[]model.Reference]             `json:"invoicesOut,omitzero"`
	CargoName               model.Field[string]                        `json:"cargoName,omitzero"`
	Carrier                 model.Field[model.Reference]               `json:"carrier,omitzero"`
	Consignee               model.Field[model.Reference]               `json:"consignee,omitzero"`
	ConsigneeBranch         model.Field[model.Reference]               `json:"consigneeBranch,omitzero"`
	Consignor               model.Field[model.Reference]               `json:"consignor,omitzero"`
	ConsignorBranch         model.Field[model.Reference]               `json:"consignorBranch,omitzero"`
	GoodPackQuantity        model.Field[int]                           `json:"goodPackQuantity,omitzero"`
	ShippingInstructions    model.Field[string]                        `json:"shippingInstructions,omitzero"`
	StateContractID         model.Field[string]                        `json:"stateContractId,omitzero"`
	TransportFacility       model.Field[string]                        `json:"transportFacility,omitzero"`
	TransportFacilityNumber model.Field[string]                        `json:"transportFacilityNumber,omitzero"`
	ShippingDocumentNumber  model.Field[string]                        `json:"shippingDocumentNumber,omitzero"`
	ShippingDocumentDate    model.Field[model.DateTime]                `json:"shippingDocumentDate,omitzero"`
}
