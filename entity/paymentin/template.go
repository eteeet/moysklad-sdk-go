package paymentin

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/paymentin/new, документ не создаётся. basis — объект основания по
// источнику (заказ покупателя, возврат поставщику, отгрузка, счёт покупателю, полученный отчёт комиссионера), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — PaymentIn без id: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in#3-shablon-vhodyashhego-platezha
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*PaymentIn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[PaymentIn](ctx, s.client, "paymentin", basis, opts)
}
