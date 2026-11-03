package publication

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest — публикация документа по шаблону печатной формы (embeddedtemplate или
// customtemplate); template обязателен. Href только для чтения и не передаётся.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
type CreateRequest struct {
	Template model.Reference `json:"template"`
}
