package saleplatform

import "github.com/eteeet/moysklad-sdk-go/model"

// SalePlatform — площадка для продаж; только чтение. SalePlatformGroup — открытый enum
// (MARKETPLACES, MESSENGERS, OFFLINE_STORES, ONLINE_STORES, SOCIAL_NETWORKS).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleplatform
type SalePlatform struct {
	Meta              model.Meta `json:"meta"`
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	SalePlatformGroup string     `json:"salePlatformGroup"`
}
