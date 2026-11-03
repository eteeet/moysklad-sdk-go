package facturein

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/facturein/new, документ не создаётся. Основания по источнику —
// приёмка или исходящий платёж; пример раздела передаёт их массивом связей, например
// {"supplies": []model.Reference{…}} или {"payments": […]}; nil — шаблон без основания. Ответ —
// FactureIn без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/facturein#3-shablon-scheta-faktury-poluchennogo
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*FactureIn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[FactureIn](ctx, s.client, "facturein", basis, opts)
}
