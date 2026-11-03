package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: dashboard service requires a client")
	}
	return &Service{client: client}, nil
}

// Day читает показатели за день.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/dashboard#3-poluchit-pokazateli-za-den
func (s *Service) Day(ctx context.Context, opts ...transport.RequestOption) (*Dashboard, error) {
	transport.ResetResponseInfo(opts)
	return s.get(ctx, "day", opts)
}

// Week читает показатели за неделю.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/dashboard#3-poluchit-pokazateli-za-nedelyu
func (s *Service) Week(ctx context.Context, opts ...transport.RequestOption) (*Dashboard, error) {
	transport.ResetResponseInfo(opts)
	return s.get(ctx, "week", opts)
}

// Month читает показатели за месяц.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/dashboard#3-poluchit-pokazateli-za-mesyac
func (s *Service) Month(ctx context.Context, opts ...transport.RequestOption) (*Dashboard, error) {
	transport.ResetResponseInfo(opts)
	return s.get(ctx, "month", opts)
}

// DayAsync создаёт асинхронную задачу показателей за день: GET report/dashboard/day?async=true.
// Результат — через async.Service.Result; создание задачи не повторяется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) DayAsync(ctx context.Context, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return s.submit(ctx, "day", opts)
}

// WeekAsync создаёт асинхронную задачу показателей за неделю: GET report/dashboard/week?async=true.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) WeekAsync(ctx context.Context, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return s.submit(ctx, "week", opts)
}

// MonthAsync создаёт асинхронную задачу показателей за месяц: GET report/dashboard/month?async=true.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) MonthAsync(ctx context.Context, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	return s.submit(ctx, "month", opts)
}

func (s *Service) submit(ctx context.Context, period string, opts []transport.RequestOption) (*async.Submitted, error) {
	return write.Submit(ctx, s.client, "get dashboard "+period+" async", "report/dashboard/"+period, url.Values{}, opts)
}

func (s *Service) get(ctx context.Context, period string, opts []transport.RequestOption) (*Dashboard, error) {
	operation := "get dashboard " + period
	// Указатели отличают отсутствующий или null раздел от нулевых показателей.
	var result struct {
		Sales  *Sales  `json:"sales"`
		Orders *Orders `json:"orders"`
		Money  *Money  `json:"money"`
	}
	resp, err := s.client.Send(ctx, http.MethodGet, "report/dashboard/"+period, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.Sales == nil || result.Orders == nil || result.Money == nil {
		return nil, fmt.Errorf("moysklad: %s: response must contain sales, orders and money", operation)
	}
	return &Dashboard{Sales: *result.Sales, Orders: *result.Orders, Money: *result.Money}, nil
}
