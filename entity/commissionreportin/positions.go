package commissionreportin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Два набора строк: «реализовано комиссионером» — entity/commissionreportin/{id}/positions,
// «возврат на склад комиссионера» — …/{id}/returntocommissionerpositions. Повторов нет.

const returnPath = "returntocommissionerpositions"

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "commissionreportin")
	return reader
}

func (s *Service) sold(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "commissionreportin", ParentID: id, Path: "positions", Item: "position", MetaType: "commissionreportinposition"}
}

func (s *Service) returned(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "commissionreportin", ParentID: id, Path: returnPath, Item: "return position", MetaType: "commissionreportinreturnedposition"}
}

// GetPositions читает одну страницу строк «реализовано» без скрытого Get и автообхода;
// fields источник для этого документа не документирует.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-poluchit-pozicii-realizovano-komissionerom
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-poluchit-poziciyu-realizovano-komissionerom
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().GetPosition(ctx, id, positionID, opts...)
}

// CreatePositions добавляет строки «реализовано» к существующим одним POST массивом (одна
// строка — массив из одного элемента); ответ — по правилам batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-sozdat-poziciyu-realizovano-komissionerom
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	return nested.CreateMany[Position](ctx, s.sold(id), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-izmenit-poziciyu-realizovano-komissionerom
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.sold(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-udalit-poziciyu-realizovano-komissionerom
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return s.rows().DeletePosition(ctx, id, positionID, opts...)
}

// DeleteManyPositions — POST …/positions/delete; ответ без тела, результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-massovoe-udalenie-pozicij-realizovano-komissionerom
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return s.rows().DeletePositions(ctx, id, positionIDs, opts...)
}

// GetReturnPositions читает одну страницу строк «возврат на склад комиссионера»: только limit и
// offset, без скрытого Get и автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-poluchit-pozicii-vozvrata-na-sklad-komissionera
func (s *Service) GetReturnPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[ReturnPosition], error) {
	transport.ResetResponseInfo(opts)
	operation := "get commissionreportin " + id + " return positions"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: document id must be a UUID", operation)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[ReturnPosition](ctx, s.client, operation, "entity/commissionreportin/"+id+"/"+returnPath, params, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-poluchit-poziciyu-vozvrata-na-sklad-komissionera
func (s *Service) GetReturnPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*ReturnPosition, error) {
	return nested.Get[ReturnPosition](ctx, s.returned(id), positionID, opts)
}

// CreateReturnPositions добавляет строки возврата к существующим одним POST массивом; ответ — по
// правилам batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-sozdat-poziciyu-vozvrata-na-sklad-komissionera
func (s *Service) CreateReturnPositions(ctx context.Context, id string, items []ReturnPositionRequest, opts ...transport.RequestOption) ([]batch.Result[ReturnPosition], error) {
	return nested.CreateMany[ReturnPosition](ctx, s.returned(id), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-izmenit-poziciyu-vozvrata-na-sklad-komissionera
func (s *Service) UpdateReturnPosition(ctx context.Context, id, positionID string, input ReturnPositionRequest, opts ...transport.RequestOption) (*ReturnPosition, error) {
	return nested.Update[ReturnPosition](ctx, s.returned(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-udalit-poziciyu-vozvrata-na-sklad-komissionera
func (s *Service) DeleteReturnPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.returned(id), positionID, opts)
}

// DeleteManyReturnPositions — POST …/returntocommissionerpositions/delete с meta строк, 1–1000 id.
// Источник отвечает 200 с пустым телом, поэтому результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-massovoe-udalenie-pozicij-vozvrata-na-sklad-komissionera
func (s *Service) DeleteManyReturnPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete commissionreportin " + id + " return positions"
	if !validID(id) {
		return fmt.Errorf("moysklad: %s: document id must be a UUID", operation)
	}
	if len(positionIDs) == 0 || len(positionIDs) > batch.MaxItems {
		return fmt.Errorf("moysklad: %s: ids must contain 1 to %d elements, got %d", operation, batch.MaxItems, len(positionIDs))
	}
	type item struct {
		Meta model.Meta `json:"meta"`
	}
	endpoint := "entity/commissionreportin/" + id + "/" + returnPath
	body := make([]item, len(positionIDs))
	for i, positionID := range positionIDs {
		if !validID(positionID) {
			return fmt.Errorf("moysklad: %s: ids[%d]: id must be a UUID", operation, i)
		}
		body[i] = item{model.Meta{Href: s.client.Href(endpoint + "/" + positionID), Type: "commissionreportinreturnedposition", MediaType: "application/json"}}
	}
	resp, err := s.client.Send(ctx, http.MethodPost, endpoint+"/delete", nil, body, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}
