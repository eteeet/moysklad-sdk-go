// Package list читает одну страницу коллекции {context, meta, rows} без перехода по nextHref.
package list

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/model"
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
