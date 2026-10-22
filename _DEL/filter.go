package moysklad

import (
	"slices"
	"strings"
	"time"
)

// Filter собирает параметр filter МойСклад. Условия соединяются через ";":
// условия на разные поля — по И, равенства на одно поле — по ИЛИ.
// Например, archived=true;archived=false отдаёт и архивные, и обычные записи.
type Filter struct{ conds []string }

// Eq добавляет условие field=value.
func (f Filter) Eq(field, value string) Filter { return f.add(field + "=" + value) }

// Gte добавляет условие field>=t. Время передаётся по МСК.
func (f Filter) Gte(field string, t time.Time) Filter {
	return f.add(field + ">=" + t.In(moscow()).Format(timeLayout))
}

// String возвращает значение параметра filter.
func (f Filter) String() string { return strings.Join(f.conds, ";") }

func (f Filter) add(cond string) Filter {
	return Filter{conds: append(slices.Clip(f.conds), cond)}
}
