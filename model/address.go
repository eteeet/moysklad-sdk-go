package model

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#5-adres
type Address struct {
	AddInfo    string     `json:"addInfo"`
	Apartment  string     `json:"apartment"`
	City       string     `json:"city"`
	Comment    string     `json:"comment"`
	FIASCodeRU string     `json:"fiasCode__ru"`
	Country    *Reference `json:"country,omitzero"`
	House      string     `json:"house"`
	PostalCode string     `json:"postalCode"`
	Region     *Reference `json:"region,omitzero"`
	Street     string     `json:"street"`
}
