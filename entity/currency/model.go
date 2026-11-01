package currency

import "github.com/eteeet/moysklad-sdk-go/model"

// Currency — валюта. Rate — курс, а не денежная сумма: передаётся как есть, без копеек.
// System — валюта из системного справочника: часть её полей сервер менять не даёт.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency
type Currency struct {
	Meta           model.Meta `json:"meta"`
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	FullName       string     `json:"fullName"`
	Code           string     `json:"code"`
	ISOCode        string     `json:"isoCode"`
	Archived       bool       `json:"archived"`
	Default        bool       `json:"default"`
	System         bool       `json:"system"`
	Rate           float64    `json:"rate"`
	Multiplicity   int        `json:"multiplicity"`
	Indirect       bool       `json:"indirect"`
	Margin         float64    `json:"margin"`
	RateUpdateType string     `json:"rateUpdateType"` // auto или manual; набор открыт
	MajorUnit      *Unit      `json:"majorUnit,omitzero"`
	MinorUnit      *Unit      `json:"minorUnit,omitzero"`
}

// Unit — формы названия единицы валюты для числительных 1, 2 и 5; Gender — masculine или feminine.
type Unit struct {
	Gender string `json:"gender"`
	S1     string `json:"s1"`
	S2     string `json:"s2"`
	S5     string `json:"s5"`
}
