package product

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Модель чтения без expand; для записи используются CreateRequest и UpdateRequest.
// PaymentItemType, TaxSystem, TrackingType и PPEType — открытые enum. Weighed и OnTap сервер
// выводит только со значением true. Несочетаемые признаки (isSerialTrackable с weighed, alcoholic,
// ppeType, trackingType, onTap и т.д.) проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product
type Product struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	AccountID    string         `json:"accountId"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Code         string         `json:"code"`
	ExternalCode string         `json:"externalCode"`
	Article      string         `json:"article"`
	Archived     bool           `json:"archived"`
	Shared       bool           `json:"shared"`
	PathName     string         `json:"pathName"`
	Updated      model.DateTime `json:"updated,omitzero"`
	SyncID       string         `json:"syncId"`

	BuyPrice   *model.Price      `json:"buyPrice,omitzero"`
	MinPrice   *model.Price      `json:"minPrice,omitzero"`
	SalePrices []model.SalePrice `json:"salePrices"`
	Attributes []model.Attribute `json:"attributes"`
	Barcodes   []model.Barcode   `json:"barcodes"`
	Packs      []Pack            `json:"packs"`

	Weight float64 `json:"weight"`
	Volume float64 `json:"volume"`

	VAT                 int  `json:"vat"`
	VATEnabled          bool `json:"vatEnabled"`
	UseParentVAT        bool `json:"useParentVat"`
	EffectiveVAT        int  `json:"effectiveVat"`
	EffectiveVATEnabled bool `json:"effectiveVatEnabled"`
	DiscountProhibited  bool `json:"discountProhibited"`
	IsSerialTrackable   bool `json:"isSerialTrackable"`

	Alcoholic       *Alcoholic `json:"alcoholic,omitzero"`
	MinimumBalance  float64    `json:"minimumBalance"`
	PartialDisposal bool       `json:"partialDisposal"`
	PaymentItemType string     `json:"paymentItemType"`
	PPEType         string     `json:"ppeType"`
	TaxSystem       string     `json:"taxSystem"`
	TrackingType    string     `json:"trackingType"`
	TNVED           string     `json:"tnved"`
	Things          []string   `json:"things"`
	Weighed         bool       `json:"weighed"`
	OnTap           bool       `json:"onTap"`
	VariantsCount   int        `json:"variantsCount"`
	// MinimumStock приходит только при List с fields=minimumStock.
	MinimumStock *model.MinimumStock `json:"minimumStock,omitzero"`

	// Регионоспецифичные модули: TasnifUZ и MarkingUZ — с transport.WithAcceptRegion(RegionUZ).
	TasnifUZ  *model.TasnifUZ  `json:"mod__tasnif__uz,omitzero"`
	MarkingUZ *model.MarkingUZ `json:"mod__marking__uz,omitzero"`
	NctKZ     *model.NctKZ     `json:"mod__nct__kz,omitzero"`

	// Images и Files — только метаданные коллекций (изображения — отдельный ресурс, файлы — files).
	Images *model.Collection[jsontext.Value] `json:"images,omitzero"`
	Files  *model.Collection[jsontext.Value] `json:"files,omitzero"`

	Owner         *model.Reference `json:"owner,omitzero"`
	Group         *model.Reference `json:"group,omitzero"`
	ProductFolder *model.Reference `json:"productFolder,omitzero"`
	UOM           *model.Reference `json:"uom,omitzero"`
	Supplier      *model.Reference `json:"supplier,omitzero"`
	Country       *model.Reference `json:"country,omitzero"`
}

// Pack — упаковка товара: единица измерения (не совпадает с единицей товара — проверяет
// сервер), количество товара в упаковке и не более одного штрихкода (без штрихкода поля нет).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-upakovki-tovara
type Pack struct {
	ID       string           `json:"id"`
	UOM      *model.Reference `json:"uom,omitzero"`
	Quantity float64          `json:"quantity"`
	Barcodes []model.Barcode  `json:"barcodes"`
}

// Alcoholic — алкогольная продукция: Type — код вида продукции, Strength — крепость, Volume —
// объём тары.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-obekt-soderzhashij-polya-alkogolnoj-produkcii
type Alcoholic struct {
	Excise   bool    `json:"excise"`
	Type     int     `json:"type"`
	Strength float64 `json:"strength"`
	Volume   float64 `json:"volume"`
}
