package model

// Набор значений открыт: неизвестные типы сохраняются без валидации и значения по умолчанию.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#4-tip-kontragenta
type CompanyType string

const (
	CompanyTypeLegal          CompanyType = "legal"
	CompanyTypeEntrepreneur   CompanyType = "entrepreneur"
	CompanyTypeIndividual     CompanyType = "individual"
	CompanyTypeLegalUZ        CompanyType = "legalUZ"
	CompanyTypeEntrepreneurUZ CompanyType = "entrepreneurUZ"
	CompanyTypeIndividualUZ   CompanyType = "individualUZ"
	CompanyTypeLegalKZ        CompanyType = "legalKZ"
	CompanyTypeEntrepreneurKZ CompanyType = "entrepreneurKZ"
	CompanyTypeIndividualKZ   CompanyType = "individualKZ"
)
