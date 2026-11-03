package salesreturn

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/salesreturn/new, документ не создаётся. basis — объект основания по
// источнику (отгрузка, розничная продажа), например {"demand": model.Reference{…}}; nil — шаблон без
// основания. Ответ — SalesReturn без id: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-shablon-vozvrata-pokupatelya-na-osnove
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*SalesReturn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[SalesReturn](ctx, s.client, "salesreturn", basis, opts)
}
