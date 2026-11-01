package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// PageOptions — одна страница вложенной коллекции (позиции документов, компоненты комплекта).
// Limit=0 означает 1000; допустимы Limit от 1 до 1000 и неотрицательный Offset. Fields — одно
// значение; какие принимает коллекция, проверяет её метод.
type PageOptions struct {
	Limit  int
	Offset int
	Fields string
}

func (p PageOptions) Values() (url.Values, error) {
	if p.Limit < 0 || p.Limit > MaxLimit {
		return nil, fmt.Errorf("moysklad: limit must be between 1 and %d, or 0 for default", MaxLimit)
	}
	if p.Offset < 0 {
		return nil, fmt.Errorf("moysklad: offset must not be negative")
	}
	if strings.Contains(p.Fields, ",") {
		return nil, fmt.Errorf("moysklad: fields accepts a single value")
	}
	limit := p.Limit
	if limit == 0 {
		limit = MaxLimit
	}
	values := url.Values{"limit": {strconv.Itoa(limit)}, "offset": {strconv.Itoa(p.Offset)}}
	if p.Fields != "" {
		values.Set("fields", p.Fields)
	}
	return values, nil
}
