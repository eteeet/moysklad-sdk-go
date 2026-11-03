package counterparty

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/single"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// ListOptions — параметры report/counterparty. Limit=0 — 1000. Filter — результат query.Filter по
// атрибутам с фильтрацией: averageReceipt, balance, demandsCount, demandsSum, discountsSum,
// firstDemandDate, lastDemandDate, lastEventDate, lastEventText, profit, returnsCount, returnsSum,
// updated, id и counterparty.name/companyType/description/email/inn/phone; операторы проверяет
// сервер. Extra — прочие параметры как есть, без повторов typed-полей.
type ListOptions struct {
	Limit  int
	Offset int
	Filter string
	Extra  url.Values
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: counterparty report service requires a client")
	}
	return &Service{client: client}, nil
}

// List читает одну страницу показателей всех контрагентов без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-counterparty#3-poluchit-pokazateli-kontragentov
func (s *Service) List(ctx context.Context, options ListOptions, opts ...transport.RequestOption) (*model.Collection[Row], error) {
	transport.ResetResponseInfo(opts)
	const operation = "get counterparty report"
	if options.Extra.Has(query.ParamLimit) || options.Extra.Has(query.ParamOffset) {
		return nil, fmt.Errorf("moysklad: %s: use Limit and Offset instead of Extra", operation)
	}
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	page, err := query.PageOptions{Limit: options.Limit, Offset: options.Offset}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	for key, values := range page {
		params[key] = values
	}
	return list.Page[Row](ctx, s.client, operation, "report/counterparty", params, opts...)
}

// ListAsync создаёт асинхронную задачу показателей всех контрагентов: GET report/counterparty с
// async=true и Filter/Extra; Limit и Offset источник у async запрещает — заданные отклоняются до
// запроса. Результат — через async.Service.Result; создание задачи не повторяется. У Selected
// (POST) async нет: источник исключает выборочные показатели.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ListAsync(ctx context.Context, options ListOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "get counterparty report async"
	if options.Limit != 0 || options.Offset != 0 {
		return nil, fmt.Errorf("moysklad: %s: limit and offset are not allowed with async", operation)
	}
	params, err := options.values(operation)
	if err != nil {
		return nil, err
	}
	return write.Submit(ctx, s.client, operation, "report/counterparty", params, opts)
}

func (o ListOptions) values(operation string) (url.Values, error) {
	params, err := query.ReportOptions{Filter: o.Filter, Extra: o.Extra}.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return params, nil
}

type selectedRequest struct {
	Counterparties []selectedItem `json:"counterparties"`
}

type selectedItem struct {
	Counterparty model.Reference `json:"counterparty"`
}

// Selected читает показатели указанных контрагентов: POST report/counterparty с массивом
// counterparties — это чтение, ничего не изменяет. Нужна хотя бы одна ссылка с непустым
// meta.href; тип и существование контрагентов проверяет сервер. Параметров запроса и лимита на
// число ссылок источник не описывает. Transport не повторяет POST: одна отправка.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-counterparty#3-vyborochnye-pokazateli-kontragentov
func (s *Service) Selected(ctx context.Context, counterparties []model.Reference, opts ...transport.RequestOption) (*model.Collection[Row], error) {
	transport.ResetResponseInfo(opts)
	const operation = "select counterparty report"
	if len(counterparties) == 0 {
		return nil, fmt.Errorf("moysklad: %s: counterparties must not be empty", operation)
	}
	body := selectedRequest{Counterparties: make([]selectedItem, len(counterparties))}
	for i, ref := range counterparties {
		if ref.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: counterparties[%d]: meta.href must not be empty", operation, i)
		}
		body.Counterparties[i] = selectedItem{Counterparty: ref}
	}
	var result model.Collection[Row]
	resp, err := s.client.Send(ctx, http.MethodPost, "report/counterparty", nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no meta.href", operation)
	}
	if result.Rows == nil {
		return nil, fmt.Errorf("moysklad: %s: response rows must be a non-null array", operation)
	}
	return &result, nil
}

// Get читает показатели одного контрагента по его id. Параметры limit и offset, которые источник
// перечисляет у этого запроса, к одному объекту неприменимы — SDK их не передаёт.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-counterparty#3-pokazateli-kontragenta
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Row, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get counterparty report: id must be a UUID")
	}
	return single.Get[Row](ctx, s.client, "get counterparty report "+id, "report/counterparty/"+id, opts)
}

// GetAsync создаёт асинхронную задачу показателей одного контрагента: GET report/counterparty/{id}
// с одним параметром async=true. Результат — через async.Service.Result; без повторов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) GetAsync(ctx context.Context, id string, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get counterparty report async: id must be a UUID")
	}
	return write.Submit(ctx, s.client, "get counterparty report "+id+" async", "report/counterparty/"+id, url.Values{}, opts)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
