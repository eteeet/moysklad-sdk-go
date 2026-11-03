package customerorder

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов; у заказа — reserve, shipped (readonly), pack, taxSystem.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// CustomerOrder — заказ покупателя. Модель чтения без expand: суммы — копейки и только для
// чтения; связанные документы (demands, payments …) — ссылки, пишутся в Create/Update; строки
// позиций — GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/customerOrder
type CustomerOrder struct {
	Meta                  model.Meta                         `json:"meta"`
	ID                    string                             `json:"id"`
	AccountID             string                             `json:"accountId"`
	Name                  string                             `json:"name"`
	Description           string                             `json:"description"`
	Code                  string                             `json:"code"`
	ExternalCode          string                             `json:"externalCode"`
	SyncID                string                             `json:"syncId"`
	Created               model.DateTime                     `json:"created,omitzero"`
	Updated               model.DateTime                     `json:"updated,omitzero"`
	Deleted               model.DateTime                     `json:"deleted,omitzero"`
	Moment                model.DateTime                     `json:"moment,omitzero"`
	DeliveryPlannedMoment model.DateTime                     `json:"deliveryPlannedMoment,omitzero"`
	Applicable            bool                               `json:"applicable"`
	Shared                bool                               `json:"shared"`
	Printed               bool                               `json:"printed"`
	Published             bool                               `json:"published"`
	Organization          *model.Reference                   `json:"organization,omitzero"`
	OrganizationAccount   *model.Reference                   `json:"organizationAccount,omitzero"`
	Agent                 *model.Reference                   `json:"agent,omitzero"`
	AgentAccount          *model.Reference                   `json:"agentAccount,omitzero"`
	Contract              *model.Reference                   `json:"contract,omitzero"`
	Owner                 *model.Reference                   `json:"owner,omitzero"`
	Group                 *model.Reference                   `json:"group,omitzero"`
	Project               *model.Reference                   `json:"project,omitzero"`
	State                 *model.Reference                   `json:"state,omitzero"`
	Store                 *model.Reference                   `json:"store,omitzero"`
	SalesChannel          *model.Reference                   `json:"salesChannel,omitzero"`
	Attributes            []model.Attribute                  `json:"attributes"`
	Rate                  *model.DocumentRate                `json:"rate,omitzero"`
	TaxSystem             string                             `json:"taxSystem"`
	VATEnabled            bool                               `json:"vatEnabled"`
	VATIncluded           bool                               `json:"vatIncluded"`
	Sum                   float64                            `json:"sum"`
	VATSum                float64                            `json:"vatSum"`
	PayedSum              float64                            `json:"payedSum"`
	ShippedSum            float64                            `json:"shippedSum"`
	InvoicedSum           float64                            `json:"invoicedSum"`
	ReservedSum           float64                            `json:"reservedSum"`
	ShipmentAddress       string                             `json:"shipmentAddress"`
	ShipmentAddressFull   *model.Address                     `json:"shipmentAddressFull,omitzero"`
	Positions             *model.Collection[Position]        `json:"positions,omitzero"`
	Files                 *model.Collection[model.Reference] `json:"files,omitzero"`

	PurchaseOrders  []model.Reference `json:"purchaseOrders"`
	Demands         []model.Reference `json:"demands"`
	Payments        []model.Reference `json:"payments"`
	ProductionTasks []model.Reference `json:"productionTasks"`
	InvoicesOut     []model.Reference `json:"invoicesOut"`
	Moves           []model.Reference `json:"moves"`
	Prepayments     []model.Reference `json:"prepayments"`
}
