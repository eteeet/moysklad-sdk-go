package factureout

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/factureout/new, документ не создаётся. Основания по источнику —
// отгрузка, возврат поставщику или входящий платёж; примеры раздела передают их массивом связей:
// {"demands": []model.Reference{…}}, {"returns": […]} или {"payments": […]}; nil — шаблон без
// основания. Ответ — FactureOut без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-shablon-scheta-faktury-vydannogo
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*FactureOut, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[FactureOut](ctx, s.client, "factureout", basis, opts)
}
