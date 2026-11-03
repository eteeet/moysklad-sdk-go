package inventory

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/inventory/new, документ не создаётся. Источник описывает
// шаблон только «без связи с каким-либо документом» (basis nil); непустой basis передаётся как
// есть, и его принятие определяет сервер. Ответ — Inventory без id: его можно доработать и
// передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-shablon-inventarizacii
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Inventory, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Inventory](ctx, s.client, "inventory", basis, opts)
}
