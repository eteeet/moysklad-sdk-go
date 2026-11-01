package move

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки для всех товарных документов.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/move
type Move struct {
	ID              string                             `json:"id"`
	AccountID       string                             `json:"accountId"`
	Name            string                             `json:"name"`
	Description     string                             `json:"description"`
	Code            string                             `json:"code"`
	ExternalCode    string                             `json:"externalCode"`
	SyncID          string                             `json:"syncId"`
	Created         model.DateTime                     `json:"created,omitzero"`
	Updated         model.DateTime                     `json:"updated,omitzero"`
	Deleted         model.DateTime                     `json:"deleted,omitzero"`
	Moment          model.DateTime                     `json:"moment,omitzero"`
	Applicable      bool                               `json:"applicable"`
	Shared          bool                               `json:"shared"`
	Printed         bool                               `json:"printed"`
	Published       bool                               `json:"published"`
	Organization    *model.Reference                   `json:"organization,omitzero"`
	Owner           *model.Reference                   `json:"owner,omitzero"`
	Group           *model.Reference                   `json:"group,omitzero"`
	Project         *model.Reference                   `json:"project,omitzero"`
	State           *model.Reference                   `json:"state,omitzero"`
	Meta            model.Meta                         `json:"meta"`
	Attributes      []model.Attribute                  `json:"attributes"`
	Rate            *model.DocumentRate                `json:"rate,omitzero"`
	Sum             float64                            `json:"sum"`
	Positions       *model.Collection[Position]        `json:"positions,omitzero"`
	Files           *model.Collection[model.Reference] `json:"files,omitzero"`
	SourceStore     *model.Reference                   `json:"sourceStore,omitzero"`
	TargetStore     *model.Reference                   `json:"targetStore,omitzero"`
	InternalOrder   *model.Reference                   `json:"internalOrder,omitzero"`
	CustomerOrder   *model.Reference                   `json:"customerOrder,omitzero"`
	Demand          *model.Reference                   `json:"demand,omitzero"`
	Supply          *model.Reference                   `json:"supply,omitzero"`
	ProductionTasks []model.Reference                  `json:"productionTasks"`
	Overhead        *model.DocumentOverhead            `json:"overhead,omitzero"`
}
