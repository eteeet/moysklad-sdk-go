// Package pages явно обходит страницы коллекции {context, meta, rows} по meta.nextHref.
// Обычные List/GetPositions читают одну страницу; обход — отдельное решение вызывающего кода.
package pages

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"strconv"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// All отдаёт first и следующие страницы по одной, держа в памяти только текущую.
// Каждая следующая страница — отдельный GET по meta.nextHref, если он ведёт на тот же origin
// и путь коллекции, что и first.Meta.Href, а его offset больше offset текущей страницы;
// offset не вычисляется. Отмена ctx проверяется перед каждой страницей.
// Набор может меняться между запросами: пропуски и дубли строк возможны, снимка нет,
// страница может оказаться пустой.
// opts применяются к каждому запросу; ResponseInfo описывает последнюю загруженную страницу.
func All[T any](ctx context.Context, client *transport.Client, first *model.Collection[T],
	opts ...transport.RequestOption) iter.Seq2[*model.Collection[T], error] {
	return func(yield func(*model.Collection[T], error) bool) {
		if first == nil {
			yield(nil, fmt.Errorf("moysklad: pages: first page is nil"))
			return
		}
		base, err := url.Parse(first.Meta.Href)
		if err != nil || !base.IsAbs() || base.Path == "" {
			yield(nil, fmt.Errorf("moysklad: pages: first page has invalid meta.href %q", first.Meta.Href))
			return
		}
		offset, err := offsetOf(base)
		if err != nil {
			yield(nil, err)
			return
		}
		page := first
		for {
			if err := ctx.Err(); err != nil {
				yield(nil, fmt.Errorf("moysklad: pages: %w", err))
				return
			}
			if !yield(page, nil) || page.Meta.NextHref == "" {
				return
			}
			next, err := url.Parse(page.Meta.NextHref)
			if err != nil || !next.IsAbs() || next.Scheme != base.Scheme || next.Host != base.Host || next.Path != base.Path {
				yield(nil, fmt.Errorf("moysklad: pages: nextHref %q leaves the collection %s", page.Meta.NextHref, base.Path))
				return
			}
			nextOffset, err := offsetOf(next)
			if err != nil {
				yield(nil, err)
				return
			}
			if nextOffset <= offset {
				yield(nil, fmt.Errorf("moysklad: pages: nextHref offset %d does not advance past %d", nextOffset, offset))
				return
			}
			if err := ctx.Err(); err != nil {
				yield(nil, fmt.Errorf("moysklad: pages: %w", err))
				return
			}
			page, err = list.Page[T](ctx, client, "pages: get next page", page.Meta.NextHref, nil, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			offset = nextOffset
		}
	}
}

// Отсутствующий offset — начало коллекции.
func offsetOf(href *url.URL) (int, error) {
	value := href.Query().Get("offset")
	if value == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("moysklad: pages: invalid offset %q in %s", value, href)
	}
	return n, nil
}
