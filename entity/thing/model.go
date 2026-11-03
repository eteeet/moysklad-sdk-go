package thing

import "github.com/eteeet/moysklad-sdk-go/model"

// Thing — серийный номер; только чтение (записи источник не описывает).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/thing
type Thing struct {
	Meta        model.Meta `json:"meta"`
	ID          string     `json:"id"`
	AccountID   string     `json:"accountId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}
