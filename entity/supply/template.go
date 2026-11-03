package supply

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/supply/new, документ не создаётся. basis — объект основания по
// источнику (без основания (раздел шаблона приёмки)); nil — шаблон без
// основания. Ответ — Supply без id: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply#3-shablon-priemki
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Supply, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Supply](ctx, s.client, "supply", basis, opts)
}
