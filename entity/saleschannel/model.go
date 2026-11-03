package saleschannel

import "github.com/eteeet/moysklad-sdk-go/model"

// SalesChannel — канал продаж. Type — открытый enum (MESSENGER, SOCIAL_NETWORK, MARKETPLACE, ECOMMERCE,
// CLASSIFIED_ADS, DIRECT_SALES, RETAIL_SALES, OTHER).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel
type SalesChannel struct {
	Meta         model.Meta       `json:"meta"`
	ID           string           `json:"id"`
	AccountID    string           `json:"accountId"`
	Updated      model.DateTime   `json:"updated,omitzero"`
	Name         string           `json:"name"`
	Type         string           `json:"type"`
	Description  string           `json:"description"`
	Code         string           `json:"code"`
	ExternalCode string           `json:"externalCode"`
	Archived     bool             `json:"archived"`
	Shared       bool             `json:"shared"`
	Owner        *model.Reference `json:"owner,omitzero"`
	Group        *model.Reference `json:"group,omitzero"`
}
