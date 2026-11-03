package retaildrawercashin

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/retaildrawercashin/new, документ не создаётся. basis — объект
// основания по источнику (розничная смена), например {"retailShift": model.Reference{…}}
// (retailShift — с большой S); nil — шаблон без основания. Ответ — RetailDrawerCashIn без id:
// его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashin#3-shablon-vneseniya-deneg-na-osnove
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*RetailDrawerCashIn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[RetailDrawerCashIn](ctx, s.client, "retaildrawercashin", basis, opts)
}
