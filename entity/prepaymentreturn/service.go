package prepaymentreturn

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Service — возвраты предоплаты: чтение, удаление и корзина. Создания, изменения и массовых операций
// источник не описывает.
type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: prepaymentreturn service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment-return#3-poluchit-vozvrat-predoplaty
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*PrepaymentReturn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get prepaymentreturn: id must be a UUID")
	}
	operation := "get prepaymentreturn " + id
	var result PrepaymentReturn
	resp, err := s.client.Send(ctx, http.MethodGet, "entity/prepaymentreturn/"+id, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment-return#3-poluchit-spisok-vozvratov-predoplaty
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[PrepaymentReturn], error) {
	return list.Entity[PrepaymentReturn](ctx, s.client, "prepaymentreturn", options, nil, opts...)
}

// Delete — «Удалить Возврат предоплаты»: DELETE entity/prepaymentreturn/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment-return#3-udalit-vozvrat-predoplaty
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete prepaymentreturn: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete prepaymentreturn "+id, "entity/prepaymentreturn/"+id, opts)
}

// Trash — удаление документа в корзину: POST entity/prepaymentreturn/{id}/trash. В отличие от Delete
// документ восстанавливается из корзины (только в интерфейсе); права и настройку корзины
// проверяет сервер. Раздела в _prepaymentreturn.md нет — общее правило документов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/common-info#3-udalenie-v-korzinu
func (s *Service) Trash(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	return write.Trash(ctx, s.client, "prepaymentreturn", id, opts)
}

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "prepaymentreturn")
	return reader
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields не поддержан.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment-return#3-poluchit-pozicii-vozvrata-predoplaty
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment-return#3-poluchit-poziciyu-vozvrata-predoplaty
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().GetPosition(ctx, id, positionID, opts...)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
