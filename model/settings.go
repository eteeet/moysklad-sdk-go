package model

// UniqueCodeRules — настройки уникальности кода справочника: проверка уникальности и
// автоматическое заполнение кода при создании.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-nastrojki-spravochnika-kontragentov
type UniqueCodeRules struct {
	CheckUniqueCode bool `json:"checkUniqueCode"`
	FillUniqueCode  bool `json:"fillUniqueCode"`
}

type UniqueCodeRulesRequest struct {
	CheckUniqueCode Field[bool] `json:"checkUniqueCode,omitzero"`
	FillUniqueCode  Field[bool] `json:"fillUniqueCode,omitzero"`
}
