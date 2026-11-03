package stock

import (
	"context"
	"encoding/hex"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Значения groupBy расширенного отчёта и отчёта по складам; пустое — значение сервера (variant).
const (
	GroupByProduct     = "product"     // только товары
	GroupByVariant     = "variant"     // товары и модификации
	GroupByConsignment = "consignment" // товары, модификации и партии
)

// Значения условия filter stockMode (query.Eq("stockMode", …)); набор открыт, сервер может
// добавить новые. По умолчанию у All — all, у ByStore — nonEmpty.
const (
	StockModeAll          = "all"
	StockModePositiveOnly = "positiveOnly"
	StockModeNegativeOnly = "negativeOnly"
	StockModeEmpty        = "empty"
	StockModeNonEmpty     = "nonEmpty"
	StockModeUnderMinimum = "underMinimum"
)

// Значения stockType кратких отчётов; пустое — значение сервера (stock). Набор открыт.
const (
	StockTypeStock     = "stock"     // физический остаток без учёта резерва и ожидания
	StockTypeFreeStock = "freeStock" // остаток за вычетом резерва
	StockTypeQuantity  = "quantity"  // доступно: с учётом резерва и ожидания
	StockTypeReserve   = "reserve"
	StockTypeInTransit = "inTransit"
)

// IncludeZeroLines — значение include кратких отчётов: вывести и нулевые остатки.
const IncludeZeroLines = "zeroLines"

// AllOptions — параметры report/stock/all. Limit=0 — 1000. Filter — результат query.Filter:
// archived, consignment, inTransitOnly, moment, product, productFolder, withSubFolders,
// quantityMode, reserveOnly, search, soldByWeight, stockDaysFrom, stockDaysTo, stockMode, store,
// supplier, variant, доп. поля и характеристики; совместимость операторов проверяет сервер. Order —
// результат query.Order (avgStockDays, code, inTransit, minimumBalance, name, pathName, price,
// productCode, quantity, reserve, salePrice, stock, sumTotal). IncludeRelated требует в Filter
// условие product или variant — без него запрос не уходит. Extra — прочие параметры как есть,
// без повторов typed-полей.
type AllOptions struct {
	Limit          int
	Offset         int
	GroupBy        string
	IncludeRelated bool
	Filter         string
	Order          string
	Extra          url.Values
}

// ByStoreOptions — параметры report/stock/bystore. Limit=0 — 1000. Filter — результат
// query.Filter: consignment, moment, product, productFolder, search, soldByWeight, stockMode,
// store, supplier, variant, доп. поля и характеристики (только =, без пустых значений; только одно
// из product, variant, consignment, productFolder, search — проверяет сервер). Order — pathName,
// name, code, productCode, stockOnAllStores. includeRelated источник для этого отчёта не описывает.
type ByStoreOptions struct {
	Limit   int
	Offset  int
	GroupBy string
	Filter  string
	Order   string
	Extra   url.Values
}

// CurrentOptions — параметры кратких отчётов AllCurrent и ByStoreCurrent. StockType — один вид
// значения (StockType*); Include — IncludeZeroLines; ChangedSince — остатки, изменившиеся с этого
// момента (не дальше 24 часов в прошлое и не в будущем — проверяет сервер), передаётся в
// московском времени. WithRecalculate — только для ByStoreCurrent. AssortmentIDs и StoreIDs —
// UUID для filter assortmentId и storeId. Extra — прочие параметры как есть.
type CurrentOptions struct {
	StockType       string
	Include         string
	ChangedSince    model.DateTime
	WithRecalculate bool
	AssortmentIDs   []string
	StoreIDs        []string
	Extra           url.Values
}

// SlotOptions — параметры BySlotCurrent: UUID для filter assortmentId и storeId; источник требует
// хотя бы одно из них. Extra — прочие параметры как есть.
type SlotOptions struct {
	AssortmentIDs []string
	StoreIDs      []string
	Extra         url.Values
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: stock service requires a client")
	}
	return &Service{client: client}, nil
}

// All читает одну страницу расширенного отчёта об остатках без автообхода. Запрос весит 5 единиц
// лимита (transport.DefaultWeight).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-rasshirennyj-otchet-ob-ostatkah
func (s *Service) All(ctx context.Context, options AllOptions, opts ...transport.RequestOption) (*model.Collection[Row], error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock all"
	params, err := options.values(operation, true)
	if err != nil {
		return nil, err
	}
	return list.Page[Row](ctx, s.client, operation, "report/stock/all", params, opts...)
}

// AllAsync создаёт асинхронную задачу расширенного отчёта об остатках: GET report/stock/all с
// async=true и параметрами options, кроме Limit и Offset — источник запрещает их у async, заданные
// отклоняются до запроса. Ответ — созданная задача; отчёт целиком — через async.Service.Result.
// Создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) AllAsync(ctx context.Context, options AllOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock all async"
	params, err := options.values(operation, false)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/stock/all", params, opts)
}

func (o AllOptions) values(operation string, paged bool) (url.Values, error) {
	if o.IncludeRelated && !filtersAssortment(o.Filter) {
		return nil, fmt.Errorf("moysklad: %s: includeRelated requires a product or variant filter", operation)
	}
	return pageValues(operation, paged, o.Limit, o.Offset, o.GroupBy, o.Filter, o.Order, o.IncludeRelated, o.Extra)
}

// ByStore читает одну страницу отчёта «Остатки по складам» без автообхода. Запрос весит 5 единиц
// лимита (transport.DefaultWeight).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-ostatki-po-skladam
func (s *Service) ByStore(ctx context.Context, options ByStoreOptions, opts ...transport.RequestOption) (*model.Collection[StoreRow], error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock bystore"
	params, err := pageValues(operation, true, options.Limit, options.Offset, options.GroupBy, options.Filter, options.Order, false, options.Extra)
	if err != nil {
		return nil, err
	}
	return list.Page[StoreRow](ctx, s.client, operation, "report/stock/bystore", params, opts...)
}

// ByStoreAsync создаёт асинхронную задачу отчёта «Остатки по складам»: GET report/stock/bystore с
// async=true и параметрами options, кроме Limit и Offset (источник запрещает их у async — заданные
// отклоняются до запроса). Результат — через async.Service.Result; создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByStoreAsync(ctx context.Context, options ByStoreOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock bystore async"
	params, err := pageValues(operation, false, options.Limit, options.Offset, options.GroupBy, options.Filter, options.Order, false, options.Extra)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/stock/bystore", params, opts)
}

// AllCurrent читает краткий отчёт об остатках по товарам: весь ответ — один массив, без страниц и
// сортировки. Отчёт доступен администратору или сотруднику с правом на остатки (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-kratkij-otchet-ob-ostatkah
func (s *Service) AllCurrent(ctx context.Context, options CurrentOptions, opts ...transport.RequestOption) ([]CurrentRow, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock all current"
	params, err := options.allValues(operation)
	if err != nil {
		return nil, err
	}
	return current[CurrentRow](ctx, s.client, operation, "report/stock/all/current", params, opts)
}

// AllCurrentAsync создаёт асинхронную задачу краткого отчёта об остатках по товарам: GET
// report/stock/all/current с async=true и параметрами options. Источник относит краткие отчёты к
// «Отчету Остатки» из списка async; limit и offset (в Extra) отклоняются до запроса. Результат —
// через async.Service.Result; создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) AllCurrentAsync(ctx context.Context, options CurrentOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock all current async"
	params, err := options.allValues(operation)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/stock/all/current", params, opts)
}

func (o CurrentOptions) allValues(operation string) (url.Values, error) {
	if o.WithRecalculate {
		return nil, fmt.Errorf("moysklad: %s: withRecalculate is supported only by bystore/current", operation)
	}
	params, err := o.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return params, nil
}

// ByStoreCurrent читает краткий отчёт об остатках по товарам и складам одним массивом. Строка без
// склада (StoreID "") содержит резерв или ожидание заказа без склада.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-kratkij-otchet-ob-ostatkah
func (s *Service) ByStoreCurrent(ctx context.Context, options CurrentOptions, opts ...transport.RequestOption) ([]StoreCurrentRow, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock bystore current"
	params, err := options.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return current[StoreCurrentRow](ctx, s.client, operation, "report/stock/bystore/current", params, opts)
}

// ByStoreCurrentAsync создаёт асинхронную задачу краткого отчёта об остатках по товарам и складам:
// GET report/stock/bystore/current с async=true и параметрами options; limit и offset (в Extra)
// отклоняются до запроса. Результат — через async.Service.Result; создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByStoreCurrentAsync(ctx context.Context, options CurrentOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock bystore current async"
	params, err := options.values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return write.Submit(ctx, s.client, operation, "report/stock/bystore/current", params, opts)
}

// BySlotCurrent читает краткий отчёт об остатках по ячейкам одним массивом; остатки вне ячеек в
// него не входят.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-kratkij-otchet-ob-ostatkah-po-yachejkam
func (s *Service) BySlotCurrent(ctx context.Context, options SlotOptions, opts ...transport.RequestOption) ([]SlotRow, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock byslot current"
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	return current[SlotRow](ctx, s.client, operation, "report/stock/byslot/current", params, opts)
}

// BySlotCurrentAsync создаёт асинхронную задачу краткого отчёта об остатках по ячейкам: GET
// report/stock/byslot/current с async=true; как и у BySlotCurrent, нужен хотя бы один UUID, limit и
// offset (в Extra) отклоняются до запроса. Результат — через async.Service.Result.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) BySlotCurrentAsync(ctx context.Context, options SlotOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock byslot current async"
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/stock/byslot/current", params, opts)
}

func (o SlotOptions) values(operation string) (url.Values, error) {
	if len(o.AssortmentIDs) == 0 && len(o.StoreIDs) == 0 {
		return nil, fmt.Errorf("moysklad: %s: assortment or store IDs are required", operation)
	}
	filter, err := idFilter(o.AssortmentIDs, o.StoreIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	params, err := query.ReportOptions{Filter: filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return params, nil
}

// ByOperation читает остатки по позициям документа operationID (отгрузка, заказ покупателя,
// розничная продажа, счета, заказ поставщику, приёмка, возвраты): на момент документа или на
// текущий момент — по типу документа. Тип документа сервер определяет по UUID.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-ostatki-po-dokumentu
func (s *Service) ByOperation(ctx context.Context, operationID string, opts ...transport.RequestOption) (*model.Collection[OperationRow], error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock byoperation"
	if !validID(operationID) {
		return nil, fmt.Errorf("moysklad: %s: operation ID must be a UUID", operation)
	}
	return list.Page[OperationRow](ctx, s.client, operation, "report/stock/byoperation", url.Values{"operation.id": {operationID}}, opts...)
}

// ByOperationAsync создаёт асинхронную задачу отчёта «Остатки по документам»: GET
// report/stock/byoperation?operation.id=…&async=true. Результат — через async.Service.Result;
// создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ByOperationAsync(ctx context.Context, operationID string, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get stock byoperation async"
	if !validID(operationID) {
		return nil, fmt.Errorf("moysklad: %s: operation ID must be a UUID", operation)
	}
	return write.Submit(ctx, s.client, operation, "report/stock/byoperation", url.Values{"operation.id": {operationID}}, opts)
}

// pageValues собирает параметры All и ByStore; paged=false — для async, где источник запрещает
// limit и offset: заданные Limit/Offset — ошибка, в запрос они не попадают.
func pageValues(operation string, paged bool, limit, offset int, groupBy, filter, order string, includeRelated bool, extra url.Values) (url.Values, error) {
	if !paged && (limit != 0 || offset != 0) {
		return nil, fmt.Errorf("moysklad: %s: limit and offset are not allowed with async", operation)
	}
	for _, name := range []string{query.ParamLimit, query.ParamOffset, query.ParamOrder, query.ParamIncludeRelated} {
		if extra.Has(name) {
			return nil, fmt.Errorf("moysklad: %s: query parameter %s must not be set in Extra", operation, name)
		}
	}
	params, err := query.ReportOptions{GroupBy: groupBy, Filter: filter, Extra: extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if paged {
		page, err := query.PageOptions{Limit: limit, Offset: offset}.Values()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		maps.Copy(params, page)
	}
	if order != "" {
		params.Set(query.ParamOrder, order)
	}
	if includeRelated {
		params.Set(query.ParamIncludeRelated, "true")
	}
	return params, nil
}

// idFilter собирает filter кратких отчётов: assortmentId и storeId — UUID через запятую.
func idFilter(assortmentIDs, storeIDs []string) (string, error) {
	var conditions []query.Condition
	for _, f := range []struct {
		name string
		ids  []string
	}{{"assortmentId", assortmentIDs}, {"storeId", storeIDs}} {
		for _, id := range f.ids {
			if !validID(id) {
				return "", fmt.Errorf("moysklad: %s %q must be a UUID", f.name, id)
			}
		}
		if len(f.ids) > 0 {
			conditions = append(conditions, query.Eq(f.name, strings.Join(f.ids, ",")))
		}
	}
	if len(conditions) == 0 {
		return "", nil
	}
	return query.Filter(conditions...)
}

func (o CurrentOptions) values() (url.Values, error) {
	filter, err := idFilter(o.AssortmentIDs, o.StoreIDs)
	if err != nil {
		return nil, err
	}
	params, err := query.ReportOptions{Filter: filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, err
	}
	changedSince, recalculate := "", ""
	if !o.ChangedSince.IsZero() {
		changedSince = o.ChangedSince.QueryValue()
	}
	if o.WithRecalculate {
		recalculate = "true"
	}
	for name, value := range map[string]string{query.ParamStockType: o.StockType, query.ParamInclude: o.Include,
		query.ParamChangedSince: changedSince, query.ParamWithRecalculate: recalculate} {
		if params.Has(name) {
			return nil, fmt.Errorf("moysklad: query parameter %s is set twice", name)
		}
		if value != "" {
			params.Set(name, value)
		}
	}
	return params, nil
}

// current выполняет один GET краткого отчёта: ответ — JSON-массив, не Collection.
func current[T any](ctx context.Context, client *transport.Client, operation, endpoint string, params url.Values, opts []transport.RequestOption) ([]T, error) {
	var rows []T
	resp, err := client.Send(ctx, http.MethodGet, endpoint, params, nil, &rows, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if rows == nil {
		return nil, fmt.Errorf("moysklad: %s: response must be a JSON array", operation)
	}
	return rows, nil
}

// filtersAssortment сообщает, есть ли в filter условие product или variant (= или !=).
func filtersAssortment(filter string) bool {
	start := 0
	for i := 0; i <= len(filter); i++ {
		if i < len(filter) && filter[i] == '\\' {
			i++
			continue
		}
		if i == len(filter) || filter[i] == ';' {
			condition := filter[start:i]
			if end := strings.IndexAny(condition, "!="); end > 0 {
				if name := strings.TrimSpace(condition[:end]); name == "product" || name == "variant" {
					return true
				}
			}
			start = i + 1
		}
	}
	return false
}

// validID проверяет форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
