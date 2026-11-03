package retaildemand

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/retaildemand/new, документ не создаётся. basis — объект основания по
// источнику (розничная смена, заказ покупателя), например {"retailShift": model.Reference{…}}; nil — шаблон без
// основания. Ответ — RetailDemand без id: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildemand#3-shablon-roznichnoj-prodazhi
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*RetailDemand, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[RetailDemand](ctx, s.client, "retaildemand", basis, opts)
}
