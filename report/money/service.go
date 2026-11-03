package money

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Значения interval из источника; набор открыт — другое значение уходит как есть.
const (
	IntervalHour  = "hour"
	IntervalDay   = "day"
	IntervalMonth = "month"
)

// PlotSeriesOptions — параметры report/money/plotseries. MomentFrom, MomentTo и Interval
// обязательны: пустые отклоняются до запроса. Filter — результат query.Filter по organization
// и project (ссылки, оператор =, можно несколько значений, пустые нельзя — проверяет сервер).
// Extra — прочие параметры как есть, без повторов typed-полей.
type PlotSeriesOptions struct {
	MomentFrom model.DateTime
	MomentTo   model.DateTime
	Interval   string
	Filter     string
	Extra      url.Values
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: money report service requires a client")
	}
	return &Service{client: client}, nil
}

// PlotSeries читает движение денежных средств за период с разбиением по Interval. Зону расчёта
// задаёт transport.WithAcceptTimezone; ответ сообщает её в X-Lognex-Content-Timezone (ResponseInfo).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-money#3-dvizhenie-denezhnyh-sredstv
func (s *Service) PlotSeries(ctx context.Context, options PlotSeriesOptions, opts ...transport.RequestOption) (*PlotSeries, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get money plotseries"
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	var result PlotSeries
	resp, err := s.client.Send(ctx, http.MethodGet, "report/money/plotseries", params, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Series == nil {
		return nil, fmt.Errorf("moysklad: %s: response series must be a non-null array", operation)
	}
	return &result, nil
}

// PlotSeriesAsync создаёт асинхронную задачу «Движение денежных средств»: GET
// report/money/plotseries с async=true и теми же обязательными параметрами, что у PlotSeries;
// limit и offset (в Extra) источник у async запрещает — ошибка до запроса. Результат — через
// async.Service.Result; создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) PlotSeriesAsync(ctx context.Context, options PlotSeriesOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get money plotseries async"
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/money/plotseries", params, opts)
}

func (o PlotSeriesOptions) values(operation string) (url.Values, error) {
	switch {
	case o.MomentFrom.IsZero():
		return nil, fmt.Errorf("moysklad: %s: momentFrom is required", operation)
	case o.MomentTo.IsZero():
		return nil, fmt.Errorf("moysklad: %s: momentTo is required", operation)
	case o.Interval == "":
		return nil, fmt.Errorf("moysklad: %s: interval is required", operation)
	}
	params, err := query.ReportOptions{MomentFrom: o.MomentFrom, MomentTo: o.MomentTo,
		Interval: o.Interval, Filter: o.Filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return params, nil
}

// ByAccount читает текущие остатки по кассам и счетам. Параметров источник не описывает; ответ —
// {context, meta, rows} без size/limit (в CollectionMeta они нулевые) и без страниц.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-money#3-ostatki-denezhnyh-sredstv
func (s *Service) ByAccount(ctx context.Context, opts ...transport.RequestOption) (*model.Collection[AccountRow], error) {
	transport.ResetResponseInfo(opts)
	return list.Page[AccountRow](ctx, s.client, "get money byaccount", "report/money/byaccount", nil, opts...)
}

// ByAccountAsync создаёт асинхронную задачу «Остатки денежных средств»: GET report/money/byaccount
// с одним параметром async=true. Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByAccountAsync(ctx context.Context, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return write.Submit(ctx, s.client, "get money byaccount async", "report/money/byaccount", url.Values{}, opts)
}
