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
}
