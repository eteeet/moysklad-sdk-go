package positions

import "github.com/eteeet/moysklad-sdk-go/model"

// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L62-L89
type MovementRequest struct {
	ItemRequest
	Price      model.Field[float64]         `json:"price,omitzero"`
	SourceSlot model.Field[model.Reference] `json:"sourceSlot,omitzero"`
	TargetSlot model.Field[model.Reference] `json:"targetSlot,omitzero"`
}
