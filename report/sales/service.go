package sales

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Значения Interval из источника; набор не проверяется SDK — неизвестное значение отклонит сервер.
const (
	IntervalHour  = "hour"
	IntervalDay   = "day"
	IntervalMonth = "month"
)

// Options — параметры графиков. MomentFrom, MomentTo и Interval источник называет обязательными:
// пустое значение отклоняется до HTTP. Filter — результат query.Filter: organization, store,
// project у обоих графиков и retailStore только у Sales (тогда учитываются лишь розничные
// продажи); у каждого поля может быть несколько значений, пустые сервер не принимает. Extra —
// прочие параметры как есть, без повторов typed-полей.
type Options struct {
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
		return nil, fmt.Errorf("moysklad: sales service requires a client")
	}
	return &Service{client: client}, nil
}

// Sales читает показатели продаж: отгрузки, розничные продажи и полученные отчёты комиссионера.
// Зону расчёта задаёт transport.WithAcceptTimezone; ответ сообщает её в X-Lognex-Content-Timezone
// (ResponseInfo). Пример фильтра в источнике разделяет условия «&» — SDK передаёт filter, как его
// собрал query.Filter, через «;».
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-sales-orders#3-pokazateli-prodazh
func (s *Service) Sales(ctx context.Context, options Options, opts ...transport.RequestOption) (*Plot, error) {
	transport.ResetResponseInfo(opts)
	return s.plot(ctx, "sales", options, opts)
}

// Orders читает показатели заказов: учитываются только заказы покупателей.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-sales-orders#3-pokazateli-zakazov
func (s *Service) Orders(ctx context.Context, options Options, opts ...transport.RequestOption) (*Plot, error) {
	transport.ResetResponseInfo(opts)
	return s.plot(ctx, "orders", options, opts)
}

// SalesAsync создаёт асинхронную задачу показателей продаж: GET report/sales/plotseries с
// async=true и теми же обязательными параметрами, что у Sales; limit и offset (в Extra) источник у
// async запрещает — ошибка до запроса. Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) SalesAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return s.submit(ctx, "sales", options, opts)
}

// OrdersAsync создаёт асинхронную задачу показателей заказов: GET report/orders/plotseries с
// async=true; параметры и ограничения — как у SalesAsync.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) OrdersAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return s.submit(ctx, "orders", options, opts)
}

func (s *Service) submit(ctx context.Context, report string, options Options, opts []transport.RequestOption) (*async.Submitted, error) {
	operation := "get " + report + " plotseries async"
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/"+report+"/plotseries", params, opts)
}

func (o Options) values(operation string) (url.Values, error) {
	if o.MomentFrom.IsZero() || o.MomentTo.IsZero() || o.Interval == "" {
		return nil, fmt.Errorf("moysklad: %s: MomentFrom, MomentTo and Interval are required", operation)
	}
	params, err := query.ReportOptions{MomentFrom: o.MomentFrom, MomentTo: o.MomentTo,
		Interval: o.Interval, Filter: o.Filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return params, nil
}

func (s *Service) plot(ctx context.Context, report string, options Options, opts []transport.RequestOption) (*Plot, error) {
	operation := "get " + report + " plotseries"
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	var result Plot
	resp, err := s.client.Send(ctx, http.MethodGet, "report/"+report+"/plotseries", params, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	if result.Series == nil {
		return nil, fmt.Errorf("moysklad: %s: response series must be a non-null array", operation)
	}
	return &result, nil
}
