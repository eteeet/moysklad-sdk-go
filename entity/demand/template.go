package demand

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/demand/new, документ не создаётся. basis — объект основания по
// источнику (заказ покупателя (customerOrder), счета покупателю (invoicesOut)), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — Demand без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/demand#3-shablon-otgruzki-na-osnove
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Demand, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Demand](ctx, s.client, "demand", basis, opts)
}
