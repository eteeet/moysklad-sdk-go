package model

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-ceny-prodazhi
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-zakupochnaya-cena
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-minimalnaya-cena

type Price struct {
	Value    float64   `json:"value"` // В копейках (сотых долях валюты).
	Currency Reference `json:"currency,omitzero"`
}

type SalePrice struct {
	Value     float64   `json:"value"` // В копейках (сотых долях валюты).
	Currency  Reference `json:"currency,omitzero"`
	PriceType PriceType `json:"priceType,omitzero"`
}

type PriceType struct {
	Meta         Meta   `json:"meta"`
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	ExternalCode string `json:"externalCode,omitempty"`
}
