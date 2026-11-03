package internalorder

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/internalorder/new, документ не создаётся. Источник
// показывает только пустое тело (basis nil); основания для внутреннего заказа он не описывает,
// непустой basis передаётся как есть, и его принятие определяет сервер. Ответ — InternalOrder
// без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-shablon-vnutrennego-zakaza
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*InternalOrder, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[InternalOrder](ctx, s.client, "internalorder", basis, opts)
}
