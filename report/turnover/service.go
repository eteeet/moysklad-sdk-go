package turnover

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Значения groupBy отчёта по товарам; пустое — значение сервера по умолчанию (product).
const (
	GroupByProduct = "product" // только товары
	GroupByVariant = "variant" // товары и модификации
)

// AllOptions — параметры report/turnover/all. Период: без обеих границ — последний месяц, только
// MomentTo — с начала года, только MomentFrom — по текущий день (решает сервер). Filter —
// результат query.Filter: agent, agentTag, contract, organization, product, project, retailStore,
// store (несколько значений), supplier, type, variant, withoutturnover, archived и доп. поля;
// большинство полей принимает одно непустое значение — проверяет сервер. Limit=0 — 1000.
// Extra — прочие параметры как есть, без повторов typed-полей.
type AllOptions struct {
	Limit      int
	Offset     int
	MomentFrom model.DateTime
	MomentTo   model.DateTime
	GroupBy    string
	Filter     string
	Extra      url.Values
}

// DetailOptions — параметры отчётов по одному товару или модификации (ByStore). Assortment
// обязателен: meta.type product или variant выбирает фильтр product=… или variant=…. Filter —
// остальные условия через query.Filter (agent, agentTag, contract, organization, project,
// retailStore, store, type). Период — как у AllOptions; limit, offset и groupBy источник для этих
// отчётов не описывает.
type DetailOptions struct {
	Assortment model.Reference
	MomentFrom model.DateTime
	MomentTo   model.DateTime
	Filter     string
}

func (o DetailOptions) values() (url.Values, error) {
	key := string(o.Assortment.Meta.Type)
	if o.Assortment.Meta.Href == "" || key != "product" && key != "variant" {
		return nil, fmt.Errorf("moysklad: assortment must be a product or variant reference")
	}
	filter, err := query.Filter(query.Eq(key, o.Assortment.Meta.Href))
	if err != nil {
		return nil, err
	}
	if o.Filter != "" {
		// Обязательное условие задаёт Assortment; второй product или variant сервер не примет.
		for _, condition := range conditions(o.Filter) {
			if name := strings.TrimSpace(strings.SplitN(condition, "=", 2)[0]); name == "product" || name == "variant" {
				return nil, fmt.Errorf("moysklad: filter must not repeat product or variant; use Assortment")
			}
		}
		filter += ";" + o.Filter
	}
	return query.ReportOptions{MomentFrom: o.MomentFrom, MomentTo: o.MomentTo, Filter: filter}.Values()
}

// conditions делит filter по ; без экранирования \.
func conditions(filter string) []string {
	var out []string
	start := 0
	for i := 0; i < len(filter); i++ {
		switch filter[i] {
		case '\\':
			i++
		case ';':
			out = append(out, filter[start:i])
			start = i + 1
		}
	}
	return append(out, filter[start:])
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: turnover service requires a client")
	}
	return &Service{client: client}, nil
}

// All читает одну страницу отчёта «Обороты по товарам» без автообхода. Зону расчёта периода
// задаёт transport.WithAcceptTimezone; ответ сообщает её в X-Lognex-Content-Timezone (ResponseInfo).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-turnover#3-poluchit-oboroty-po-tovaram
func (s *Service) All(ctx context.Context, options AllOptions, opts ...transport.RequestOption) (*model.Collection[Row], error) {
	transport.ResetResponseInfo(opts)
	const operation = "get turnover all"
	if options.Extra.Has(query.ParamLimit) || options.Extra.Has(query.ParamOffset) {
		return nil, fmt.Errorf("moysklad: %s: use Limit and Offset instead of Extra", operation)
	}
	params, err := query.ReportOptions{MomentFrom: options.MomentFrom, MomentTo: options.MomentTo,
		GroupBy: options.GroupBy, Filter: options.Filter, Extra: options.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	page, err := query.PageOptions{Limit: options.Limit, Offset: options.Offset}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	for key, values := range page {
		params[key] = values
	}
	return list.Page[Row](ctx, s.client, operation, "report/turnover/all", params, opts...)
}

// ByStore читает отчёт по товару или модификации с детализацией по складам: одна строка на
// товар и каждую его модификацию; показатели складов SDK не суммирует и не выводит из All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-turnover#3-poluchit-oboroty-po-tovaru-s-detalizaciej-po-skladu
func (s *Service) ByStore(ctx context.Context, options DetailOptions, opts ...transport.RequestOption) (*model.Collection[StoreRow], error) {
	transport.ResetResponseInfo(opts)
	const operation = "get turnover bystore"
	params, err := options.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[StoreRow](ctx, s.client, operation, "report/turnover/bystore", params, opts...)
}
