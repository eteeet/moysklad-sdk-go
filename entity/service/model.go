package service

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Service — услуга. Модель чтения без expand. Деньги в ценах — копейки; VAT — процент.
// PaymentItemType и TaxSystem открыты (SERVICE, WORK …; GENERAL_TAX_SYSTEM …).
// Противоречия источника: таблица называет поле minPrice, примеры — minprice (SDK следует
// таблице и upstream-фикстуре); вступление говорит «услуга без закупочной цены», но buyPrice
// в таблице есть — поле читается и пишется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service
type Service struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	AccountID    string         `json:"accountId"`
	Updated      model.DateTime `json:"updated,omitzero"`
	Name         string         `json:"name"`
	Code         string         `json:"code"`
	ExternalCode string         `json:"externalCode"`
	Description  string         `json:"description"`
	PathName     string         `json:"pathName"`
	Archived     bool           `json:"archived"`
	Shared       bool           `json:"shared"`
	// SyncID только для чтения; задаётся при создании.
	SyncID string `json:"syncId"`

	VAT                 int    `json:"vat"`
	VATEnabled          bool   `json:"vatEnabled"`
	UseParentVAT        bool   `json:"useParentVat"`
	EffectiveVAT        int    `json:"effectiveVat"`
	EffectiveVATEnabled bool   `json:"effectiveVatEnabled"`
	DiscountProhibited  bool   `json:"discountProhibited"`
	PaymentItemType     string `json:"paymentItemType"`
	TaxSystem           string `json:"taxSystem"`

	BuyPrice   *model.Price      `json:"buyPrice,omitzero"`
	MinPrice   *model.Price      `json:"minPrice,omitzero"`
	SalePrices []model.SalePrice `json:"salePrices"`
	Barcodes   []model.Barcode   `json:"barcodes"`
	Attributes []model.Attribute `json:"attributes"`
	TasnifUZ   *model.TasnifUZ   `json:"mod__tasnif__uz,omitzero"`

	Owner         *model.Reference `json:"owner,omitzero"`
	Group         *model.Reference `json:"group,omitzero"`
	ProductFolder *model.Reference `json:"productFolder,omitzero"`
	UOM           *model.Reference `json:"uom,omitzero"`
	// Files — только метаданные коллекции файлов; сами файлы читаются отдельным ресурсом.
	Files *model.Collection[jsontext.Value] `json:"files,omitzero"`
}
