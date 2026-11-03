package paymentout

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/paymentout/new, документ не создаётся. basis — объект основания по
// источнику (возврат покупателя, приёмка, счёт поставщика, заказ поставщику, выданный отчёт комиссионера), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — PaymentOut без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out#3-shablon-ishodyashhego-platezha
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*PaymentOut, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[PaymentOut](ctx, s.client, "paymentout", basis, opts)
}
