package product

import "github.com/eteeet/moysklad-sdk-go/model"

// Частичная модель чтения без expand; для записи будут отдельные типы.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product
// TODO: даты, цены, штрихкоды, изображения, файлы, упаковки, доп. поля и остальные атрибуты.
type Product struct {
	Meta         model.Meta `json:"meta"`
	ID           string     `json:"id"`
	AccountID    string     `json:"accountId"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Code         string     `json:"code"`
	ExternalCode string     `json:"externalCode"`
	Article      string     `json:"article"`
	Archived     bool       `json:"archived"`
	Shared       bool       `json:"shared"`
	PathName     string     `json:"pathName"`

	Weight float64 `json:"weight"`
	Volume float64 `json:"volume"`

	VAT                 int  `json:"vat"`
	VATEnabled          bool `json:"vatEnabled"`
	UseParentVAT        bool `json:"useParentVat"`
	EffectiveVAT        int  `json:"effectiveVat"`
	EffectiveVATEnabled bool `json:"effectiveVatEnabled"`
	DiscountProhibited  bool `json:"discountProhibited"`
	IsSerialTrackable   bool `json:"isSerialTrackable"`

	Owner         *model.Reference `json:"owner,omitzero"`
	Group         *model.Reference `json:"group,omitzero"`
	ProductFolder *model.Reference `json:"productFolder,omitzero"`
	UOM           *model.Reference `json:"uom,omitzero"`
	Supplier      *model.Reference `json:"supplier,omitzero"`
	Country       *model.Reference `json:"country,omitzero"`
}
