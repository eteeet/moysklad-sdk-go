package productiontask

import "github.com/eteeet/moysklad-sdk-go/model"

// ProductionTask — производственное задание. Модель чтения без expand: денежных полей у документа
// нет. Позиции (техкарты) — GetProductionRows, продукты — GetProducts, производственные этапы —
// ListStages, материалы этапа — GetMaterials.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask
type ProductionTask struct {
	Meta                  model.Meta                         `json:"meta"`
	ID                    string                             `json:"id"`
	AccountID             string                             `json:"accountId"`
	Name                  string                             `json:"name"`
	Description           string                             `json:"description"`
	Code                  string                             `json:"code"`
	ExternalCode          string                             `json:"externalCode"`
	Created               model.DateTime                     `json:"created,omitzero"`
	Updated               model.DateTime                     `json:"updated,omitzero"`
	Deleted               model.DateTime                     `json:"deleted,omitzero"`
	Moment                model.DateTime                     `json:"moment,omitzero"`
	DeliveryPlannedMoment model.DateTime                     `json:"deliveryPlannedMoment,omitzero"`
	ProductionStart       model.DateTime                     `json:"productionStart,omitzero"`
	ProductionEnd         model.DateTime                     `json:"productionEnd,omitzero"`
	Applicable            bool                               `json:"applicable"`
	Awaiting              bool                               `json:"awaiting"`
	Reserve               bool                               `json:"reserve"`
	Shared                bool                               `json:"shared"`
	Printed               bool                               `json:"printed"`
	Published             bool                               `json:"published"`
	Organization          *model.Reference                   `json:"organization,omitzero"`
	MaterialsStore        *model.Reference                   `json:"materialsStore,omitzero"`
	ProductsStore         *model.Reference                   `json:"productsStore,omitzero"`
	Owner                 *model.Reference                   `json:"owner,omitzero"`
	Group                 *model.Reference                   `json:"group,omitzero"`
	State                 *model.Reference                   `json:"state,omitzero"`
	Attributes            []model.Attribute                  `json:"attributes"`
	ProductionRows        *model.Collection[ProductionRow]   `json:"productionRows,omitzero"`
	Products              *model.Collection[Product]         `json:"products,omitzero"`
	Files                 *model.Collection[model.Reference] `json:"files,omitzero"`
	// Связи с другими документами.
	CustomerOrders         []model.Reference `json:"customerOrders"`
	Demands                []model.Reference `json:"demands"`
	EmissionOrders         []model.Reference `json:"emissionOrders"`
	InternalOrders         []model.Reference `json:"internalOrders"`
	Moves                  []model.Reference `json:"moves"`
	ProductionTasks        []model.Reference `json:"productionTasks"`
	ProductionTaskSupplies []model.Reference `json:"productionTaskSupplies"`
	PurchaseOrders         []model.Reference `json:"purchaseOrders"`
	Supplies               []model.Reference `json:"supplies"`
}
