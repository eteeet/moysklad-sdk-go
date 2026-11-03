package profit

import (
	"context"
	"fmt"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Options — параметры всех разрезов отчёта. Период: без обеих границ — последний месяц, только
// MomentTo — с начала года, только MomentFrom — по текущий день (решает сервер). Filter — результат
// query.Filter: product, productFolder (= и !=), withSubFolders, agentTag, entityType (demand,
// retaildemand), counterparty, organization, store, project, retailStore (кроме BySalesChannel),
// supplier, salesChannel. Несколько значений допускают product, productFolder и salesChannel;
// одновременно product и productFolder сервер не принимает — проверяет сервер. Limit=0 — 1000.
// Extra — прочие параметры как есть, без повторов typed-полей.
type Options struct {
	Limit      int
	Offset     int
	MomentFrom model.DateTime
	MomentTo   model.DateTime
	Filter     string
	Extra      url.Values
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: profit service requires a client")
	}
	return &Service{client: client}, nil
}

// ByProduct читает одну страницу отчёта по товарам, услугам и комплектам без автообхода;
// модификации учтены в своём товаре. Зону расчёта периода задаёт transport.WithAcceptTimezone.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-tovaram
func (s *Service) ByProduct(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[AssortmentRow], error) {
	transport.ResetResponseInfo(opts)
	return page[AssortmentRow](ctx, s.client, "byproduct", options, opts)
}

// ByProductAsync создаёт асинхронную задачу отчёта «Прибыльность по товарам»: GET report/profit/byproduct с
// async=true и параметрами options, кроме Limit и Offset (источник запрещает их у async —
// заданные отклоняются до запроса). Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByProductAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return submit(ctx, s.client, "byproduct", options, opts)
}

// ByVariant читает одну страницу отчёта по модификациям: модификация выводится сама, не
// родительский товар; товары без модификаций, услуги и комплекты тоже входят.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-modifikaciyam
func (s *Service) ByVariant(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[AssortmentRow], error) {
	transport.ResetResponseInfo(opts)
	return page[AssortmentRow](ctx, s.client, "byvariant", options, opts)
}

// ByVariantAsync создаёт асинхронную задачу отчёта «Прибыльность по модификациям»: GET report/profit/byvariant с
// async=true и параметрами options, кроме Limit и Offset (источник запрещает их у async —
// заданные отклоняются до запроса). Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByVariantAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return submit(ctx, s.client, "byvariant", options, opts)
}

// ByEmployee читает одну страницу отчёта по сотрудникам. Путь — report/profit/byemployee, как в
// запросе раздела и upstream-фикстуре; meta.href примера ответа пишет byEmployee — противоречие
// источника, SDK отправляет нижний регистр, как остальные разрезы.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-sotrudnikam
func (s *Service) ByEmployee(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[EmployeeRow], error) {
	transport.ResetResponseInfo(opts)
	return page[EmployeeRow](ctx, s.client, "byemployee", options, opts)
}

// ByEmployeeAsync создаёт асинхронную задачу отчёта «Прибыльность по сотрудникам»: GET report/profit/byemployee с
// async=true и параметрами options, кроме Limit и Offset (источник запрещает их у async —
// заданные отклоняются до запроса). Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByEmployeeAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return submit(ctx, s.client, "byemployee", options, opts)
}

// ByCounterparty читает одну страницу отчёта по покупателям.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-pokupatelyam
func (s *Service) ByCounterparty(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[CounterpartyRow], error) {
	transport.ResetResponseInfo(opts)
	return page[CounterpartyRow](ctx, s.client, "bycounterparty", options, opts)
}

// ByCounterpartyAsync создаёт асинхронную задачу отчёта «Прибыльность по покупателям»: GET report/profit/bycounterparty с
// async=true и параметрами options, кроме Limit и Offset (источник запрещает их у async —
// заданные отклоняются до запроса). Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByCounterpartyAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return submit(ctx, s.client, "bycounterparty", options, opts)
}

// BySalesChannel читает одну страницу отчёта по каналам продаж.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-kanalam-prodazh
func (s *Service) BySalesChannel(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[SalesChannelRow], error) {
	transport.ResetResponseInfo(opts)
	return page[SalesChannelRow](ctx, s.client, "bysaleschannel", options, opts)
}

// BySalesChannelAsync создаёт асинхронную задачу отчёта «Прибыльность по каналам продаж»: GET report/profit/bysaleschannel с
// async=true и параметрами options, кроме Limit и Offset (источник запрещает их у async —
// заданные отклоняются до запроса). Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) BySalesChannelAsync(ctx context.Context, options Options, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return submit(ctx, s.client, "bysaleschannel", options, opts)
}

func page[T any](ctx context.Context, client *transport.Client, slice string, options Options,
	opts []transport.RequestOption) (*model.Collection[T], error) {
	operation := "get profit " + slice
	if options.Extra.Has(query.ParamLimit) || options.Extra.Has(query.ParamOffset) {
		return nil, fmt.Errorf("moysklad: %s: use Limit and Offset instead of Extra", operation)
	}
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	paging, err := query.PageOptions{Limit: options.Limit, Offset: options.Offset}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	for key, values := range paging {
		params[key] = values
	}
	return list.Page[T](ctx, client, operation, "report/profit/"+slice, params, opts...)
}

func submit(ctx context.Context, client *transport.Client, slice string, options Options, opts []transport.RequestOption) (*async.Submitted, error) {
	operation := "get profit " + slice + " async"
	if options.Limit != 0 || options.Offset != 0 {
		return nil, fmt.Errorf("moysklad: %s: limit and offset are not allowed with async", operation)
	}
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, client, operation, "report/profit/"+slice, params, opts)
}

func (o Options) values(operation string) (url.Values, error) {
	params, err := query.ReportOptions{MomentFrom: o.MomentFrom, MomentTo: o.MomentTo, Filter: o.Filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return params, nil
}
