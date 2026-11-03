package demand

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки для всех товарных документов.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/demand
type Demand struct {
	ID                  string                             `json:"id"`
	AccountID           string                             `json:"accountId"`
	Name                string                             `json:"name"`
	Description         string                             `json:"description"`
	Code                string                             `json:"code"`
	ExternalCode        string                             `json:"externalCode"`
	SyncID              string                             `json:"syncId"`
	Created             model.DateTime                     `json:"created,omitzero"`
	Updated             model.DateTime                     `json:"updated,omitzero"`
	Deleted             model.DateTime                     `json:"deleted,omitzero"`
	Moment              model.DateTime                     `json:"moment,omitzero"`
	Applicable          bool                               `json:"applicable"`
	Shared              bool                               `json:"shared"`
	Printed             bool                               `json:"printed"`
	Published           bool                               `json:"published"`
	Organization        *model.Reference                   `json:"organization,omitzero"`
	Owner               *model.Reference                   `json:"owner,omitzero"`
	Group               *model.Reference                   `json:"group,omitzero"`
	Project             *model.Reference                   `json:"project,omitzero"`
	State               *model.Reference                   `json:"state,omitzero"`
	Meta                model.Meta                         `json:"meta"`
	Attributes          []model.Attribute                  `json:"attributes"`
	Rate                *model.DocumentRate                `json:"rate,omitzero"`
	Sum                 float64                            `json:"sum"`
	Agent               *model.Reference                   `json:"agent,omitzero"`
	OrganizationAccount *model.Reference                   `json:"organizationAccount,omitzero"`
	AgentAccount        *model.Reference                   `json:"agentAccount,omitzero"`
	Contract            *model.Reference                   `json:"contract,omitzero"`
	VATSum              float64                            `json:"vatSum"`
	Positions           *model.Collection[Position]        `json:"positions,omitzero"`
	VATEnabled          bool                               `json:"vatEnabled"`
	VATIncluded         bool                               `json:"vatIncluded"`
	PayedSum            float64                            `json:"payedSum"`
	Store               *model.Reference                   `json:"store,omitzero"`
	Payments            []model.Reference                  `json:"payments"`
	Files               *model.Collection[model.Reference] `json:"files,omitzero"`
	SalesChannel        *model.Reference                   `json:"salesChannel,omitzero"`
	CustomerOrder       *model.Reference                   `json:"customerOrder,omitzero"`
	FactureOut          *model.Reference                   `json:"factureOut,omitzero"`
	ShipmentAddress     string                             `json:"shipmentAddress"`
	ShipmentAddressFull *model.Address                     `json:"shipmentAddressFull,omitzero"`
	Returns             []model.Reference                  `json:"returns"`
	ProductionTasks     []model.Reference                  `json:"productionTasks"`
	InvoicesOut         []model.Reference                  `json:"invoicesOut"`
	Overhead            *model.DocumentOverhead            `json:"overhead,omitzero"`

	// Транспортные поля («Другие поля» источника): carrier, consignee и consignor — контрагент или
	// юрлицо, *Branch — подразделение юрлица; GoodPackQuantity — всего мест.
	CargoName               string           `json:"cargoName"`
	Carrier                 *model.Reference `json:"carrier,omitzero"`
	Consignee               *model.Reference `json:"consignee,omitzero"`
	ConsigneeBranch         *model.Reference `json:"consigneeBranch,omitzero"`
	Consignor               *model.Reference `json:"consignor,omitzero"`
	ConsignorBranch         *model.Reference `json:"consignorBranch,omitzero"`
	GoodPackQuantity        int              `json:"goodPackQuantity"`
	ShippingInstructions    string           `json:"shippingInstructions"`
	StateContractID         string           `json:"stateContractId"`
	TransportFacility       string           `json:"transportFacility"`
	TransportFacilityNumber string           `json:"transportFacilityNumber"`
	ShippingDocumentNumber  string           `json:"shippingDocumentNumber"`
	ShippingDocumentDate    model.DateTime   `json:"shippingDocumentDate,omitzero"`
}
