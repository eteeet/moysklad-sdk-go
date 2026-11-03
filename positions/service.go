package positions

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"slices"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

var readerDocumentTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// positionFields — значения fields, которые источник документирует для позиций: stock — остатки
// и себестоимость (общий раздел), declaration — прослеживаемость у списания и розничной продажи.
var positionFields = map[string][]string{
	"demand": {"stock"}, "customerorder": {"stock"}, "invoiceout": {"stock"}, "purchaseorder": {"stock"},
	"invoicein": {"stock"}, "supply": {"stock"}, "purchasereturn": {"stock"}, "salesreturn": {"stock"},
	"retailsalesreturn": {"stock"}, "retaildemand": {"stock", "declaration"}, "loss": {"declaration"},
}

// Service читает отдельные страницы позиций документа в выбранную вызывающим кодом форму T.
// Generic reader не определяет полноту полей позиций конкретного документа.
type Service[T any] struct {
	client       *transport.Client
	documentType string
}

// NewService создаёт reader для API-кода документа вида [a-z][a-z0-9]*.
// Поддержку endpoint positions для этого кода определяет сервер.
func NewService[T any](client *transport.Client, documentType string) (*Service[T], error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: positions service requires a client")
	}
	if !readerDocumentTypePattern.MatchString(documentType) {
		return nil, fmt.Errorf("moysklad: positions document type must match [a-z][a-z0-9]*")
	}
	return &Service[T]{client: client, documentType: documentType}, nil
}

// GetPage читает одну страницу без повторов и переходов по href; limit/offset задаёт вызывающий код.
// fields принимается только из positionFields для этого кода документа.
// Metadata positions в заголовке документа не означает, что строки уже прочитаны.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2468-L2525
func (s *Service[T]) GetPage(ctx context.Context, documentID string, page PageOptions, opts ...transport.RequestOption) (*model.Collection[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "get " + s.documentType + " " + documentID + " positions"
	if !validID(documentID) {
		return nil, fmt.Errorf("moysklad: %s: document id must be a UUID", operation)
	}
	if page.Fields != "" && !slices.Contains(positionFields[s.documentType], page.Fields) {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	query, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[T](ctx, s.client, operation, "entity/"+s.documentType+"/"+documentID+"/positions", query, opts...)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}

// positionMetaTypes — meta.type строк, отличный от {документ}position: так в примерах источника.
var positionMetaTypes = map[string]model.MetaType{
	"invoicein": "invoiceposition", "invoiceout": "invoiceposition",
	"retaildemand": "demandposition", "retailsalesreturn": "salesreturnposition",
}

func (s *Service[T]) rows(documentID string) nested.Collection {
	metaType, ok := positionMetaTypes[s.documentType]
	if !ok {
		metaType = model.MetaType(s.documentType + "position")
	}
	return nested.Collection{Client: s.client, Parent: s.documentType, ParentID: documentID, Path: "positions", Item: "position", MetaType: metaType}
}

// GetPosition читает одну строку документа.
func (s *Service[T]) GetPosition(ctx context.Context, documentID, positionID string, opts ...transport.RequestOption) (*T, error) {
	return nested.Get[T](ctx, s.rows(documentID), positionID, opts)
}

// CreatePositions добавляет строки к существующим одним POST массивом (одна строка — массив из
// одного элемента); assortment и положительный quantity проверяет сервер. Ответ — по правилам
// batch; повторов нет.
func (s *Service[T]) CreatePositions(ctx context.Context, documentID string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[T], error) {
	return nested.CreateMany[T](ctx, s.rows(documentID), items, opts)
}

// UpdatePosition изменяет переданные поля одной строки; пустой запрос — ошибка до HTTP.
func (s *Service[T]) UpdatePosition(ctx context.Context, documentID, positionID string, input PositionRequest, opts ...transport.RequestOption) (*T, error) {
	return nested.Update[T](ctx, s.rows(documentID), positionID, input, opts)
}

// DeletePosition удаляет одну строку; повторов нет.
func (s *Service[T]) DeletePosition(ctx context.Context, documentID, positionID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.rows(documentID), positionID, opts)
}

// DeletePositions удаляет строки одним POST …/positions/delete с их meta. Источник отвечает 200 с
// пустым телом, поэтому результата по элементам нет; 1–1000 строк, повторов нет.
func (s *Service[T]) DeletePositions(ctx context.Context, documentID string, positionIDs []string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	rows := s.rows(documentID)
	operation := "delete " + s.documentType + " " + documentID + " positions"
	if !validID(documentID) {
		return fmt.Errorf("moysklad: %s: document id must be a UUID", operation)
	}
	if len(positionIDs) == 0 || len(positionIDs) > batch.MaxItems {
		return fmt.Errorf("moysklad: %s: ids must contain 1 to %d elements, got %d", operation, batch.MaxItems, len(positionIDs))
	}
	type item struct {
		Meta model.Meta `json:"meta"`
	}
	endpoint := "entity/" + s.documentType + "/" + documentID + "/positions"
	body := make([]item, len(positionIDs))
	for i, id := range positionIDs {
		if !validID(id) {
			return fmt.Errorf("moysklad: %s: ids[%d]: id must be a UUID", operation, i)
		}
		body[i] = item{model.Meta{Href: s.client.Href(endpoint + "/" + id), Type: rows.MetaType, MediaType: "application/json"}}
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
