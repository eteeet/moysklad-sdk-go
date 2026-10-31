package counterparty

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#4-atributy-sushnosti
type Counterparty struct {
	Meta               model.Meta        `json:"meta"`
	ID                 string            `json:"id"`
	AccountID          string            `json:"accountId"`
	Created            model.DateTime    `json:"created,omitzero"`
	Updated            model.DateTime    `json:"updated,omitzero"`
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Code               string            `json:"code"`
	ExternalCode       string            `json:"externalCode"`
	Archived           bool              `json:"archived"`
	Shared             bool              `json:"shared"`
	CompanyType        model.CompanyType `json:"companyType"`
	Email              string            `json:"email"`
	Phone              string            `json:"phone"`
	Fax                string            `json:"fax"`
	LegalTitle         string            `json:"legalTitle"`
	LegalFirstName     string            `json:"legalFirstName"`
	LegalLastName      string            `json:"legalLastName"`
	LegalMiddleName    string            `json:"legalMiddleName"`
	INN                string            `json:"inn"`
	KPP                string            `json:"kpp"`
	OGRN               string            `json:"ogrn"`
	OGRNIP             string            `json:"ogrnip"`
	OKPO               string            `json:"okpo"`
	CertificateNumber  string            `json:"certificateNumber"`
	CertificateDate    model.DateTime    `json:"certificateDate,omitzero"`
	BirthDate          model.DateTime    `json:"birthDate,omitzero"`
	Sex                string            `json:"sex"`
	ActualAddress      string            `json:"actualAddress"`
	LegalAddress       string            `json:"legalAddress"`
	ActualAddressFull  *model.Address    `json:"actualAddressFull,omitzero"`
	LegalAddressFull   *model.Address    `json:"legalAddressFull,omitzero"`
	Attributes         []model.Attribute `json:"attributes"`
	Tags               []string          `json:"tags"`
	DiscountCardNumber string            `json:"discountCardNumber"`
	Owner              *model.Reference  `json:"owner,omitzero"`
	Group              *model.Reference  `json:"group,omitzero"`
	BonusProgram       *model.Reference  `json:"bonusProgram,omitzero"`
	State              *model.Reference  `json:"state,omitzero"`
	PriceType          *model.PriceType  `json:"priceType,omitzero"`
	SyncID             string            `json:"syncId"`
	BonusPoints        int               `json:"bonusPoints"`
	SalesAmount        float64           `json:"salesAmount"`
}
