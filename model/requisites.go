package model

// RequisitesUZ — реквизиты Узбекистана (mod__requisites__uz, +UZ): валидны только у сущностей
// узбекского типа (legalUZ, entrepreneurUZ, individualUZ), это проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#5-rekvizity-uzbekistana
type RequisitesUZ struct {
	INN               string   `json:"inn"`
	PINFL             string   `json:"pinfl"`
	OKED              string   `json:"oked"`
	VATPayerRegCode   string   `json:"vatPayerRegCode"`
	CertificateNumber string   `json:"certificateNumber"`
	CertificateDate   DateTime `json:"certificateDate,omitzero"`
}

type RequisitesUZRequest struct {
	INN               Field[string]   `json:"inn,omitzero"`
	PINFL             Field[string]   `json:"pinfl,omitzero"`
	OKED              Field[string]   `json:"oked,omitzero"`
	VATPayerRegCode   Field[string]   `json:"vatPayerRegCode,omitzero"`
	CertificateNumber Field[string]   `json:"certificateNumber,omitzero"`
	CertificateDate   Field[DateTime] `json:"certificateDate,omitzero"`
}

// RequisitesKZ — реквизиты Казахстана (mod__requisites__kz, +KZ): валидны только у сущностей
// казахского типа, это проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#5-rekvizity-kazahstana
type RequisitesKZ struct {
	BIN                  string   `json:"bin"`
	IIN                  string   `json:"iin"`
	OKED                 string   `json:"oked"`
	KBE                  string   `json:"kbe"`
	VATCertificateNumber string   `json:"vatCertificateNumber"`
	VATCertificateDate   DateTime `json:"vatCertificateDate,omitzero"`
}

type RequisitesKZRequest struct {
	BIN                  Field[string]   `json:"bin,omitzero"`
	IIN                  Field[string]   `json:"iin,omitzero"`
	OKED                 Field[string]   `json:"oked,omitzero"`
	KBE                  Field[string]   `json:"kbe,omitzero"`
	VATCertificateNumber Field[string]   `json:"vatCertificateNumber,omitzero"`
	VATCertificateDate   Field[DateTime] `json:"vatCertificateDate,omitzero"`
}
