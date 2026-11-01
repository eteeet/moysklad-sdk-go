package move

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/move
type CreateRequest struct {
	Name          model.Field[string]                        `json:"name,omitzero"`
	Description   model.Field[string]                        `json:"description,omitzero"`
	Code          model.Field[string]                        `json:"code,omitzero"`
	ExternalCode  model.Field[string]                        `json:"externalCode,omitzero"`
	SyncID        model.Field[string]                        `json:"syncId,omitzero"`
	Moment        model.Field[model.DateTime]                `json:"moment,omitzero"`
	Applicable    model.Field[bool]                          `json:"applicable,omitzero"`
	Shared        model.Field[bool]                          `json:"shared,omitzero"`
	Organization  model.Reference                            `json:"organization"`
	Owner         model.Field[model.Reference]               `json:"owner,omitzero"`
	Group         model.Field[model.Reference]               `json:"group,omitzero"`
	Project       model.Field[model.Reference]               `json:"project,omitzero"`
	State         model.Field[model.Reference]               `json:"state,omitzero"`
	Attributes    model.Field[[]model.Attribute]             `json:"attributes,omitzero"`
	Rate          model.Field[model.DocumentRateRequest]     `json:"rate,omitzero"`
	Positions     model.Field[[]PositionRequest]             `json:"positions,omitzero"`
	SourceStore   model.Reference                            `json:"sourceStore"`
	TargetStore   model.Reference                            `json:"targetStore"`
	InternalOrder model.Field[model.Reference]               `json:"internalOrder,omitzero"`
	CustomerOrder model.Field[model.Reference]               `json:"customerOrder,omitzero"`
	Overhead      model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/move
type UpdateRequest struct {
	Name          model.Field[string]                        `json:"name,omitzero"`
	Description   model.Field[string]                        `json:"description,omitzero"`
	Code          model.Field[string]                        `json:"code,omitzero"`
	ExternalCode  model.Field[string]                        `json:"externalCode,omitzero"`
	SyncID        model.Field[string]                        `json:"syncId,omitzero"`
	Moment        model.Field[model.DateTime]                `json:"moment,omitzero"`
	Applicable    model.Field[bool]                          `json:"applicable,omitzero"`
	Shared        model.Field[bool]                          `json:"shared,omitzero"`
	Organization  model.Field[model.Reference]               `json:"organization,omitzero"`
	Owner         model.Field[model.Reference]               `json:"owner,omitzero"`
	Group         model.Field[model.Reference]               `json:"group,omitzero"`
	Project       model.Field[model.Reference]               `json:"project,omitzero"`
	State         model.Field[model.Reference]               `json:"state,omitzero"`
	Attributes    model.Field[[]model.Attribute]             `json:"attributes,omitzero"`
	Rate          model.Field[model.DocumentRateRequest]     `json:"rate,omitzero"`
	Positions     model.Field[[]PositionRequest]             `json:"positions,omitzero"`
	SourceStore   model.Field[model.Reference]               `json:"sourceStore,omitzero"`
	TargetStore   model.Field[model.Reference]               `json:"targetStore,omitzero"`
	InternalOrder model.Field[model.Reference]               `json:"internalOrder,omitzero"`
	CustomerOrder model.Field[model.Reference]               `json:"customerOrder,omitzero"`
	Overhead      model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
}
