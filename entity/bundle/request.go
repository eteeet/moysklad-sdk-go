package bundle

import "github.com/eteeet/moysklad-sdk-go/model"

// ComponentRequest — компонент при записи комплекта. ID указывает существующий компонент для
// Update, без ID элемент — новый. Таблица источника помечает quantity «только для чтения», но
// примеры Create и добавления компонента его передают: без количества компонент не имеет
// смысла, поэтому SDK его пишет (исключение утверждено, см. docs/catalog.md).
type ComponentRequest struct {
	ID         string          `json:"id,omitempty"`
	Assortment model.Reference `json:"assortment,omitzero"`
	Quantity   float64         `json:"quantity"`
}

// CreateRequest: необходимы непустое name и непустой components; лимит 50 компонентов
// и обязательность assortment проверяет сервер. Штрихкоды: если barcodes
// не передан или пуст, сервер создаёт один случайный EAN13; список из пустых значений — ни одного.
// Images (загрузка Base64) и files в SDK v1 не пишутся.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-sozdat-komplekt
type CreateRequest struct {
	Name               string               `json:"name"`
	Components         []ComponentRequest   `json:"components"`
	Code               model.Field[string]  `json:"code,omitzero"`
	ExternalCode       model.Field[string]  `json:"externalCode,omitzero"`
	Description        model.Field[string]  `json:"description,omitzero"`
	Article            model.Field[string]  `json:"article,omitzero"`
	Archived           model.Field[bool]    `json:"archived,omitzero"`
	Shared             model.Field[bool]    `json:"shared,omitzero"`
	VAT                model.Field[int]     `json:"vat,omitzero"`
	VATEnabled         model.Field[bool]    `json:"vatEnabled,omitzero"`
	UseParentVAT       model.Field[bool]    `json:"useParentVat,omitzero"`
	DiscountProhibited model.Field[bool]    `json:"discountProhibited,omitzero"`
	PaymentItemType    model.Field[string]  `json:"paymentItemType,omitzero"`
	TaxSystem          model.Field[string]  `json:"taxSystem,omitzero"`
	TrackingType       model.Field[string]  `json:"trackingType,omitzero"`
	PartialDisposal    model.Field[bool]    `json:"partialDisposal,omitzero"`
	TNVED              model.Field[string]  `json:"tnved,omitzero"`
	Weight             model.Field[float64] `json:"weight,omitzero"`
	Volume             model.Field[float64] `json:"volume,omitzero"`

	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	UOM           model.Field[model.Reference] `json:"uom,omitzero"`
	Country       model.Field[model.Reference] `json:"country,omitzero"`

	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	Overhead   model.Field[model.Price]       `json:"overhead,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Barcodes   model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	TasnifUZ   model.Field[model.TasnifUZ]    `json:"mod__tasnif__uz,omitzero"`
	MarkingUZ  model.Field[model.MarkingUZ]   `json:"mod__marking__uz,omitzero"`

	// SyncID — UUID клиента; сервер записывает его один раз при создании, в UpdateRequest его нет.
	// SDK сам Create не повторяет.
	SyncID model.Field[string] `json:"syncId,omitzero"`
}

// UpdateRequest. Доп. поля, цены продажи и компоненты сервер сливает с текущими: переданный
// элемент записывается или перезаписывается (компонент — по ID), непереданный не меняется.
// Семантику null источник не описывает: Null в Field отправляется как есть, ответ — за сервером. Barcodes, напротив,
// заменяются целиком: пустой список удаляет все штрихкоды.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-izmenit-komplekt
type UpdateRequest struct {
	Name               model.Field[string]             `json:"name,omitzero"`
	Components         model.Field[[]ComponentRequest] `json:"components,omitzero"`
	Code               model.Field[string]             `json:"code,omitzero"`
	ExternalCode       model.Field[string]             `json:"externalCode,omitzero"`
	Description        model.Field[string]             `json:"description,omitzero"`
	Article            model.Field[string]             `json:"article,omitzero"`
	Archived           model.Field[bool]               `json:"archived,omitzero"`
	Shared             model.Field[bool]               `json:"shared,omitzero"`
	VAT                model.Field[int]                `json:"vat,omitzero"`
	VATEnabled         model.Field[bool]               `json:"vatEnabled,omitzero"`
	UseParentVAT       model.Field[bool]               `json:"useParentVat,omitzero"`
	DiscountProhibited model.Field[bool]               `json:"discountProhibited,omitzero"`
	PaymentItemType    model.Field[string]             `json:"paymentItemType,omitzero"`
	TaxSystem          model.Field[string]             `json:"taxSystem,omitzero"`
	TrackingType       model.Field[string]             `json:"trackingType,omitzero"`
	PartialDisposal    model.Field[bool]               `json:"partialDisposal,omitzero"`
	TNVED              model.Field[string]             `json:"tnved,omitzero"`
	Weight             model.Field[float64]            `json:"weight,omitzero"`
	Volume             model.Field[float64]            `json:"volume,omitzero"`

	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	UOM           model.Field[model.Reference] `json:"uom,omitzero"`
	Country       model.Field[model.Reference] `json:"country,omitzero"`

	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	Overhead   model.Field[model.Price]       `json:"overhead,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Barcodes   model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	TasnifUZ   model.Field[model.TasnifUZ]    `json:"mod__tasnif__uz,omitzero"`
	MarkingUZ  model.Field[model.MarkingUZ]   `json:"mod__marking__uz,omitzero"`
}
