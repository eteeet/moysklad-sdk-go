// Package query собирает значения query-параметров JSON API: filter, order и опции списков.
// Поддержку конкретного поля и оператора определяет endpoint; пакет не хранит их перечень.
package query

import (
	"fmt"
	"slices"
	"strings"
)

// Condition — одно условие filter. Поле — имя поля, href доп. поля или любое имя, которое
// поддерживает endpoint; значение — строка, UUID или полный href ссылки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-filtraciya-vyborki-s-pomoshyu-parametra-filter
type Condition struct {
	field, op, value string
	raw              bool
}

func Eq(field, value string) Condition       { return Condition{field: field, op: "=", value: value} }
func NotEq(field, value string) Condition    { return Condition{field: field, op: "!=", value: value} }
func Less(field, value string) Condition     { return Condition{field: field, op: "<", value: value} }
func LessOrEq(field, value string) Condition { return Condition{field: field, op: "<=", value: value} }
func Greater(field, value string) Condition  { return Condition{field: field, op: ">", value: value} }
func GreaterOrEq(field, value string) Condition {
	return Condition{field: field, op: ">=", value: value}
}

// Contains, StartsWith и EndsWith — операторы подобия ~, ~= и =~.
func Contains(field, value string) Condition { return Condition{field: field, op: "~", value: value} }
func StartsWith(field, value string) Condition {
	return Condition{field: field, op: "~=", value: value}
}
func EndsWith(field, value string) Condition { return Condition{field: field, op: "=~", value: value} }

// IsEmpty отбирает отсутствующее или пустое значение (field=), NotEmpty — заполненное (field!=).
func IsEmpty(field string) Condition  { return Eq(field, "") }
func NotEmpty(field string) Condition { return NotEq(field, "") }

// RawCondition передаёт условие как есть, без экранирования и проверок,
// например для грамматики отчётов. Буквальная ';' в нём уже должна быть экранирована как «\;».
func RawCondition(text string) Condition { return Condition{value: text, raw: true} }

// Filter соединяет условия через ';' и экранирует ';' в значениях как «\;».
// Результат передаётся в url.Values как есть: URL-кодирование выполняется один раз при отправке.
//
// Повтор = одного поля даёт OR, повтор != — AND, разные поля — AND.
// Filter отклоняет то, что API объявляет ошибкой или будущей ошибкой: = вместе с неравенством
// того же поля и повтор <, <=, >, >= одного поля (сейчас API берёт первое значение).
// Значение с «\» в конце тоже отклоняется: вместе с разделителем оно дало бы «\;».
// Архив и корзина: Eq("archived", "true") и Eq("archived", "false") вместе отбирают все
// сущности; так же isDeleted для документов.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md
func Filter(conditions ...Condition) (string, error) {
	parts := make([]string, 0, len(conditions))
	ops := map[string]map[string]bool{}
	for _, c := range conditions {
		if c.raw {
			if c.value == "" || hasUnescapedSemicolon(c.value) {
				return "", fmt.Errorf("moysklad: raw filter condition must be non-empty and escape ';' as \\;")
			}
			parts = append(parts, c.value)
			continue
		}
		if c.field == "" || strings.Contains(c.field, ";") {
			return "", fmt.Errorf("moysklad: filter field %q must be non-empty and contain no ';'", c.field)
		}
		if strings.HasSuffix(c.value, `\`) {
			return "", fmt.Errorf("moysklad: filter value of %q must not end with a backslash", c.field)
		}
		seen := ops[c.field]
		if seen == nil {
			seen = map[string]bool{}
			ops[c.field] = seen
		}
		if conflicts(seen, c.op) {
			return "", fmt.Errorf("moysklad: filter field %q: %q conflicts with another condition on this field", c.field, c.op)
		}
		seen[c.op] = true
		parts = append(parts, c.field+c.op+strings.ReplaceAll(c.value, ";", `\;`))
	}
	return strings.Join(parts, ";"), nil
}

var comparisons = []string{"<", "<=", ">", ">="}

func conflicts(seen map[string]bool, op string) bool {
	switch op {
	case "=":
		return seen["!="] || slices.ContainsFunc(comparisons, func(c string) bool { return seen[c] })
	case "!=":
		return seen["="]
	case "<", "<=", ">", ">=":
		return seen["="] || seen[op]
	}
	return false
}

func hasUnescapedSemicolon(text string) bool {
	for i := range len(text) {
		if text[i] == ';' && (i == 0 || text[i-1] != '\\') {
			return true
		}
	}
	return false
}
