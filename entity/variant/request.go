package variant

import "github.com/eteeet/moysklad-sdk-go/model"

// CharacteristicValue — значение характеристики при записи; поиск описания по ID или Name —
// см. Characteristic. Новые описания не создаются: см. CreateCharacteristic.
type CharacteristicValue struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

// PackRequest — упаковка при записи; массив barcodes может содержать только один штрихкод.
type PackRequest struct {
	Barcodes   []model.Barcode `json:"barcodes"`
	ParentPack model.Reference `json:"parentpack"`
}

// CreateRequest: необходимы product и непустой characteristics. Правила создания штрихкодов
// раздел variant копирует из bundle под заголовком «при создании комплекта» (нет barcodes или
// пустой список — один случайный EAN13; список пустых значений — ни одного); для модификации
// это не подтверждено. НДС и налоговых полей у модификации в источнике нет. Images и
// minimumStock в SDK v1 не пишутся.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-sozdat-modifikaciyu
type CreateRequest struct {
	Product         model.Reference       `json:"product"`
	Characteristics []CharacteristicValue `json:"characteristics"`

	Name               model.Field[string] `json:"name,omitzero"`
	Code               model.Field[string] `json:"code,omitzero"`
	ExternalCode       model.Field[string] `json:"externalCode,omitzero"`
	Description        model.Field[string] `json:"description,omitzero"`
	Article            model.Field[string] `json:"article,omitzero"`
	Archived           model.Field[bool]   `json:"archived,omitzero"`
	DiscountProhibited model.Field[bool]   `json:"discountProhibited,omitzero"`

	BuyPrice   model.Field[model.Price]       `json:"buyPrice,omitzero"`
	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Barcodes   model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Packs      model.Field[[]PackRequest]     `json:"packs,omitzero"`
}

// UpdateRequest. Characteristics заменяются целиком: непереданная характеристика теряет значение.
// Barcodes тоже заменяются целиком, пустой список удаляет все. Цены продажи сервер сливает по
// типу цены. Смену Product источник не запрещает и не описывает — ограничения проверяет сервер.
// Семантику null источник не описывает: Null в Field отправляется как есть.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/variant#3-izmenit-modifikaciyu
type UpdateRequest struct {
	Product         model.Field[model.Reference]       `json:"product,omitzero"`
	Characteristics model.Field[[]CharacteristicValue] `json:"characteristics,omitzero"`

	Name               model.Field[string] `json:"name,omitzero"`
	Code               model.Field[string] `json:"code,omitzero"`
	ExternalCode       model.Field[string] `json:"externalCode,omitzero"`
	Description        model.Field[string] `json:"description,omitzero"`
	Article            model.Field[string] `json:"article,omitzero"`
	Archived           model.Field[bool]   `json:"archived,omitzero"`
	DiscountProhibited model.Field[bool]   `json:"discountProhibited,omitzero"`

	BuyPrice   model.Field[model.Price]       `json:"buyPrice,omitzero"`
	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Barcodes   model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Packs      model.Field[[]PackRequest]     `json:"packs,omitzero"`
}
