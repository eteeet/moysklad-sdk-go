package turnover

import (
	"context"
	"fmt"
	"net/url"

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
