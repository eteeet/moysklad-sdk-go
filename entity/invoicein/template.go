package invoicein

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/invoicein/new, документ не создаётся. basis — объект основания по
// источнику (заказ поставщику, приёмка), например {"purchaseOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — InvoiceIn без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in#3-shablon-scheta-postavshhika
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*InvoiceIn, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[InvoiceIn](ctx, s.client, "invoicein", basis, opts)
}
