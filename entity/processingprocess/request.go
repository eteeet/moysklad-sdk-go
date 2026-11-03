package processingprocess

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязательны непустые name и positions (1–100 позиций, у каждой — processingstage);
// проверяется до запроса. Порядок позиций в теле — порядок среза.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-sozdat-tehprocess
type CreateRequest struct {
	Name         string                       `json:"name"`
	Positions    []PositionRequest            `json:"positions"`
	Description  model.Field[string]          `json:"description,omitzero"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Archived     model.Field[bool]            `json:"archived,omitzero"`
	Shared       model.Field[bool]            `json:"shared,omitzero"`
	Owner        model.Field[model.Reference] `json:"owner,omitzero"`
	Group        model.Field[model.Reference] `json:"group,omitzero"`
}

// UpdateRequest. Positions (1–100, проверяется до запроса) передаёт набор позиций в порядке среза:
// строка с Meta изменяет существующую позицию (смена этапа её не пересоздаёт), строка без Meta —
// новая позиция. Что сервер делает с позициями, не переданными в массиве, источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-izmenit-tehprocess
type UpdateRequest struct {
	Name         model.Field[string]                  `json:"name,omitzero"`
	Positions    model.Field[[]PositionUpdateRequest] `json:"positions,omitzero"`
	Description  model.Field[string]                  `json:"description,omitzero"`
	ExternalCode model.Field[string]                  `json:"externalCode,omitzero"`
	Archived     model.Field[bool]                    `json:"archived,omitzero"`
	Shared       model.Field[bool]                    `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]         `json:"owner,omitzero"`
	Group        model.Field[model.Reference]         `json:"group,omitzero"`
}

// PositionRequest — новая позиция в CreateRequest и CreatePositions: processingstage обязателен
// (непустой meta.href проверяется до запроса).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-sozdat-pozicii-tehprocessa
type PositionRequest struct {
	ProcessingStage model.Reference                `json:"processingstage"`
	NextPositions   model.Field[[]model.Reference] `json:"nextPositions,omitzero"`
}

// PositionUpdateRequest — тело UpdatePosition и строка UpdateRequest.Positions. Meta (meta позиции
// из Position.Meta) адресует существующую позицию только в UpdateRequest.Positions, как в примере
// источника; UpdatePosition с Meta отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-izmenit-poziciyu-tehprocessa
type PositionUpdateRequest struct {
	Meta            model.Field[model.Meta]        `json:"meta,omitzero"`
	ProcessingStage model.Field[model.Reference]   `json:"processingstage,omitzero"`
	NextPositions   model.Field[[]model.Reference] `json:"nextPositions,omitzero"`
}
