package counterparty

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-kontragenta
type CreateRequest struct {
	Name               string                            `json:"name"`
	Description        model.Field[string]               `json:"description,omitzero"`
	Code               model.Field[string]               `json:"code,omitzero"`
	ExternalCode       model.Field[string]               `json:"externalCode,omitzero"`
	Archived           model.Field[bool]                 `json:"archived,omitzero"`
	Shared             model.Field[bool]                 `json:"shared,omitzero"`
	CompanyType        model.Field[model.CompanyType]    `json:"companyType,omitzero"`
	Email              model.Field[string]               `json:"email,omitzero"`
	Phone              model.Field[string]               `json:"phone,omitzero"`
	Fax                model.Field[string]               `json:"fax,omitzero"`
	LegalTitle         model.Field[string]               `json:"legalTitle,omitzero"`
	LegalFirstName     model.Field[string]               `json:"legalFirstName,omitzero"`
	LegalLastName      model.Field[string]               `json:"legalLastName,omitzero"`
	LegalMiddleName    model.Field[string]               `json:"legalMiddleName,omitzero"`
	INN                model.Field[string]               `json:"inn,omitzero"`
	KPP                model.Field[string]               `json:"kpp,omitzero"`
	OGRN               model.Field[string]               `json:"ogrn,omitzero"`
	OGRNIP             model.Field[string]               `json:"ogrnip,omitzero"`
	OKPO               model.Field[string]               `json:"okpo,omitzero"`
	CertificateNumber  model.Field[string]               `json:"certificateNumber,omitzero"`
	CertificateDate    model.Field[model.DateTime]       `json:"certificateDate,omitzero"`
	BirthDate          model.Field[model.DateTime]       `json:"birthDate,omitzero"`
	Sex                model.Field[string]               `json:"sex,omitzero"`
	ActualAddress      model.Field[string]               `json:"actualAddress,omitzero"`
	LegalAddress       model.Field[string]               `json:"legalAddress,omitzero"`
	ActualAddressFull  model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	LegalAddressFull   model.Field[model.AddressRequest] `json:"legalAddressFull,omitzero"`
	Attributes         model.Field[[]model.Attribute]    `json:"attributes,omitzero"`
	Tags               model.Field[[]string]             `json:"tags,omitzero"`
	DiscountCardNumber model.Field[string]               `json:"discountCardNumber,omitzero"`
	Owner              model.Field[model.Reference]      `json:"owner,omitzero"`
	Group              model.Field[model.Reference]      `json:"group,omitzero"`
	BonusProgram       model.Field[model.Reference]      `json:"bonusProgram,omitzero"`
	State              model.Field[model.Reference]      `json:"state,omitzero"`
	PriceType          model.Field[model.PriceType]      `json:"priceType,omitzero"`
	SyncID             model.Field[string]               `json:"syncId,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-kontragenta
type UpdateRequest struct {
	Name               model.Field[string]               `json:"name,omitzero"`
	Description        model.Field[string]               `json:"description,omitzero"`
	Code               model.Field[string]               `json:"code,omitzero"`
	ExternalCode       model.Field[string]               `json:"externalCode,omitzero"`
	Archived           model.Field[bool]                 `json:"archived,omitzero"`
	Shared             model.Field[bool]                 `json:"shared,omitzero"`
	CompanyType        model.Field[model.CompanyType]    `json:"companyType,omitzero"`
	Email              model.Field[string]               `json:"email,omitzero"`
	Phone              model.Field[string]               `json:"phone,omitzero"`
	Fax                model.Field[string]               `json:"fax,omitzero"`
	LegalTitle         model.Field[string]               `json:"legalTitle,omitzero"`
	LegalFirstName     model.Field[string]               `json:"legalFirstName,omitzero"`
	LegalLastName      model.Field[string]               `json:"legalLastName,omitzero"`
	LegalMiddleName    model.Field[string]               `json:"legalMiddleName,omitzero"`
	INN                model.Field[string]               `json:"inn,omitzero"`
	KPP                model.Field[string]               `json:"kpp,omitzero"`
	OGRN               model.Field[string]               `json:"ogrn,omitzero"`
	OGRNIP             model.Field[string]               `json:"ogrnip,omitzero"`
	OKPO               model.Field[string]               `json:"okpo,omitzero"`
	CertificateNumber  model.Field[string]               `json:"certificateNumber,omitzero"`
	CertificateDate    model.Field[model.DateTime]       `json:"certificateDate,omitzero"`
	BirthDate          model.Field[model.DateTime]       `json:"birthDate,omitzero"`
	Sex                model.Field[string]               `json:"sex,omitzero"`
	ActualAddress      model.Field[string]               `json:"actualAddress,omitzero"`
	LegalAddress       model.Field[string]               `json:"legalAddress,omitzero"`
	ActualAddressFull  model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	LegalAddressFull   model.Field[model.AddressRequest] `json:"legalAddressFull,omitzero"`
	Attributes         model.Field[[]model.Attribute]    `json:"attributes,omitzero"`
	Tags               model.Field[[]string]             `json:"tags,omitzero"`
	DiscountCardNumber model.Field[string]               `json:"discountCardNumber,omitzero"`
	Owner              model.Field[model.Reference]      `json:"owner,omitzero"`
	Group              model.Field[model.Reference]      `json:"group,omitzero"`
	BonusProgram       model.Field[model.Reference]      `json:"bonusProgram,omitzero"`
	State              model.Field[model.Reference]      `json:"state,omitzero"`
	PriceType          model.Field[model.PriceType]      `json:"priceType,omitzero"`
}
