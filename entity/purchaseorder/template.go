package purchaseorder

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/purchaseorder/new, документ не создаётся. basis — объект основания по
// источнику (внутренний заказ, заказ покупателя), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — PurchaseOrder без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder#3-shablon-zakaza-postavshhiku
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*PurchaseOrder, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[PurchaseOrder](ctx, s.client, "purchaseorder", basis, opts)
}
