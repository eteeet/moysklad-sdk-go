package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Имена query-параметров JSON API без expand. Typed-поля есть только у ListOptions и ReportOptions;
// остальные параметры поддерживают отдельные endpoints и передаются через Extra.
// https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md
const (
	ParamLimit           = "limit"
	ParamOffset          = "offset"
	ParamFilter          = "filter"
	ParamOrder           = "order"
	ParamSearch          = "search"
	ParamFields          = "fields"
	ParamNamedFilter     = "namedfilter"
	ParamOwner           = "owner"
	ParamGroupBy         = "groupBy"
	ParamIncludeRelated  = "includeRelated"
	ParamInclude         = "include"
	ParamChangedSince    = "changedSince"
	ParamStockType       = "stockType"
	ParamWithRecalculate = "withRecalculate"
	ParamMomentFrom      = "momentFrom"
	ParamMomentTo        = "momentTo"
	ParamInterval        = "interval"
	ParamAction          = "action"
	ParamMiniature       = "miniature"
)

// MaxLimit — максимум limit большинства списков сущностей и документов.
const MaxLimit = 1000

// ListOptions — параметры одной страницы списка. Нулевое поле не передаётся:
// сервер применяет свой default (у списков сущностей limit=1000, offset=0).
// Поддержку search, namedfilter, fields и полей filter/order определяет endpoint.
type ListOptions struct {
	Limit       int
	Offset      int
	Filter      string // результат Filter
	Order       string // результат Order
	Search      string
	NamedFilter string // абсолютный href сохранённого фильтра этого типа сущности
	Fields      string // одно значение, например "stock"; несколько через запятую API отклоняет

	// Extra — прочие параметры как есть. Не может повторять typed-поля, expand и async.
	Extra url.Values
}

// Values проверяет опции и возвращает query; maxLimit задаёт endpoint (у audit — 100).
func (o ListOptions) Values(maxLimit int) (url.Values, error) {
	if o.Limit < 0 || o.Limit > maxLimit {
		return nil, fmt.Errorf("moysklad: limit must be between 1 and %d, or 0 for server default", maxLimit)
	}
	if o.Offset < 0 {
		return nil, fmt.Errorf("moysklad: offset must not be negative")
	}
	if strings.Contains(o.Fields, ",") {
		return nil, fmt.Errorf("moysklad: fields accepts a single value")
	}
	if o.NamedFilter != "" {
		if u, err := url.Parse(o.NamedFilter); err != nil || !u.IsAbs() {
			return nil, fmt.Errorf("moysklad: namedfilter must be an absolute href")
		}
	}
	return build(o.Extra, map[string]string{
		ParamLimit: nonZero(o.Limit), ParamOffset: nonZero(o.Offset), ParamFilter: o.Filter, ParamOrder: o.Order,
		ParamSearch: o.Search, ParamNamedFilter: o.NamedFilter, ParamFields: o.Fields,
	})
}

// ReportOptions — общие параметры отчётов. Обязательность периода, допустимые groupBy,
// interval и прочие параметры (limit, search …) определяет конкретный отчёт; их передают в Extra.
type ReportOptions struct {
	// Момент в московском времени; миллисекунды передаются, только если они не нулевые.
	MomentFrom model.DateTime
	MomentTo   model.DateTime
	GroupBy    string
	Interval   string // hour, day или month у графиков
	Filter     string // своя грамматика отчёта, см. RawCondition

	// Extra — прочие параметры как есть. Не может повторять typed-поля, expand и async.
	Extra url.Values
}

// Values возвращает query отчёта.
func (o ReportOptions) Values() (url.Values, error) {
	return build(o.Extra, map[string]string{
		ParamMomentFrom: moment(o.MomentFrom), ParamMomentTo: moment(o.MomentTo),
		ParamGroupBy: o.GroupBy, ParamInterval: o.Interval, ParamFilter: o.Filter,
	})
}

// build копирует extra и добавляет непустые typed-значения; typed-имена в extra запрещены.
func build(extra url.Values, typed map[string]string) (url.Values, error) {
	values := url.Values{}
	for name, list := range extra {
		if _, ok := typed[name]; ok {
			return nil, fmt.Errorf("moysklad: query parameter %s is set twice", name)
		}
		if name == "expand" || name == "async" {
			return nil, fmt.Errorf("moysklad: query parameter %s is not supported", name)
		}
		values[name] = append([]string(nil), list...)
	}
	for name, value := range typed {
		if value != "" {
			values.Set(name, value)
		}
	}
	return values, nil
}

func nonZero(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}

func moment(d model.DateTime) string {
	if d.IsZero() {
		return ""
	}
	return d.QueryValue()
}
