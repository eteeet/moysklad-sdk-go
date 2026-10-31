package model

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#5-adres
// Незаданная часть адреса не отправляется; Null кодируется, но не обещает очистку на сервере.
type AddressRequest struct {
	AddInfo    Field[string]    `json:"addInfo,omitzero"`
	Apartment  Field[string]    `json:"apartment,omitzero"`
	City       Field[string]    `json:"city,omitzero"`
	Comment    Field[string]    `json:"comment,omitzero"`
	FIASCodeRU Field[string]    `json:"fiasCode__ru,omitzero"`
	Country    Field[Reference] `json:"country,omitzero"`
	House      Field[string]    `json:"house,omitzero"`
	PostalCode Field[string]    `json:"postalCode,omitzero"`
	Region     Field[Reference] `json:"region,omitzero"`
	Street     Field[string]    `json:"street,omitzero"`
}
