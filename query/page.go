package query

import (
	"fmt"
	"net/url"
	"strconv"
)

// PageOptions — одна страница вложенной коллекции, у которой источник объявляет только limit
// и offset (позиции документов, компоненты комплекта). Limit=0 означает 1000; допустимы Limit
// от 1 до 1000 и неотрицательный Offset.
type PageOptions struct {
	Limit  int
	Offset int
}

func (p PageOptions) Values() (url.Values, error) {
	if p.Limit < 0 || p.Limit > MaxLimit {
		return nil, fmt.Errorf("moysklad: limit must be between 1 and %d, or 0 for default", MaxLimit)
	}
	if p.Offset < 0 {
		return nil, fmt.Errorf("moysklad: offset must not be negative")
	}
	limit := p.Limit
	if limit == 0 {
		limit = MaxLimit
	}
	return url.Values{"limit": {strconv.Itoa(limit)}, "offset": {strconv.Itoa(p.Offset)}}, nil
}
