package organization

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-sozdat-yurlico
type CreateRequest struct {
	Name         string                         `json:"name"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Archived     model.Field[bool]              `json:"archived,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	CompanyType  model.Field[model.CompanyType] `json:"companyType,omitzero"`
	Email        model.Field[string]            `json:"email,omitzero"`
	Phone        model.Field[string]            `json:"phone,omitzero"`
	Fax          model.Field[string]            `json:"fax,omitzero"`

	LegalTitle        model.Field[string]         `json:"legalTitle,omitzero"`
	LegalFirstName    model.Field[string]         `json:"legalFirstName,omitzero"`
	LegalLastName     model.Field[string]         `json:"legalLastName,omitzero"`
	LegalMiddleName   model.Field[string]         `json:"legalMiddleName,omitzero"`
	INN               model.Field[string]         `json:"inn,omitzero"`
	KPP               model.Field[string]         `json:"kpp,omitzero"`
	OGRN              model.Field[string]         `json:"ogrn,omitzero"`
	OGRNIP            model.Field[string]         `json:"ogrnip,omitzero"`
	OKPO              model.Field[string]         `json:"okpo,omitzero"`
	CertificateNumber model.Field[string]         `json:"certificateNumber,omitzero"`
	CertificateDate   model.Field[model.DateTime] `json:"certificateDate,omitzero"`

	ActualAddress     model.Field[string]               `json:"actualAddress,omitzero"`
	LegalAddress      model.Field[string]               `json:"legalAddress,omitzero"`
	ActualAddressFull model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	LegalAddressFull  model.Field[model.AddressRequest] `json:"legalAddressFull,omitzero"`
	Attributes        model.Field[[]model.Attribute]    `json:"attributes,omitzero"`
	Owner             model.Field[model.Reference]      `json:"owner,omitzero"`
	Group             model.Field[model.Reference]      `json:"group,omitzero"`
	BonusProgram      model.Field[model.Reference]      `json:"bonusProgram,omitzero"`
	SyncID            model.Field[string]               `json:"syncId,omitzero"`

	PayerVAT               model.Field[bool]           `json:"payerVat,omitzero"`
	Director               model.Field[string]         `json:"director,omitzero"`
	DirectorPosition       model.Field[string]         `json:"directorPosition,omitzero"`
	ChiefAccountant        model.Field[string]         `json:"chiefAccountant,omitzero"`
	TrackingContractNumber model.Field[string]         `json:"trackingContractNumber,omitzero"`
	TrackingContractDate   model.Field[model.DateTime] `json:"trackingContractDate,omitzero"`

	CompanyVat        model.Field[CompanyVatRequest]         `json:"companyVat__ru,omitzero"`
	AdvancePaymentVat model.Field[float64]                   `json:"advancePaymentVat,omitzero"`
	DirectorSign      model.Field[ImageRequest]              `json:"directorSign,omitzero"`
	ChiefAccountSign  model.Field[ImageRequest]              `json:"chiefAccountSign,omitzero"`
	Stamp             model.Field[ImageRequest]              `json:"stamp,omitzero"`
	RequisitesUZ      model.Field[model.RequisitesUZRequest] `json:"mod__requisites__uz,omitzero"`
	RequisitesKZ      model.Field[model.RequisitesKZRequest] `json:"mod__requisites__kz,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/organization#3-izmenit-yurlico
type UpdateRequest struct {
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Archived     model.Field[bool]              `json:"archived,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	CompanyType  model.Field[model.CompanyType] `json:"companyType,omitzero"`
	Email        model.Field[string]            `json:"email,omitzero"`
	Phone        model.Field[string]            `json:"phone,omitzero"`
	Fax          model.Field[string]            `json:"fax,omitzero"`

	LegalTitle        model.Field[string]         `json:"legalTitle,omitzero"`
	LegalFirstName    model.Field[string]         `json:"legalFirstName,omitzero"`
	LegalLastName     model.Field[string]         `json:"legalLastName,omitzero"`
	LegalMiddleName   model.Field[string]         `json:"legalMiddleName,omitzero"`
	INN               model.Field[string]         `json:"inn,omitzero"`
	KPP               model.Field[string]         `json:"kpp,omitzero"`
	OGRN              model.Field[string]         `json:"ogrn,omitzero"`
	OGRNIP            model.Field[string]         `json:"ogrnip,omitzero"`
	OKPO              model.Field[string]         `json:"okpo,omitzero"`
	CertificateNumber model.Field[string]         `json:"certificateNumber,omitzero"`
	CertificateDate   model.Field[model.DateTime] `json:"certificateDate,omitzero"`

	ActualAddress     model.Field[string]               `json:"actualAddress,omitzero"`
	LegalAddress      model.Field[string]               `json:"legalAddress,omitzero"`
	ActualAddressFull model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	LegalAddressFull  model.Field[model.AddressRequest] `json:"legalAddressFull,omitzero"`
	Attributes        model.Field[[]model.Attribute]    `json:"attributes,omitzero"`
	Owner             model.Field[model.Reference]      `json:"owner,omitzero"`
	Group             model.Field[model.Reference]      `json:"group,omitzero"`
	BonusProgram      model.Field[model.Reference]      `json:"bonusProgram,omitzero"`
	SyncID            model.Field[string]               `json:"syncId,omitzero"`

	PayerVAT               model.Field[bool]           `json:"payerVat,omitzero"`
	Director               model.Field[string]         `json:"director,omitzero"`
	DirectorPosition       model.Field[string]         `json:"directorPosition,omitzero"`
	ChiefAccountant        model.Field[string]         `json:"chiefAccountant,omitzero"`
	TrackingContractNumber model.Field[string]         `json:"trackingContractNumber,omitzero"`
	TrackingContractDate   model.Field[model.DateTime] `json:"trackingContractDate,omitzero"`

	CompanyVat        model.Field[CompanyVatRequest]         `json:"companyVat__ru,omitzero"`
	AdvancePaymentVat model.Field[float64]                   `json:"advancePaymentVat,omitzero"`
	DirectorSign      model.Field[ImageRequest]              `json:"directorSign,omitzero"`
	ChiefAccountSign  model.Field[ImageRequest]              `json:"chiefAccountSign,omitzero"`
	Stamp             model.Field[ImageRequest]              `json:"stamp,omitzero"`
	RequisitesUZ      model.Field[model.RequisitesUZRequest] `json:"mod__requisites__uz,omitzero"`
	RequisitesKZ      model.Field[model.RequisitesKZRequest] `json:"mod__requisites__kz,omitzero"`
}
