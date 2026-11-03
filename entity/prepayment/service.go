package prepayment

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

// Service — предоплаты: чтение, удаление и корзина. Создания, изменения и массовых операций
// источник не описывает.
type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: prepayment service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment#3-poluchit-predoplatu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Prepayment, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get prepayment: id must be a UUID")
	}
	operation := "get prepayment " + id
	var result Prepayment
	resp, err := s.client.Send(ctx, http.MethodGet, "entity/prepayment/"+id, nil, nil, &result, opts...)
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment#3-poluchit-spisok-predoplat
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Prepayment], error) {
	return list.Entity[Prepayment](ctx, s.client, "prepayment", options, nil, opts...)
}

// Delete — «Удалить Предоплату»: DELETE entity/prepayment/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment#3-udalit-predoplatu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete prepayment: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete prepayment "+id, "entity/prepayment/"+id, opts)
}

// Trash — удаление документа в корзину: POST entity/prepayment/{id}/trash. В отличие от Delete
// документ восстанавливается из корзины (только в интерфейсе); права и настройку корзины
// проверяет сервер. Раздела в _prepayment.md нет — общее правило документов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/common-info#3-udalenie-v-korzinu
func (s *Service) Trash(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	return write.Trash(ctx, s.client, "prepayment", id, opts)
}

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "prepayment")
	return reader
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields не поддержан.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment#3-poluchit-pozicii-predoplaty
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment#3-poluchit-poziciyu-predoplaty
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
