package product

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-sozdat-tovar
// Изображения (images: filename и content) в теле товара SDK не пишет — их загружает пакет images
// (Add, Replace).

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
	Barcodes          model.Field[[]model.Barcode]           `json:"barcodes,omitzero"`
	Packs             model.Field[[]PackRequest]             `json:"packs,omitzero"`
	IsSerialTrackable model.Field[bool]                      `json:"isSerialTrackable,omitzero"`
	Alcoholic         model.Field[AlcoholicRequest]          `json:"alcoholic,omitzero"`
	MinimumBalance    model.Field[float64]                   `json:"minimumBalance,omitzero"`
	PartialDisposal   model.Field[bool]                      `json:"partialDisposal,omitzero"`
	PaymentItemType   model.Field[string]                    `json:"paymentItemType,omitzero"`
	PPEType           model.Field[string]                    `json:"ppeType,omitzero"`
	TaxSystem         model.Field[string]                    `json:"taxSystem,omitzero"`
	TrackingType      model.Field[string]                    `json:"trackingType,omitzero"`
	TNVED             model.Field[string]                    `json:"tnved,omitzero"`
	Things            model.Field[[]string]                  `json:"things,omitzero"`
	Weighed           model.Field[bool]                      `json:"weighed,omitzero"`
	OnTap             model.Field[bool]                      `json:"onTap,omitzero"`
	TasnifUZ          model.Field[model.TasnifUZ]            `json:"mod__tasnif__uz,omitzero"`
	MarkingUZ         model.Field[model.MarkingUZ]           `json:"mod__marking__uz,omitzero"`
	NctKZ             model.Field[model.NctKZ]               `json:"mod__nct__kz,omitzero"`
	MinimumStock      model.Field[model.MinimumStockRequest] `json:"minimumStock,omitzero"`

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
	// Packs заменяет коллекцию упаковок целиком: упаковка с ID — существующая (пример изменения
	// товара сохраняет её id), без ID — новая; не переданные удаляются.
	Packs             model.Field[[]PackRequest]             `json:"packs,omitzero"`
	IsSerialTrackable model.Field[bool]                      `json:"isSerialTrackable,omitzero"`
	Alcoholic         model.Field[AlcoholicRequest]          `json:"alcoholic,omitzero"`
	MinimumBalance    model.Field[float64]                   `json:"minimumBalance,omitzero"`
	PartialDisposal   model.Field[bool]                      `json:"partialDisposal,omitzero"`
	PaymentItemType   model.Field[string]                    `json:"paymentItemType,omitzero"`
	PPEType           model.Field[string]                    `json:"ppeType,omitzero"`
	TaxSystem         model.Field[string]                    `json:"taxSystem,omitzero"`
	TrackingType      model.Field[string]                    `json:"trackingType,omitzero"`
	TNVED             model.Field[string]                    `json:"tnved,omitzero"`
	Things            model.Field[[]string]                  `json:"things,omitzero"`
	Weighed           model.Field[bool]                      `json:"weighed,omitzero"`
	OnTap             model.Field[bool]                      `json:"onTap,omitzero"`
	TasnifUZ          model.Field[model.TasnifUZ]            `json:"mod__tasnif__uz,omitzero"`
	MarkingUZ         model.Field[model.MarkingUZ]           `json:"mod__marking__uz,omitzero"`
	NctKZ             model.Field[model.NctKZ]               `json:"mod__nct__kz,omitzero"`
	MinimumStock      model.Field[model.MinimumStockRequest] `json:"minimumStock,omitzero"`
}

// PackRequest — упаковка в запросе: uom и quantity необходимы, штрихкодов не больше одного;
// переданный список штрихкодов заменяет прежний. ID адресует существующую упаковку (id таблица
// помечает «только для чтения»: значение не меняется, им выбирается упаковка, как в примере
// изменения товара); пустой ID — новая упаковка.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-upakovki-tovara
type PackRequest struct {
	ID       string                       `json:"id,omitempty"`
	UOM      model.Reference              `json:"uom"`
	Quantity float64                      `json:"quantity"`
	Barcodes model.Field[[]model.Barcode] `json:"barcodes,omitzero"`
}

// AlcoholicRequest — поля алкогольной продукции; объект целиком или частично — как в таблице.
type AlcoholicRequest struct {
	Excise   model.Field[bool]    `json:"excise,omitzero"`
	Type     model.Field[int]     `json:"type,omitzero"`
	Strength model.Field[float64] `json:"strength,omitzero"`
	Volume   model.Field[float64] `json:"volume,omitzero"`
}
