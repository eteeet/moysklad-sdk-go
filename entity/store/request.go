package store

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен только name. Адрес — строкой или структурой; при обеих строка
// игнорируется сервером. null адреса источник не поддерживает: удаляет адрес Address = "".
// parent и pathName маркера «только для чтения» не имеют и передаются примером PUT.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-sklad
type CreateRequest struct {
	Name         string                            `json:"name"`
	Description  model.Field[string]               `json:"description,omitzero"`
	Code         model.Field[string]               `json:"code,omitzero"`
	ExternalCode model.Field[string]               `json:"externalCode,omitzero"`
	Archived     model.Field[bool]                 `json:"archived,omitzero"`
	Shared       model.Field[bool]                 `json:"shared,omitzero"`
	PathName     model.Field[string]               `json:"pathName,omitzero"`
	Address      model.Field[string]               `json:"address,omitzero"`
	AddressFull  model.Field[model.AddressRequest] `json:"addressFull,omitzero"`
	Attributes   model.Field[[]model.Attribute]    `json:"attributes,omitzero"`
	Parent       model.Field[model.Reference]      `json:"parent,omitzero"`
	Owner        model.Field[model.Reference]      `json:"owner,omitzero"`
	Group        model.Field[model.Reference]      `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-sklad
type UpdateRequest struct {
	Name         model.Field[string]               `json:"name,omitzero"`
	Description  model.Field[string]               `json:"description,omitzero"`
	Code         model.Field[string]               `json:"code,omitzero"`
	ExternalCode model.Field[string]               `json:"externalCode,omitzero"`
	Archived     model.Field[bool]                 `json:"archived,omitzero"`
	Shared       model.Field[bool]                 `json:"shared,omitzero"`
	PathName     model.Field[string]               `json:"pathName,omitzero"`
	Address      model.Field[string]               `json:"address,omitzero"`
	AddressFull  model.Field[model.AddressRequest] `json:"addressFull,omitzero"`
	Attributes   model.Field[[]model.Attribute]    `json:"attributes,omitzero"`
	Parent       model.Field[model.Reference]      `json:"parent,omitzero"`
	Owner        model.Field[model.Reference]      `json:"owner,omitzero"`
	Group        model.Field[model.Reference]      `json:"group,omitzero"`
}

// ZoneCreateRequest: обязателен name; имена зон уникальны в пределах склада (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-zonu-sklada
type ZoneCreateRequest struct {
	Name         string              `json:"name"`
	ExternalCode model.Field[string] `json:"externalCode,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-zonu-sklada
type ZoneUpdateRequest struct {
	Name         model.Field[string] `json:"name,omitzero"`
	ExternalCode model.Field[string] `json:"externalCode,omitzero"`
}

// SlotCreateRequest: обязателен name; имена и штрихкоды ячеек уникальны в пределах склада.
// Zone таблица помечает «только для чтения», но источник даёт отдельные примеры создания и
// изменения ячейки с зоной — SDK поле передаёт: без него ячейку нельзя привязать к зоне.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-yachejku-sklada
type SlotCreateRequest struct {
	Name         string                       `json:"name"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Barcode      model.Field[string]          `json:"barcode,omitzero"`
	Zone         model.Field[model.Reference] `json:"zone,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-yachejku-sklada
type SlotUpdateRequest struct {
	Name         model.Field[string]          `json:"name,omitzero"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Barcode      model.Field[string]          `json:"barcode,omitzero"`
	Zone         model.Field[model.Reference] `json:"zone,omitzero"`
}
