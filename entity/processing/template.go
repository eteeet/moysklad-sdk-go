package processing

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон техоперации: PUT entity/processing/new, документ не создаётся. basis — основание по
// источнику: заказ на производство ({"processingOrder": model.Reference{…}}) или техкарта
// ({"processingPlan": model.Reference{…}}); nil — шаблон без основания. Ответ — Processing без id;
// materials/products в нём — только rows без meta. Нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#4-shablon-tehoperacii-na-osnove
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*Processing, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[Processing](ctx, s.client, "processing", basis, opts)
}
