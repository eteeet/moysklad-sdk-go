package processingorder

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// CreateRequest: таблица помечает «Необходимо при создании» name, organization, store,
// processingPlan и positions; раздел «Создать» и минимальный пример называют только organization,
// processingPlan и positions. SDK следует таблице: все пять проверяются до запроса. Positions —
// от 1 до 1000 строк. attributes, shared и syncId таблица помечает «Только для чтения» (примеры
// передают attributes и shared) — SDK их не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-sozdat-zakaz-na-proizvodstvo
type CreateRequest struct {
	Name                  string                         `json:"name"`
	Organization          model.Reference                `json:"organization"`
	Store                 model.Reference                `json:"store"`
	ProcessingPlan        model.Reference                `json:"processingPlan"`
	Positions             []PositionRequest              `json:"positions"`
	Description           model.Field[string]            `json:"description,omitzero"`
	Code                  model.Field[string]            `json:"code,omitzero"`
	ExternalCode          model.Field[string]            `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]    `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]    `json:"deliveryPlannedMoment,omitzero"`
	Applicable            model.Field[bool]              `json:"applicable,omitzero"`
	Quantity              model.Field[float64]           `json:"quantity,omitzero"`
	Owner                 model.Field[model.Reference]   `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]   `json:"group,omitzero"`
	Project               model.Field[model.Reference]   `json:"project,omitzero"`
	State                 model.Field[model.Reference]   `json:"state,omitzero"`
	Processings           model.Field[[]model.Reference] `json:"processings,omitzero"`
}

// UpdateRequest. Источник требует в каждом Update meta processingPlan, даже если техкарта не
// меняется, — Update проверяет его до запроса (UpdateMany — нет, проверяет сервер). Positions
// заменяют все строки: строка с ID — существующая, без — новая; позиции с новыми товарами в
// заказ добавить нельзя (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-izmenit-zakaz-na-proizvodstvo
type UpdateRequest struct {
	ProcessingPlan        model.Reference                `json:"processingPlan,omitzero"`
	Name                  model.Field[string]            `json:"name,omitzero"`
	Organization          model.Field[model.Reference]   `json:"organization,omitzero"`
	Store                 model.Field[model.Reference]   `json:"store,omitzero"`
	Positions             model.Field[[]PositionRequest] `json:"positions,omitzero"`
	Description           model.Field[string]            `json:"description,omitzero"`
	Code                  model.Field[string]            `json:"code,omitzero"`
	ExternalCode          model.Field[string]            `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]    `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]    `json:"deliveryPlannedMoment,omitzero"`
	Applicable            model.Field[bool]              `json:"applicable,omitzero"`
	Quantity              model.Field[float64]           `json:"quantity,omitzero"`
	Owner                 model.Field[model.Reference]   `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]   `json:"group,omitzero"`
	Project               model.Field[model.Reference]   `json:"project,omitzero"`
	State                 model.Field[model.Reference]   `json:"state,omitzero"`
	Processings           model.Field[[]model.Reference] `json:"processings,omitzero"`
}

// PositionRequest — строка в Create/Update документа и UpdatePosition. Для новой строки сервер
// требует assortment и quantity; в UpdatePosition обязательных полей нет.
type PositionRequest struct {
	// ID сопоставляет существующую строку в Update документа; для новой строки его не задают.
	ID         model.Field[string]                `json:"id,omitzero"`
	Assortment model.Field[model.Reference]       `json:"assortment,omitzero"`
	Pack       model.Field[positions.PackRequest] `json:"pack,omitzero"`
	Quantity   model.Field[float64]               `json:"quantity,omitzero"`
	Reserve    model.Field[float64]               `json:"reserve,omitzero"`
}
