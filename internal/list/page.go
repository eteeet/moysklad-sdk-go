// Package list читает одну страницу коллекции {context, meta, rows} без перехода по nextHref.
package list

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Page выполняет один GET. Голые массивы и ответы отчётов — другие контракты, не Collection.
func Page[T any](ctx context.Context, client *transport.Client, operation, endpoint string, params url.Values,
	opts ...transport.RequestOption) (*model.Collection[T], error) {
	var result model.Collection[T]
	resp, err := client.Send(ctx, http.MethodGet, endpoint, params, nil, &result, opts...)
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

// Entity читает одну страницу /entity/{entity}. fields — значения fields, которые List этого
// endpoint принимает; остальные отклоняются до запроса. fields=stock у списков документов не
// принимается: остатки лежат в positions и в списке появляются только с expand=positions,
// а expand вне SDK; строки с остатками читает GetPositions.
func Entity[T any](ctx context.Context, client *transport.Client, entity string, options query.ListOptions,
	fields []string, opts ...transport.RequestOption) (*model.Collection[T], error) {
	operation := "list " + entity
	if options.Fields != "" && !slices.Contains(fields, options.Fields) {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, options.Fields)
	}
	params, err := options.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return Page[T](ctx, client, operation, "entity/"+entity, params, opts...)
}
