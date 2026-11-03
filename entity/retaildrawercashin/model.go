package retaildrawercashin

import "github.com/eteeet/moysklad-sdk-go/model"

// RetailDrawerCashIn — внесение денег в кассу розничной смены. Модель чтения без expand: sum —
// копейки, только для чтения.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashin
type RetailDrawerCashIn struct {
	Meta               model.Meta                         `json:"meta"`
	ID                 string                             `json:"id"`
	AccountID          string                             `json:"accountId"`
	Name               string                             `json:"name"`
	Description        string                             `json:"description"`
	ExternalCode       string                             `json:"externalCode"`
	SyncID             string                             `json:"syncId"`
	Created            model.DateTime                     `json:"created,omitzero"`
	Updated            model.DateTime                     `json:"updated,omitzero"`
	Deleted            model.DateTime                     `json:"deleted,omitzero"`
	Moment             model.DateTime                     `json:"moment,omitzero"`
	Applicable         bool                               `json:"applicable"`
	Shared             bool                               `json:"shared"`
	Printed            bool                               `json:"printed"`
	Published          bool                               `json:"published"`
	Organization       *model.Reference                   `json:"organization,omitzero"`
	OrganizationBranch *model.Reference                   `json:"organizationBranch,omitzero"`
	Agent              *model.Reference                   `json:"agent,omitzero"`
	Owner              *model.Reference                   `json:"owner,omitzero"`
	Group              *model.Reference                   `json:"group,omitzero"`
	State              *model.Reference                   `json:"state,omitzero"`
	Attributes         []model.Attribute                  `json:"attributes"`
	Rate               *model.DocumentRate                `json:"rate,omitzero"`
	Sum                float64                            `json:"sum"`
	Files              *model.Collection[model.Reference] `json:"files,omitzero"`
	// RetailShift — розничная смена (раздел «Связи с другими документами»).
	RetailShift *model.Reference `json:"retailShift,omitzero"`
}
