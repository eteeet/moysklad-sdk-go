package publication

import "github.com/eteeet/moysklad-sdk-go/model"

// Publication — опубликованная печатная форма документа. Href — адрес страницы публикации,
// только чтение; Template — шаблон печати (без expand). Meta.Type в примерах — operationpublication.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/publication
type Publication struct {
	Meta     model.Meta      `json:"meta"`
	Template model.Reference `json:"template"`
	Href     string          `json:"href"`
}
