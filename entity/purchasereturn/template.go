package purchasereturn

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/purchasereturn/new, документ не создаётся. basis — объект основания по
// источнику (приёмка), например {"supply": model.Reference{…}}; nil — шаблон без
// основания. Ответ — PurchaseReturn без id: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchase-return#3-shablon-vozvrata-postavshhiku
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*PurchaseReturn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[PurchaseReturn](ctx, s.client, "purchasereturn", basis, opts)
}
