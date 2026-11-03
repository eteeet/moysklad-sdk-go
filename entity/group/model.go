package group

import "github.com/eteeet/moysklad-sdk-go/model"

// Group — отдел. Index — порядковый номер в списке отделов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/group
type Group struct {
	Meta      model.Meta `json:"meta"`
	ID        string     `json:"id"`
	AccountID string     `json:"accountId"`
	Name      string     `json:"name"`
	Index     int        `json:"index"`
}
