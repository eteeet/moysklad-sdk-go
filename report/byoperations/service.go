package byoperations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Options — параметры всех трёх отчётов. Assortment обязателен: ссылка на товар, модификацию
// или партию (meta.type product, variant или consignment) уходит фильтром
// filter=assortment=<href>; другой фильтрации источник не описывает. Limit=0 — 1000.
// Extra — прочие параметры как есть, без filter, limit и offset; имена общих параметров
// отчётов (momentFrom, momentTo, groupBy, interval), expand и async query.ReportOptions тоже
// отклоняет.
type Options struct {
	Assortment model.Reference
	Limit      int
	Offset     int
	Extra      url.Values
}

func (o Options) values() (url.Values, error) {
	if key := o.Assortment.Meta.Type; o.Assortment.Meta.Href == "" || key != "product" && key != "variant" && key != "consignment" {
		return nil, fmt.Errorf("moysklad: assortment must be a product, variant or consignment reference")
	}
	if o.Extra.Has(query.ParamLimit) || o.Extra.Has(query.ParamOffset) {
		return nil, fmt.Errorf("moysklad: use Limit and Offset instead of Extra")
	}
	filter, err := query.Filter(query.Eq("assortment", o.Assortment.Meta.Href))
	if err != nil {
		return nil, err
	}
	params, err := query.ReportOptions{Filter: filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, err
	}
	page, err := query.PageOptions{Limit: o.Limit, Offset: o.Offset}.Values()
	if err != nil {
		return nil, err
	}
	for key, values := range page {
		params[key] = values
	}
	return params, nil
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: byoperations service requires a client")
	}
	return &Service{client: client}, nil
}

// Stock читает одну страницу отчёта по документам, формирующим остатки номенклатуры, без
// автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-by-operations#3-poluchit-otchet-po-dokumentam-otobrazhayushij-ostatki
func (s *Service) Stock(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[StockRow], error) {
	transport.ResetResponseInfo(opts)
	return page[StockRow](ctx, s.client, "get byoperations stock", "report/byoperations/stock", options, opts)
}

// Reserve читает одну страницу отчёта по документам, формирующим резервы, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-by-operations#3-poluchit-otchet-po-dokumentam-otobrazhayushij-rezervy
func (s *Service) Reserve(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[ReserveRow], error) {
	transport.ResetResponseInfo(opts)
	return page[ReserveRow](ctx, s.client, "get byoperations reserve", "report/byoperations/reserve", options, opts)
}

// InTransit читает одну страницу отчёта по документам, формирующим ожидания, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-by-operations#3-poluchit-otchet-po-dokumentam-otobrazhayushij-ozhidanie
func (s *Service) InTransit(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[InTransitRow], error) {
	transport.ResetResponseInfo(opts)
	return page[InTransitRow](ctx, s.client, "get byoperations intransit", "report/byoperations/intransit", options, opts)
}

func page[T any](ctx context.Context, client *transport.Client, operation, endpoint string, options Options, opts []transport.RequestOption) (*model.Collection[T], error) {
	params, err := options.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[T](ctx, client, operation, endpoint, params, opts...)
}
