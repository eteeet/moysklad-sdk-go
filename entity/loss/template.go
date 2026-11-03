package loss

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/loss/new, документ не создаётся. basis — объект основания по
// источнику (возврат покупателя (salesReturn), инвентаризация (не больше 500 позиций)), например {"salesReturn": model.Reference{…}}; nil — шаблон без
// основания. Ответ — Loss без id: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/loss#3-shablon-spisaniya
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Loss, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Loss](ctx, s.client, "loss", basis, opts)
}
