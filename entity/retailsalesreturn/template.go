package retailsalesreturn

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/retailsalesreturn/new, документ не создаётся. basis — объект основания по
// источнику (розничная смена, розничная продажа), например {"retailShift": model.Reference{…}}; nil — шаблон без
// основания. Ответ — RetailSalesReturn без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retail-sales-return#3-shablon-roznichnogo-vozvrata
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*RetailSalesReturn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[RetailSalesReturn](ctx, s.client, "retailsalesreturn", basis, opts)
}
