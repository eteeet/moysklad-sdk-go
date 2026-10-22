package moysklad

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"strconv"
)

// Entity — операции с сущностью вида /entity/{type}.
type Entity[T any] struct {
	t    *transport
	path string
}

// ListOptions — параметры запроса списка.
type ListOptions struct {
	Filter Filter
	// Limit — размер страницы, до 1000. 0 — по умолчанию API.
	Limit  int
	Offset int
}

// Page — страница списка.
type Page[T any] struct {
	Meta PageMeta `json:"meta"`
	Rows []T      `json:"rows"`
}

// PageMeta — метаданные страницы списка.
type PageMeta struct {
	// Size — сколько всего записей подходит под фильтр.
	Size   int `json:"size"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	// NextHref — ссылка на следующую страницу; пустая на последней.
	NextHref string `json:"nextHref"`
}

// Get возвращает сущность по id.
func (e Entity[T]) Get(ctx context.Context, id string) (T, error) {
	var v T

	err := e.t.do(ctx, http.MethodGet, e.path+"/"+id, nil, nil, &v)
	if err != nil {
		return v, fmt.Errorf("transport.do: %w", err)
	}

	return v, nil
}

// List возвращает одну страницу списка.
func (e Entity[T]) List(ctx context.Context, opts ListOptions) (Page[T], error) {
	q := url.Values{}
	if f := opts.Filter.String(); f != "" {
		q.Set("filter", f)
	}

	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}

	if opts.Offset > 0 {
		q.Set("offset", strconv.Itoa(opts.Offset))
	}

	var p Page[T]

	err := e.t.do(ctx, http.MethodGet, e.path, q, nil, &p)
	if err != nil {
		return p, fmt.Errorf("transport.do: %w", err)
	}

	return p, nil
}

// maxLimit — наибольший размер страницы списка в API.
const maxLimit = 1000

// All обходит все страницы списка, начиная с opts.Offset. Limit 0 или больше 1000 — страницы по 1000 записей.
// Ошибка приходит последним элементом, после неё обход заканчивается.
func (e Entity[T]) All(ctx context.Context, opts ListOptions) iter.Seq2[T, error] {
	if opts.Limit <= 0 || opts.Limit > maxLimit {
		opts.Limit = maxLimit
	}

	return func(yield func(T, error) bool) {
		for {
			page, err := e.List(ctx, opts)
			if err != nil {
				var zero T

				yield(zero, fmt.Errorf("Entity.List: %w", err))

				return
			}

			for _, row := range page.Rows {
				if !yield(row, nil) {
					return
				}
			}

			if page.Meta.NextHref == "" {
				return
			}

			opts.Offset += opts.Limit
		}
	}
}
