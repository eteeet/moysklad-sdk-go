package profit

import (
	"context"
	"fmt"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
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

// ByVariant читает одну страницу отчёта по модификациям: модификация выводится сама, не
// родительский товар; товары без модификаций, услуги и комплекты тоже входят.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-modifikaciyam
func (s *Service) ByVariant(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[AssortmentRow], error) {
	transport.ResetResponseInfo(opts)
	return page[AssortmentRow](ctx, s.client, "byvariant", options, opts)
}

// ByEmployee читает одну страницу отчёта по сотрудникам. Путь — report/profit/byemployee, как в
// запросе раздела и upstream-фикстуре; meta.href примера ответа пишет byEmployee — противоречие
// источника, SDK отправляет нижний регистр, как остальные разрезы.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-sotrudnikam
func (s *Service) ByEmployee(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[EmployeeRow], error) {
	transport.ResetResponseInfo(opts)
	return page[EmployeeRow](ctx, s.client, "byemployee", options, opts)
}

// ByCounterparty читает одну страницу отчёта по покупателям.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-pokupatelyam
func (s *Service) ByCounterparty(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[CounterpartyRow], error) {
	transport.ResetResponseInfo(opts)
	return page[CounterpartyRow](ctx, s.client, "bycounterparty", options, opts)
}

// BySalesChannel читает одну страницу отчёта по каналам продаж.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#3-poluchit-pribylnost-po-kanalam-prodazh
func (s *Service) BySalesChannel(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[SalesChannelRow], error) {
	transport.ResetResponseInfo(opts)
	return page[SalesChannelRow](ctx, s.client, "bysaleschannel", options, opts)
}

func page[T any](ctx context.Context, client *transport.Client, slice string, options Options,
	opts []transport.RequestOption) (*model.Collection[T], error) {
	operation := "get profit " + slice
	if options.Extra.Has(query.ParamLimit) || options.Extra.Has(query.ParamOffset) {
		return nil, fmt.Errorf("moysklad: %s: use Limit and Offset instead of Extra", operation)
	}
	params, err := query.ReportOptions{MomentFrom: options.MomentFrom, MomentTo: options.MomentTo,
		Filter: options.Filter, Extra: options.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
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
