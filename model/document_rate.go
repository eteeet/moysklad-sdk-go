package model

// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L689-L699
type DocumentRate struct {
	Currency *Reference `json:"currency,omitzero"`
	// nil — курс 1: сервер не передаёт value, равный 1.
	Value *float64 `json:"value,omitzero"`
}

// Незаданные поля не отправляются; внешний rate:null возвращает валюту по умолчанию
// с серверным пересчётом цен и сумм. SDK передаёт JSON без расчётов.
type DocumentRateRequest struct {
	Currency Field[Reference] `json:"currency,omitzero"`
	Value    Field[float64]   `json:"value,omitzero"`
}
