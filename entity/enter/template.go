package enter

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/enter/new, документ не создаётся. basis — объект основания по
// источнику (инвентаризация (inventory; не больше 500 подходящих позиций)), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — Enter без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/enter
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Enter, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Enter](ctx, s.client, "enter", basis, opts)
}
