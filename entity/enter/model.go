package enter

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов; у оприходования — reason, slot, pack, country, gtd,
// things, overhead (readonly).
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// Enter — оприходование. Модель чтения без expand: sum и overhead.sum — копейки; overhead
// распределяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/enter
type Enter struct {
	Meta         model.Meta                         `json:"meta"`
	ID           string                             `json:"id"`
	AccountID    string                             `json:"accountId"`
	Name         string                             `json:"name"`
	Description  string                             `json:"description"`
	Code         string                             `json:"code"`
	ExternalCode string                             `json:"externalCode"`
	SyncID       string                             `json:"syncId"`
	Created      model.DateTime                     `json:"created,omitzero"`
	Updated      model.DateTime                     `json:"updated,omitzero"`
	Deleted      model.DateTime                     `json:"deleted,omitzero"`
	Moment       model.DateTime                     `json:"moment,omitzero"`
	Applicable   bool                               `json:"applicable"`
	Shared       bool                               `json:"shared"`
	Printed      bool                               `json:"printed"`
	Published    bool                               `json:"published"`
	Organization *model.Reference                   `json:"organization,omitzero"`
	Store        *model.Reference                   `json:"store,omitzero"`
	Owner        *model.Reference                   `json:"owner,omitzero"`
	Group        *model.Reference                   `json:"group,omitzero"`
	Project      *model.Reference                   `json:"project,omitzero"`
	State        *model.Reference                   `json:"state,omitzero"`
	Attributes   []model.Attribute                  `json:"attributes"`
	Rate         *model.DocumentRate                `json:"rate,omitzero"`
	Sum          float64                            `json:"sum"`
	Positions    *model.Collection[Position]        `json:"positions,omitzero"`
	Files        *model.Collection[model.Reference] `json:"files,omitzero"`
	Overhead     *model.DocumentOverhead            `json:"overhead,omitzero"`
	// Inventory — связанная инвентаризация; пишется как связь документа.
	Inventory *model.Reference `json:"inventory,omitzero"`
}
