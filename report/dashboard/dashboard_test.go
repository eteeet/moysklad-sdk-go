package dashboard_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/report/dashboard"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *dashboard.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := dashboard.NewService(c)
	return s
}

// Пример источника: продажи упали на 80000 копеек, баланс 69700 копеек.
var example = dashboard.Dashboard{
	Sales: dashboard.Sales{MovementAmount: -80000},
	Money: dashboard.Money{Balance: 69700},
}

func TestDashboardJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, path := range []string{"upstream/report_dashboard.json", "custom/dashboard/day.json", "custom/dashboard/week.json", "custom/dashboard/month.json"} {
				var got dashboard.Dashboard
				if err := codec.Unmarshal(read(t, path), &got); err != nil || got != example {
					t.Fatalf("%s: got = %+v, err = %v", path, got, err)
				}
			}
			for _, money := range []struct {
				json string
				want float64
			}{{"0", 0}, {"0.0", 0}, {"150.5", 150.5}, {"123456789012.25", 123456789012.25}} {
				body := fmt.Sprintf(`{"sales":{"count":3,"amount":%[1]s,"movementAmount":-%[1]s},"orders":{"count":2,"amount":%[1]s,"movementAmount":%[1]s},`+
					`"money":{"income":%[1]s,"outcome":%[1]s,"balance":%[1]s,"todayMovement":12.5,"movement":-7.25},"future":1}`, money.json)
				want := dashboard.Dashboard{
					Sales:  dashboard.Sales{Count: 3, Amount: money.want, MovementAmount: -money.want},
					Orders: dashboard.Orders{Count: 2, Amount: money.want, MovementAmount: money.want},
					Money:  dashboard.Money{Income: money.want, Outcome: money.want, Balance: money.want, TodayMovement: 12.5, Movement: -7.25},
				}
				var got dashboard.Dashboard
				if err := codec.Unmarshal([]byte(body), &got); err != nil || got != want {
					t.Errorf("%s: got = %+v, err = %v", money.json, got, err)
				}
			}
		})
	}
}

func TestHTTP(t *testing.T) {
	methods := map[string]func(*dashboard.Service, context.Context, ...transport.RequestOption) (*dashboard.Dashboard, error){
		"day":   (*dashboard.Service).Day,
		"week":  (*dashboard.Service).Week,
		"month": (*dashboard.Service).Month,
	}
	for period, method := range methods {
		t.Run(period, func(t *testing.T) {
			invoke := func(ctx context.Context, c *transport.Client) (*dashboard.Dashboard, error) {
				return method(service(c), ctx)
			}
			cases := documenttest.ResponseCases(string(read(t, "custom/dashboard/"+period+".json")))
			// Ответ без meta: каждый обязательный раздел проверяется отдельно.
			for _, missing := range []string{"sales", "orders", "money"} {
				body := strings.Replace(`{"sales":{},"orders":{},"money":{}}`, `"`+missing+`":{}`, `"`+missing+`":null`, 1)
				cases = append(cases, documenttest.ResponseCase{Name: "null " + missing, Status: 200, Body: body})
			}
			cases = append(cases, documenttest.ResponseCase{Name: "zero sections", Status: 200, Body: `{"sales":{},"orders":{},"money":{}}`, Success: true})
			for _, rc := range cases {
				t.Run(rc.Name, func(t *testing.T) {
					documenttest.HTTP(t, rc, "GET", "report/dashboard/"+period, nil, invoke, func(got *dashboard.Dashboard) {
						want := example
						if rc.Name == "zero sections" {
							want = dashboard.Dashboard{}
						}
						if *got != want {
							t.Fatalf("got = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, invoke)
			var calls int
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++ }))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			// nil context отклоняется до HTTP, ResponseInfo прошлого вызова сброшен.
			info := transport.ResponseInfo{StatusCode: 503}
			if _, err := method(service(c), nil, transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || calls != 0 ||
				strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get dashboard "+period) {
				t.Errorf("nil context: err = %v, info = %+v, calls = %d", err, info, calls)
			}
		})
	}
	if _, err := dashboard.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
