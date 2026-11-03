package paymentin

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Trash — удаление документа в корзину: POST entity/paymentin/{id}/trash. В отличие от Delete документ
// восстанавливается из корзины (только в интерфейсе); права и настройку корзины проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/common-info#3-udalenie-v-korzinu
func (s *Service) Trash(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	return write.Trash(ctx, s.client, "paymentin", id, opts)
}
