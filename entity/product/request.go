package product

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-sozdat-tovar

import "github.com/eteeet/moysklad-sdk-go/model"

type CreateRequest struct {
	Name               string               `json:"name"`
	Code               model.Field[string]  `json:"code,omitzero"`
	ExternalCode       model.Field[string]  `json:"externalCode,omitzero"`
	Description        model.Field[string]  `json:"description,omitzero"`
	Article            model.Field[string]  `json:"article,omitzero"`
	Archived           model.Field[bool]    `json:"archived,omitzero"`
	Shared             model.Field[bool]    `json:"shared,omitzero"`
	Weight             model.Field[float64] `json:"weight,omitzero"`
	Volume             model.Field[float64] `json:"volume,omitzero"`
	VAT                model.Field[int]     `json:"vat,omitzero"`
	VATEnabled         model.Field[bool]    `json:"vatEnabled,omitzero"`
	UseParentVAT       model.Field[bool]    `json:"useParentVat,omitzero"`
	DiscountProhibited model.Field[bool]    `json:"discountProhibited,omitzero"`

	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	UOM           model.Field[model.Reference] `json:"uom,omitzero"`
	Supplier      model.Field[model.Reference] `json:"supplier,omitzero"`
	Country       model.Field[model.Reference] `json:"country,omitzero"`

	BuyPrice   model.Field[model.Price]       `json:"buyPrice,omitzero"`
	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	// Barcodes: без поля или с пустым списком сервер создаёт один случайный EAN13; список из пустых
	// значений — ни одного штрихкода.
	Barcodes model.Field[[]model.Barcode] `json:"barcodes,omitzero"`
	Packs    model.Field[[]PackRequest]   `json:"packs,omitzero"`

	// SyncID — UUID, сгенерированный клиентом: повторный Create с тем же SyncID вернёт ранее
	// созданный товар. Это не общий Idempotency-Key, SDK сам Create не повторяет.
	// В UpdateRequest его нет: API помечает поле «только для чтения» и заполняет при создании.
	SyncID model.Field[string] `json:"syncId,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-izmenit-tovar
type UpdateRequest struct {
	Name               model.Field[string]  `json:"name,omitzero"`
	Code               model.Field[string]  `json:"code,omitzero"`
	ExternalCode       model.Field[string]  `json:"externalCode,omitzero"`
	Description        model.Field[string]  `json:"description,omitzero"`
	Article            model.Field[string]  `json:"article,omitzero"`
	Archived           model.Field[bool]    `json:"archived,omitzero"`
	Shared             model.Field[bool]    `json:"shared,omitzero"`
	Weight             model.Field[float64] `json:"weight,omitzero"`
	Volume             model.Field[float64] `json:"volume,omitzero"`
	VAT                model.Field[int]     `json:"vat,omitzero"`
	VATEnabled         model.Field[bool]    `json:"vatEnabled,omitzero"`
	UseParentVAT       model.Field[bool]    `json:"useParentVat,omitzero"`
	DiscountProhibited model.Field[bool]    `json:"discountProhibited,omitzero"`

	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	UOM           model.Field[model.Reference] `json:"uom,omitzero"`
	Supplier      model.Field[model.Reference] `json:"supplier,omitzero"`
	Country       model.Field[model.Reference] `json:"country,omitzero"`

	BuyPrice model.Field[model.Price] `json:"buyPrice,omitzero"`
	MinPrice model.Field[model.Price] `json:"minPrice,omitzero"`
	// Сервер сливает по priceType / meta; пустые массивы не очищают коллекции.
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	// Barcodes заменяет список целиком: отсутствующие штрихкоды удаляются, пустой список — все.
	Barcodes model.Field[[]model.Barcode] `json:"barcodes,omitzero"`
	// Packs заменяет коллекцию упаковок целиком; id упаковки только читается, поэтому прежние
	// упаковки сервер пересоздаёт.
	Packs model.Field[[]PackRequest] `json:"packs,omitzero"`
}

// PackRequest — упаковка в запросе: uom и quantity необходимы, штрихкодов не больше одного;
// переданный список штрихкодов заменяет прежний.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-upakovki-tovara
type PackRequest struct {
	UOM      model.Reference              `json:"uom"`
	Quantity float64                      `json:"quantity"`
	Barcodes model.Field[[]model.Barcode] `json:"barcodes,omitzero"`
}
