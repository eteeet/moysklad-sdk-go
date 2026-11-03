package productfolder

import "github.com/eteeet/moysklad-sdk-go/model"

// ProductFolder — группа товаров. Vat и EffectiveVat — проценты; TaxSystem — открытый enum
// (GENERAL_TAX_SYSTEM, SIMPLIFIED_TAX_SYSTEM_INCOME, TAX_SYSTEM_SAME_AS_GROUP, …).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/productfolder
type ProductFolder struct {
	Meta                model.Meta       `json:"meta"`
	ID                  string           `json:"id"`
	AccountID           string           `json:"accountId"`
	Updated             model.DateTime   `json:"updated,omitzero"`
	Name                string           `json:"name"`
	Description         string           `json:"description"`
	Code                string           `json:"code"`
	ExternalCode        string           `json:"externalCode"`
	Archived            bool             `json:"archived"`
	PathName            string           `json:"pathName"`
	ProductFolder       *model.Reference `json:"productFolder,omitzero"`
	Vat                 int              `json:"vat"`
	VatEnabled          bool             `json:"vatEnabled"`
	UseParentVat        bool             `json:"useParentVat"`
	EffectiveVat        int              `json:"effectiveVat"`
	EffectiveVatEnabled bool             `json:"effectiveVatEnabled"`
	TaxSystem           string           `json:"taxSystem"`
	Shared              bool             `json:"shared"`
	Owner               *model.Reference `json:"owner,omitzero"`
	Group               *model.Reference `json:"group,omitzero"`
}
