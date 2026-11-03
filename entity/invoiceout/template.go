package invoiceout

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// New — шаблон документа: PUT entity/invoiceout/new, документ не создаётся. basis — объект основания по
// источнику (заказ покупателя, отгрузка), например {"customerOrder": model.Reference{…}}; nil — шаблон без
// основания. Ответ — InvoiceOut без id: его можно доработать и передать в Create.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-out#3-shablon-scheta-pokupatelyu
func (s *Service) New(ctx context.Context, basis map[string]any, opts ...transport.RequestOption) (*InvoiceOut, error) {
	transport.ResetResponseInfo(opts)
	return write.Template[InvoiceOut](ctx, s.client, "invoiceout", basis, opts)
}
