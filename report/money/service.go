package money

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
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
	switch {
	case options.MomentFrom.IsZero():
		return nil, fmt.Errorf("moysklad: %s: momentFrom is required", operation)
	case options.MomentTo.IsZero():
		return nil, fmt.Errorf("moysklad: %s: momentTo is required", operation)
	case options.Interval == "":
		return nil, fmt.Errorf("moysklad: %s: interval is required", operation)
	}
	params, err := query.ReportOptions{MomentFrom: options.MomentFrom, MomentTo: options.MomentTo,
		Interval: options.Interval, Filter: options.Filter, Extra: options.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
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

// ByAccount читает текущие остатки по кассам и счетам. Параметров источник не описывает; ответ —
// {context, meta, rows} без size/limit (в CollectionMeta они нулевые) и без страниц.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-money#3-ostatki-denezhnyh-sredstv
func (s *Service) ByAccount(ctx context.Context, opts ...transport.RequestOption) (*model.Collection[AccountRow], error) {
	transport.ResetResponseInfo(opts)
	return list.Page[AccountRow](ctx, s.client, "get money byaccount", "report/money/byaccount", nil, opts...)
}
