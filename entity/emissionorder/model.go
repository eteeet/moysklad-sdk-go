package emissionorder

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// EmissionOrder — заказ кодов маркировки. Модель чтения без expand. TrackingType, EmissionType и
// DocumentState — открытые enum: неизвестное значение сохраняется как есть. ProductionTasks —
// связанные производственные задания («Связи с другими документами»). Строки — GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder
type EmissionOrder struct {
	Meta            model.Meta                  `json:"meta"`
	ID              string                      `json:"id"`
	AccountID       string                      `json:"accountId"`
	Name            string                      `json:"name"`
	Description     string                      `json:"description"`
	ExternalCode    string                      `json:"externalCode"`
	Created         model.DateTime              `json:"created,omitzero"`
	Updated         model.DateTime              `json:"updated,omitzero"`
	Moment          model.DateTime              `json:"moment,omitzero"`
	Shared          bool                        `json:"shared"`
	Printed         bool                        `json:"printed"`
	Published       bool                        `json:"published"`
	Organization    *model.Reference            `json:"organization,omitzero"`
	Owner           *model.Reference            `json:"owner,omitzero"`
	Group           *model.Reference            `json:"group,omitzero"`
	State           *model.Reference            `json:"state,omitzero"`
	TrackingType    string                      `json:"trackingType"`
	EmissionType    string                      `json:"emissionType"`
	DocumentState   string                      `json:"documentState"`
	Positions       *model.Collection[Position] `json:"positions,omitzero"`
	ProductionTasks []model.Reference           `json:"productionTasks"`
}

// Position — строка заказа (meta.type emissionorderposition): assortment и quantity (не больше
// 1000 шт., проверяет сервер). Status — статус кодов, открытый enum, только для чтения.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#4-pozicii-zakaza-kodov-markirovki
type Position struct {
	positions.Item
	Status string `json:"status"`
}

// PositionRequest — строка в Create/Update документа, CreatePositions и UpdatePosition: id
// (только для сопоставления в Update документа), assortment и quantity.
type PositionRequest = positions.ItemRequest
