package processingorder

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/processingorder/new, документ не создаётся. basis — объект
// основания по источнику: техкарта, например {"processingPlan": model.Reference{…}}; nil — пустое
// тело (источник его не показывает, принятие определяет сервер). Ответ — ProcessingOrder без id,
// строки в positions.rows без meta: нужные поля переносятся в CreateRequest вручную.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-shablon-zakaza-na-proizvodstvo
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*ProcessingOrder, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[ProcessingOrder](ctx, s.client, "processingorder", basis, opts)
}
