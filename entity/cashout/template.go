package cashout

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/cashout/new, документ не создаётся. basis — объект основания по
// источнику (возврат покупателя, приёмка, счёт поставщика, заказ поставщику, выданный отчёт комиссионера), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — CashOut без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashout#3-shablon-rashodnogo-ordera
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*CashOut, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[CashOut](ctx, s.client, "cashout", basis, opts)
}
