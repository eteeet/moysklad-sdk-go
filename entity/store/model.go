package store

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Store — склад. Zones и Slots — только метаданные коллекций (rows приходят с expand, который
// SDK не запрашивает); содержимое читают ListZones и ListSlots.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#2-sklad
type Store struct {
	Meta         model.Meta        `json:"meta"`
	ID           string            `json:"id"`
	AccountID    string            `json:"accountId"`
	Updated      model.DateTime    `json:"updated,omitzero"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Code         string            `json:"code"`
	ExternalCode string            `json:"externalCode"`
	Archived     bool              `json:"archived"`
	Shared       bool              `json:"shared"`
	PathName     string            `json:"pathName"`
	Address      string            `json:"address"`
	AddressFull  *model.Address    `json:"addressFull,omitzero"`
	Attributes   []model.Attribute `json:"attributes"`
	Parent       *model.Reference  `json:"parent,omitzero"`
	Owner        *model.Reference  `json:"owner,omitzero"`
	Group        *model.Reference  `json:"group,omitzero"`

	Zones *model.Collection[jsontext.Value] `json:"zones,omitzero"`
	Slots *model.Collection[jsontext.Value] `json:"slots,omitzero"`
}

// Zone — зона склада (meta.type storezone).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-zony-sklada
type Zone struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	AccountID    string         `json:"accountId"`
	Updated      model.DateTime `json:"updated,omitzero"`
	Name         string         `json:"name"`
	ExternalCode string         `json:"externalCode"`
}

// Slot — ячейка склада (meta.type slot).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-yachejki-sklada
type Slot struct {
	Meta         model.Meta       `json:"meta"`
	ID           string           `json:"id"`
	AccountID    string           `json:"accountId"`
	Updated      model.DateTime   `json:"updated,omitzero"`
	Name         string           `json:"name"`
	ExternalCode string           `json:"externalCode"`
	Barcode      string           `json:"barcode"`
	Zone         *model.Reference `json:"zone,omitzero"`
}
