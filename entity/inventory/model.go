package inventory

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Position — строка инвентаризации (meta.type inventoryposition), собственная форма, а не
// positions.Position: таблица строк — assortment, quantity, price, pack, calculatedQuantity,
// correctionAmount, correctionSum. Price и CorrectionSum — копейки; CalculatedQuantity
// (расчётный остаток) и CorrectionAmount (разница с фактическим) — количества.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#4-pozicii-inventarizacii
type Position struct {
	positions.Item
	Price              float64         `json:"price"`
	Pack               *positions.Pack `json:"pack,omitzero"`
	CalculatedQuantity float64         `json:"calculatedQuantity"`
	CorrectionAmount   float64         `json:"correctionAmount"`
	CorrectionSum      float64         `json:"correctionSum"`
}

// PositionRequest — строка в embedded Create/Update и в отдельных операциях позиций.
// correctionAmount и correctionSum помечены «Только для чтения» и не пишутся; calculatedQuantity
// маркера не имеет и пишется (так же в примере Update). Для новой строки assortment и
// положительный quantity проверяет сервер; ID сопоставляет строку embedded Update.
type PositionRequest struct {
	positions.ItemRequest
	Price              model.Field[float64]               `json:"price,omitzero"`
	Pack               model.Field[positions.PackRequest] `json:"pack,omitzero"`
	CalculatedQuantity model.Field[float64]               `json:"calculatedQuantity,omitzero"`
}

// Inventory — инвентаризация. Модель чтения без expand: sum — копейки, SDK его не пересчитывает.
// У инвентаризации в таблице нет applicable, rate, project и полей НДС. Enters и Losses —
// связи с оприходованиями и списаниями (раздел «Связи с другими документами»).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory
type Inventory struct {
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
	Shared       bool                               `json:"shared"`
	Printed      bool                               `json:"printed"`
	Published    bool                               `json:"published"`
	Organization *model.Reference                   `json:"organization,omitzero"`
	Store        *model.Reference                   `json:"store,omitzero"`
	Owner        *model.Reference                   `json:"owner,omitzero"`
	Group        *model.Reference                   `json:"group,omitzero"`
	State        *model.Reference                   `json:"state,omitzero"`
	Attributes   []model.Attribute                  `json:"attributes"`
	Sum          float64                            `json:"sum"`
	Positions    *model.Collection[Position]        `json:"positions,omitzero"`
	Files        *model.Collection[model.Reference] `json:"files,omitzero"`
	Enters       []model.Reference                  `json:"enters"`
	Losses       []model.Reference                  `json:"losses"`
}
