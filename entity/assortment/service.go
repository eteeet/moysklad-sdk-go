package assortment

import (
	"context"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/pages"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// BetaFeature — тег заголовка X-Lognex-Remap-Beta-Feature, включающий новый endpoint без остатков.
const BetaFeature = "assortmentWithoutStock"

// Значения groupBy. Источник: product — только товары, variant — «товары и модификации
// (аналогично отсутствию параметра)», consignment — все сущности и партии. Пример без groupBy
// при этом возвращает и услуги с комплектами — противоречие; SDK передаёт значение как есть.
// Значения открыты: неизвестное передаётся как есть.
type GroupBy string

const (
	GroupByProduct     GroupBy = "product"
	GroupByVariant     GroupBy = "variant"
	GroupByConsignment GroupBy = "consignment"
)

// Options — общие ListOptions и groupBy. Фильтр type=consignment сервер выполняет только
// вместе с GroupByConsignment. fields у ассортимента не документирован и не принимается.
type Options struct {
	query.ListOptions
	GroupBy GroupBy
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: assortment service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу нового endpoint: всегда передаёт X-Lognex-Remap-Beta-Feature:
// assortmentWithoutStock последним, поверх опций вызова (исключение из порядка ADR 0001:
// чужой beta-тег вызова заменяется, иначе List молча ушёл бы в legacy или чужой режим). Остатков stock, reserve, inTransit,
// quantity в строках нет, и SDK не запрашивает их скрыто: они — в отчёте report/stock/all/current
// по id строки. Все страницы — All: прямой pages.All по странице List без опции
// transport.WithBetaFeature(BetaFeature) уводит GET по nextHref в legacy endpoint.
//
// Фильтры по источнику: archived (по умолчанию только неархивные; true и false вместе — все),
// barcode (=, ~, ~=, =~; без !=), code, description, externalCode, name, pathName, id, group,
// productFolder, withSubFolders (по умолчанию true, только с productFolder), supplier,
// type (product, service, bundle, variant, consignment; consignment — только с GroupByConsignment),
// updated и доп. поля. Недоступны: alcoholic.type, article, isSerialTrackable, minimumBalance,
// owner, quantityMode, salePrice, shared, stockMode, stockMoment, stockStore, updatedBy, weighed.
// Для одного поля совместимы только LIKE-комбинации ~= и =~. Сортировка недоступна по inTransit,
// minimumBalance, quantity, reserve, salePrice, stock. Поиск — отдельный параметр Search.
// Фильтры проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment
func (s *Service) List(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[Row], error) {
	transport.ResetResponseInfo(opts)
	return s.list(ctx, "list assortment", options, append(slices.Clip(opts), transport.WithBetaFeature(BetaFeature)))
}

// All — List и обход следующих страниц через pages.All: beta-тег BetaFeature добавляется к каждому
// GET по nextHref, поверх опций вызова, как в List. Запросы начинаются при обходе; ошибка первой
// страницы — единственный элемент. opts применяются к каждому запросу; ResponseInfo описывает
// последнюю загруженную страницу. Страницы ListLegacy обходятся pages.All без beta-опции.
func (s *Service) All(ctx context.Context, options Options, opts ...transport.RequestOption) iter.Seq2[*model.Collection[Row], error] {
	return func(yield func(*model.Collection[Row], error) bool) {
		first, err := s.List(ctx, options, opts...)
		if err != nil {
			yield(nil, err)
			return
		}
		for page, err := range pages.All(ctx, s.client, first, append(slices.Clip(opts), transport.WithBetaFeature(BetaFeature))...) {
			if !yield(page, err) {
				return
			}
		}
	}
}

// ListLegacy читает одну страницу старого endpoint без beta-заголовка; если assortmentWithoutStock
// задан в Config.Header или опциях вызова, возвращает ошибку до запроса. Строки товаров и
// модификаций содержат остатки, фильтры stockMode, quantityMode (по умолчанию all), stockStore,
// stockMoment, article, owner, shared, weighed и др. работают по разделу Deprecated. Поиск
// источник описывает только как фильтр search=… (одно значение, оператор =); параметр Search
// передаётся как есть, его поддержку определяет сервер.
//
// Deprecated: источник планирует вывести endpoint из эксплуатации в 2027 году; используйте List
// и отчёт об остатках.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment-legacy
func (s *Service) ListLegacy(ctx context.Context, options Options, opts ...transport.RequestOption) (*model.Collection[Row], error) {
	transport.ResetResponseInfo(opts)
	const operation = "list legacy assortment"
	// Тег из Config.Header или опций вызова молча переключил бы сервер на новый endpoint без остатков.
	for _, value := range s.client.HeaderValues(transport.HeaderBetaFeature, opts...) {
		for tag := range strings.SplitSeq(value, ",") {
			if strings.TrimSpace(tag) == BetaFeature {
				return nil, fmt.Errorf("moysklad: %s: %s %s selects the new endpoint; use List", operation, transport.HeaderBetaFeature, BetaFeature)
			}
		}
	}
	return s.list(ctx, operation, options, opts)
}

// DeleteMany удаляет позиции ассортимента разных типов одним POST entity/assortment/delete,
// не в корзину; до batch.MaxItems элементов. Источник: «Массовое удаление позиций в Ассортименте».
// metas — meta объектов (товар, услуга, комплект, модификация) с href и type, например
// row.Product.Meta. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#3-massovoe-udalenie-pozicij-v-assortimente
func (s *Service) DeleteMany(ctx context.Context, metas []model.Meta, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMetas(ctx, s.client, "assortment", metas, opts)
}

func (s *Service) list(ctx context.Context, operation string, options Options, opts []transport.RequestOption) (*model.Collection[Row], error) {
	if options.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields is not supported", operation)
	}
	if options.Extra.Has("groupBy") {
		return nil, fmt.Errorf("moysklad: %s: use Options.GroupBy instead of Extra groupBy", operation)
	}
	params, err := options.ListOptions.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if options.GroupBy != "" {
		params.Set("groupBy", string(options.GroupBy))
	}
	return list.Page[Row](ctx, s.client, operation, "entity/assortment", params, opts...)
}
