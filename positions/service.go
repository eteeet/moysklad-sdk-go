package positions

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

var readerDocumentTypePattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

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
// Metadata positions в заголовке документа не означает, что строки уже прочитаны.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2468-L2525
func (s *Service[T]) GetPage(ctx context.Context, documentID string, page PageOptions) (*model.Collection[T], error) {
	operation := "get " + s.documentType + " " + documentID + " positions"
	if !validID(documentID) {
		return nil, fmt.Errorf("moysklad: %s: document id must be a UUID", operation)
	}
	if page.Limit < 0 || page.Limit > 1000 {
		return nil, fmt.Errorf("moysklad: %s: limit must be between 1 and 1000, or 0 for default", operation)
	}
	if page.Offset < 0 {
		return nil, fmt.Errorf("moysklad: %s: offset must not be negative", operation)
	}
	limit := page.Limit
	if limit == 0 {
		limit = 1000
	}
	query := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(page.Offset)},
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet,
		"entity/"+s.documentType+"/"+documentID+"/positions", query, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var result model.Collection[T]
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	if result.Rows == nil {
		return nil, fmt.Errorf("moysklad: %s: response rows must be a non-null array", operation)
	}
	return &result, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
