package organization

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-yurlica
type Organization struct {
	Meta         model.Meta        `json:"meta"`
	ID           string            `json:"id"`
	AccountID    string            `json:"accountId"`
	Created      model.DateTime    `json:"created,omitzero"`
	Updated      model.DateTime    `json:"updated,omitzero"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Code         string            `json:"code"`
	ExternalCode string            `json:"externalCode"`
	Archived     bool              `json:"archived"`
	Shared       bool              `json:"shared"`
	CompanyType  model.CompanyType `json:"companyType"`
	Email        string            `json:"email"`
	Phone        string            `json:"phone"`
	Fax          string            `json:"fax"`

	LegalTitle        string         `json:"legalTitle"`
	LegalFirstName    string         `json:"legalFirstName"`
	LegalLastName     string         `json:"legalLastName"`
	LegalMiddleName   string         `json:"legalMiddleName"`
	INN               string         `json:"inn"`
	KPP               string         `json:"kpp"`
	OGRN              string         `json:"ogrn"`
	OGRNIP            string         `json:"ogrnip"`
	OKPO              string         `json:"okpo"`
	CertificateNumber string         `json:"certificateNumber"`
	CertificateDate   model.DateTime `json:"certificateDate,omitzero"`

	ActualAddress     string            `json:"actualAddress"`
	LegalAddress      string            `json:"legalAddress"`
	ActualAddressFull *model.Address    `json:"actualAddressFull,omitzero"`
	LegalAddressFull  *model.Address    `json:"legalAddressFull,omitzero"`
	Attributes        []model.Attribute `json:"attributes"`
	Owner             *model.Reference  `json:"owner,omitzero"`
	Group             *model.Reference  `json:"group,omitzero"`
	BonusProgram      *model.Reference  `json:"bonusProgram,omitzero"`
	SyncID            string            `json:"syncId"`
	BonusPoints       int               `json:"bonusPoints"`

	PayerVAT               bool           `json:"payerVat"`
	Director               string         `json:"director"`
	DirectorPosition       string         `json:"directorPosition"`
	ChiefAccountant        string         `json:"chiefAccountant"`
	TrackingContractNumber string         `json:"trackingContractNumber"`
	TrackingContractDate   model.DateTime `json:"trackingContractDate,omitzero"`
	IsEGAISEnable          bool           `json:"isEgaisEnable"`

	// CompanyVat приходит только при payerVat = true. AdvancePaymentVat — процент ставки для
	// авансов (значение из существующих ставок НДС проверяет сервер).
	CompanyVat        *CompanyVat `json:"companyVat__ru,omitzero"`
	AdvancePaymentVat float64     `json:"advancePaymentVat"`
	DirectorSign      *Image      `json:"directorSign,omitzero"`
	ChiefAccountSign  *Image      `json:"chiefAccountSign,omitzero"`
	Stamp             *Image      `json:"stamp,omitzero"`

	RequisitesUZ *model.RequisitesUZ `json:"mod__requisites__uz,omitzero"`
	RequisitesKZ *model.RequisitesKZ `json:"mod__requisites__kz,omitzero"`
}

// CompanyVat — настройки НДС юрлица (companyVat__ru). DefaultCompanyVat (процент) и
// DefaultCompanyVatEnabled приходят только при UseCompanyVat = true; при «без НДС» сервер
// возвращает DefaultCompanyVatEnabled = false без DefaultCompanyVat. nil — поле не пришло.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#5-nastrojki-nds
type CompanyVat struct {
	UseCompanyVat            bool  `json:"useCompanyVat"`
	DefaultCompanyVat        *int  `json:"defaultCompanyVat,omitzero"`
	DefaultCompanyVatEnabled *bool `json:"defaultCompanyVatEnabled,omitzero"`
}

// CompanyVatRequest: указывается только при payerVat = true (проверяет сервер). При
// UseCompanyVat = false остальные поля сервер игнорирует; DefaultCompanyVat в 0–99 с
// DefaultCompanyVatEnabled = false сохраняет «без НДС».
type CompanyVatRequest struct {
	UseCompanyVat            model.Field[bool] `json:"useCompanyVat,omitzero"`
	DefaultCompanyVat        model.Field[int]  `json:"defaultCompanyVat,omitzero"`
	DefaultCompanyVatEnabled model.Field[bool] `json:"defaultCompanyVatEnabled,omitzero"`
}

// Image — подпись или печать юрлица. Meta.Href — ссылка на скачивание файла
// (transport.Client.Download), Miniature — метаданные миниатюры.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#5-podpisi-i-pechat
type Image struct {
	Meta      model.Meta     `json:"meta"`
	Title     string         `json:"title"`
	Filename  string         `json:"filename"`
	Size      int64          `json:"size"`
	Updated   model.DateTime `json:"updated,omitzero"`
	Miniature *model.Meta    `json:"miniature,omitzero"`
}

// ImageRequest — новый файл подписи или печати: имя с расширением и содержимое; JSON кодирует
// Content в Base64. Удаление (null) источник не описывает: Null отправляется как есть.
type ImageRequest struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}
