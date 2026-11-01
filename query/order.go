package query

import (
	"fmt"
	"strings"
)

// Sort — одно поле сортировки. Сортируемые поля определяет endpoint.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-sorting#2-sortirovka
type Sort struct {
	field string
	desc  bool
}

func Asc(field string) Sort  { return Sort{field: field} }
func Desc(field string) Sort { return Sort{field: field, desc: true} }

// Order собирает значение order: «field,asc;other,desc» в порядке приоритета.
func Order(sorts ...Sort) (string, error) {
	parts := make([]string, 0, len(sorts))
	for _, s := range sorts {
		if s.field == "" || strings.ContainsAny(s.field, ",;") {
			return "", fmt.Errorf("moysklad: order field %q must be non-empty and contain no ',' or ';'", s.field)
		}
		direction := "asc"
		if s.desc {
			direction = "desc"
		}
		parts = append(parts, s.field+","+direction)
	}
	return strings.Join(parts, ";"), nil
}
