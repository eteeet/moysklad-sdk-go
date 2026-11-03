package consignment

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Consignment — партия. Name собирается сервером из наименования товара и метки и только читается;
// Images — метаданные коллекции изображений товара партии, только чтение.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment
type Consignment struct {
	Meta         model.Meta                        `json:"meta"`
	ID           string                            `json:"id"`
	AccountID    string                            `json:"accountId"`
	Updated      model.DateTime                    `json:"updated,omitzero"`
	Name         string                            `json:"name"`
	Label        string                            `json:"label"`
	Assortment   *model.Reference                  `json:"assortment,omitzero"`
	ExpiryDate   model.DateTime                    `json:"expiryDate,omitzero"`
	Description  string                            `json:"description"`
	Code         string                            `json:"code"`
	ExternalCode string                            `json:"externalCode"`
	Archived     bool                              `json:"archived"`
	Barcodes     []model.Barcode                   `json:"barcodes"`
	Attributes   []model.Attribute                 `json:"attributes"`
	Images       *model.Collection[jsontext.Value] `json:"images,omitzero"`
}
