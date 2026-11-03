package internalorder

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов (meta.type internalorderposition): таблица строк
// источника — assortment, quantity, price, vat, vatEnabled, pack; остальные поля общей формы
// остаются нулевыми.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// InternalOrder — внутренний заказ. Модель чтения без expand: sum и vatSum — копейки; суммы и
// НДС SDK не пересчитывает. Moves, PurchaseOrders и ProductionTasks — связи с другими
// документами (productionTasks есть только в разделе связей, не в основной таблице).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder
type InternalOrder struct {
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
	Store                 *model.Reference                   `json:"store,omitzero"`
	Owner                 *model.Reference                   `json:"owner,omitzero"`
	Group                 *model.Reference                   `json:"group,omitzero"`
	Project               *model.Reference                   `json:"project,omitzero"`
	State                 *model.Reference                   `json:"state,omitzero"`
	Attributes            []model.Attribute                  `json:"attributes"`
	Rate                  *model.DocumentRate                `json:"rate,omitzero"`
	Sum                   float64                            `json:"sum"`
	VATEnabled            bool                               `json:"vatEnabled"`
	VATIncluded           bool                               `json:"vatIncluded"`
	VATSum                float64                            `json:"vatSum"`
	Positions             *model.Collection[Position]        `json:"positions,omitzero"`
	Files                 *model.Collection[model.Reference] `json:"files,omitzero"`
	Moves                 []model.Reference                  `json:"moves"`
	PurchaseOrders        []model.Reference                  `json:"purchaseOrders"`
	ProductionTasks       []model.Reference                  `json:"productionTasks"`
}
