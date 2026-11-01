package employee

import "github.com/eteeet/moysklad-sdk-go/model"

// Модель чтения без expand. Права и роль сотрудника читаются отдельным ресурсом security.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/employee
type Employee struct {
	Meta         model.Meta     `json:"meta"`
	ID           string         `json:"id"`
	AccountID    string         `json:"accountId"`
	Created      model.DateTime `json:"created,omitzero"`
	Updated      model.DateTime `json:"updated,omitzero"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	ExternalCode string         `json:"externalCode"`
	Archived     bool           `json:"archived"`
	Shared       bool           `json:"shared"`

	// UID — логин сотрудника; name, fullName и shortFio сервер собирает из ФИО.
	UID        string `json:"uid"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	FullName   string `json:"fullName"`
	ShortFio   string `json:"shortFio"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	INN        string `json:"inn"`
	Position   string `json:"position"`

	Salary     *Salary           `json:"salary,omitzero"`
	Image      *Image            `json:"image,omitzero"`
	Cashiers   []Cashier         `json:"cashiers"`
	Attributes []model.Attribute `json:"attributes"`
	Owner      *model.Reference  `json:"owner,omitzero"`
	Group      *model.Reference  `json:"group,omitzero"`
}

type Salary struct {
	Value float64 `json:"value"` // В копейках (сотых долях валюты).
}

// Кассир точки продаж; только для чтения.
type Cashier struct {
	Meta        model.Meta       `json:"meta"`
	ID          string           `json:"id"`
	AccountID   string           `json:"accountId"`
	Employee    *model.Reference `json:"employee,omitzero"`
	RetailStore *model.Reference `json:"retailStore,omitzero"`
}

// Фотография сотрудника; загрузка фото в SDK не поддерживается.
type Image struct {
	Meta      model.Meta     `json:"meta"`
	Title     string         `json:"title"`
	Filename  string         `json:"filename"`
	Size      int            `json:"size"`
	Updated   model.DateTime `json:"updated,omitzero"`
	Miniature model.Meta     `json:"miniature,omitzero"`
	Tiny      model.Meta     `json:"tiny,omitzero"`
}
