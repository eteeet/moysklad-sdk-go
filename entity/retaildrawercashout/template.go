package retaildrawercashout

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/retaildrawercashout/new, документ не создаётся. basis — объект
// основания по источнику (розничная смена), например {"retailShift": model.Reference{…}};
// nil — шаблон без основания. Ответ — RetailDrawerCashOut без id:
// его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-shablon-vyplaty-deneg-na-osnove
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*RetailDrawerCashOut, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[RetailDrawerCashOut](ctx, s.client, "retaildrawercashout", basis, opts)
}
