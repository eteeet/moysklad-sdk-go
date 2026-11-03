package variant

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Variant — модификация товара. Модель чтения без expand: Product — только ссылка, товар
// не загружается. Деньги в ценах — копейки. Если у модификации нет своих цен продажи, сервер
// отдаёт цены товара. minimumStock отложен, как у Product: поле не читается и не пишется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant
type Variant struct {
	Meta               model.Meta     `json:"meta"`
	ID                 string         `json:"id"`
	AccountID          string         `json:"accountId"`
	Updated            model.DateTime `json:"updated,omitzero"`
	Name               string         `json:"name"`
	Code               string         `json:"code"`
	ExternalCode       string         `json:"externalCode"`
	Description        string         `json:"description"`
	Article            string         `json:"article"`
	Archived           bool           `json:"archived"`
	DiscountProhibited bool           `json:"discountProhibited"`

	Product         model.Reference  `json:"product"`
	Characteristics []Characteristic `json:"characteristics"`

	BuyPrice   *model.Price      `json:"buyPrice,omitzero"`
	MinPrice   *model.Price      `json:"minPrice,omitzero"`
	SalePrices []model.SalePrice `json:"salePrices"`
	Barcodes   []model.Barcode   `json:"barcodes"`
	// MinimumStock приходит только при List с fields=minimumStock.
	MinimumStock *model.MinimumStock `json:"minimumStock,omitzero"`
	Packs        []Pack              `json:"packs"`
	// Things — серийные номера, только для чтения.
	Things []string `json:"things"`
	// Images — только метаданные коллекции; изображения читаются отдельным ресурсом.
	Images *model.Collection[jsontext.Value] `json:"images,omitzero"`
}

// Pack — упаковка модификации: переопределяет штрихкод упаковки товара ParentPack.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#5-upakovki-modifikacii
type Pack struct {
	ID         string          `json:"id"`
	Barcodes   []model.Barcode `json:"barcodes"`
	ParentPack model.Reference `json:"parentpack"`
}
