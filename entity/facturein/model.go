package facturein

import "github.com/eteeet/moysklad-sdk-go/model"

// FactureIn — счёт-фактура полученный. Модель чтения без expand и без позиций: sum — копейки,
// readonly. Supplies (приёмки) и Payments (исходящие платежи) — связи с документами-основаниями.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/facturein
type FactureIn struct {
	Meta           model.Meta                         `json:"meta"`
	ID             string                             `json:"id"`
	AccountID      string                             `json:"accountId"`
	Name           string                             `json:"name"`
	Description    string                             `json:"description"`
	Code           string                             `json:"code"`
	ExternalCode   string                             `json:"externalCode"`
	SyncID         string                             `json:"syncId"`
	Created        model.DateTime                     `json:"created,omitzero"`
	Updated        model.DateTime                     `json:"updated,omitzero"`
	Deleted        model.DateTime                     `json:"deleted,omitzero"`
	Moment         model.DateTime                     `json:"moment,omitzero"`
	Applicable     bool                               `json:"applicable"`
	Shared         bool                               `json:"shared"`
	Printed        bool                               `json:"printed"`
	Published      bool                               `json:"published"`
	Organization   *model.Reference                   `json:"organization,omitzero"`
	Agent          *model.Reference                   `json:"agent,omitzero"`
	Contract       *model.Reference                   `json:"contract,omitzero"`
	Owner          *model.Reference                   `json:"owner,omitzero"`
	Group          *model.Reference                   `json:"group,omitzero"`
	State          *model.Reference                   `json:"state,omitzero"`
	Attributes     []model.Attribute                  `json:"attributes"`
	Rate           *model.DocumentRate                `json:"rate,omitzero"`
	Sum            float64                            `json:"sum"`
	Files          *model.Collection[model.Reference] `json:"files,omitzero"`
	Supplies       []model.Reference                  `json:"supplies"`
	Payments       []model.Reference                  `json:"payments"`
	IncomingNumber string                             `json:"incomingNumber"`
	IncomingDate   model.DateTime                     `json:"incomingDate,omitzero"`
}
