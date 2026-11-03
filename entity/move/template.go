package move

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/move/new, документ не создаётся. basis — объект основания по
// источнику (внутренний заказ, заказ покупателя), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — Move без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/move#3-shablon-peremeshheniya
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Move, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Move](ctx, s.client, "move", basis, opts)
}
