package bundle

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Bundle — комплект. Модель чтения без expand. Деньги в ценах и overhead — копейки; VAT — процент.
// PaymentItemType, TaxSystem и TrackingType открыты (GOOD …; GENERAL_TAX_SYSTEM …; NOT_TRACKED …).
// Противоречие источника: таблица называет поле minPrice, пример ответа Update — minprice;
// SDK следует таблице и upstream-фикстуре.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle
type Bundle struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	AccountID    string         `json:"accountId"`
	Updated      model.DateTime `json:"updated,omitzero"`
	Name         string         `json:"name"`
	Code         string         `json:"code"`
	ExternalCode string         `json:"externalCode"`
	Description  string         `json:"description"`
	Article      string         `json:"article"`
	PathName     string         `json:"pathName"`
	Archived     bool           `json:"archived"`
	Shared       bool           `json:"shared"`
	// SyncID только для чтения; задаётся при создании.
	SyncID string `json:"syncId"`

	VAT                 int     `json:"vat"`
	VATEnabled          bool    `json:"vatEnabled"`
	UseParentVAT        bool    `json:"useParentVat"`
	EffectiveVAT        int     `json:"effectiveVat"`
	EffectiveVATEnabled bool    `json:"effectiveVatEnabled"`
	DiscountProhibited  bool    `json:"discountProhibited"`
	PaymentItemType     string  `json:"paymentItemType"`
	TaxSystem           string  `json:"taxSystem"`
	TrackingType        string  `json:"trackingType"`
	PartialDisposal     bool    `json:"partialDisposal"`
	TNVED               string  `json:"tnved"`
	Weight              float64 `json:"weight"`
	Volume              float64 `json:"volume"`

	MinPrice   *model.Price      `json:"minPrice,omitzero"`
	Overhead   *model.Price      `json:"overhead,omitzero"`
	SalePrices []model.SalePrice `json:"salePrices"`
	Barcodes   []model.Barcode   `json:"barcodes"`
	Attributes []model.Attribute `json:"attributes"`
	TasnifUZ   *model.TasnifUZ   `json:"mod__tasnif__uz,omitzero"`
	MarkingUZ  *model.MarkingUZ  `json:"mod__marking__uz,omitzero"`

	Owner         *model.Reference `json:"owner,omitzero"`
	Group         *model.Reference `json:"group,omitzero"`
	ProductFolder *model.Reference `json:"productFolder,omitzero"`
	UOM           *model.Reference `json:"uom,omitzero"`
	Country       *model.Reference `json:"country,omitzero"`

	// Components — метаданные коллекции компонентов; Rows заполнены, только если их вернул сервер.
	// Без скрытых запросов: строки читает GetComponents.
	Components *model.Collection[Component] `json:"components,omitzero"`
	// Images и Files — только метаданные коллекций; содержимое читается отдельными ресурсами.
	Images *model.Collection[jsontext.Value] `json:"images,omitzero"`
	Files  *model.Collection[jsontext.Value] `json:"files,omitzero"`
}

// Component — компонент комплекта: ссылка на товар, услугу или модификацию и количество.
// Не позиция документа: цены, скидки и НДС у компонента нет. Тип Assortment открыт
// (таблица называет товар/услугу/партию, вступление — товар/услугу/модификацию).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#5-komponenty-komplekta
type Component struct {
	Meta       model.Meta      `json:"meta"`
	ID         string          `json:"id"`
	AccountID  string          `json:"accountId"`
	Assortment model.Reference `json:"assortment"`
	Quantity   float64         `json:"quantity"`
}
