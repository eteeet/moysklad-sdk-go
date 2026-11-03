package cashin

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/cashin/new, документ не создаётся. basis — объект основания по
// источнику (заказ покупателя, возврат поставщику, отгрузка, счёт покупателю, полученный отчёт комиссионера), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — CashIn без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin#3-shablon-prihodnogo-ordera
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*CashIn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[CashIn](ctx, s.client, "cashin", basis, opts)
}
