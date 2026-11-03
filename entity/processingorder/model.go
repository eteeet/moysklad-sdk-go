package processingorder

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// ProcessingOrder — заказ на производство. Модель чтения без expand: денежных полей в таблице нет;
// quantity — объём производства, не деньги. sum приходит в примерах ответов, но таблица его не
// называет — в модель не входит. Строки — GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#2-zakaz-na-proizvodstvo
type ProcessingOrder struct {
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
	ProcessingPlan        *model.Reference                   `json:"processingPlan,omitzero"`
	Owner                 *model.Reference                   `json:"owner,omitzero"`
	Group                 *model.Reference                   `json:"group,omitzero"`
	Project               *model.Reference                   `json:"project,omitzero"`
	State                 *model.Reference                   `json:"state,omitzero"`
	Attributes            []model.Attribute                  `json:"attributes"`
	Quantity              float64                            `json:"quantity"`
	Positions             *model.Collection[Position]        `json:"positions,omitzero"`
	Files                 *model.Collection[model.Reference] `json:"files,omitzero"`
	// Processings — связанные техоперации (раздел «Связи с другими документами»).
	Processings []model.Reference `json:"processings"`
}

// Position — строка заказа (meta.type processingorderposition): товар или модификация,
// соответствующие материалам техкарты. Собственная форма, не positions.Position: цен, НДС и скидок
// у строки нет. Quantity и Reserve — количества, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#4-pozicii-zakaza-na-proizvodstvo
type Position struct {
	Meta       model.Meta       `json:"meta"`
	ID         string           `json:"id"`
	AccountID  string           `json:"accountId"`
	Assortment *model.Reference `json:"assortment,omitzero"`
	Pack       *positions.Pack  `json:"pack,omitzero"`
	Quantity   float64          `json:"quantity"`
	Reserve    float64          `json:"reserve"`
}
