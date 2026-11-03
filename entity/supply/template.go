package supply

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/supply/new, документ не создаётся. basis — объект основания по
// источнику (без основания (раздел шаблона приёмки)), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — Supply без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply#3-shablon-priemki
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Supply, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Supply](ctx, s.client, "supply", basis, opts)
}
