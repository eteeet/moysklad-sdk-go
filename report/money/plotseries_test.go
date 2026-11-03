package money_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/money"
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

func service(c *transport.Client) *money.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := money.NewService(c)
	return s
}

// points сравнивает серию: даты идут по дням с 2018-09-01, want — credit, debit, balance.
func points(t *testing.T, got []money.Point, want ...[3]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("series = %+v", got)
	}
	start := documenttest.DateTime("2018-09-01 00:00:00")
	for i, p := range got {
		if !time.Time(p.Date).Equal(time.Time(start).AddDate(0, 0, i)) || p.Credit != want[i][0] || p.Debit != want[i][1] || p.Balance != want[i][2] {
			t.Errorf("point %d = %+v, want %v", i, p, want[i])
		}
	}
}

func TestPlotSeriesJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, path := range []string{"upstream/report_money_plotseries.json", "custom/money/plotseries.json"} {
				var got money.PlotSeries
				if err := codec.Unmarshal(read(t, path), &got); err != nil {
					t.Fatal(err)
				}
				if got.Credit != 400 || got.Debit != 200 || got.Meta.Type != "moneyplotseries" || got.Context == nil || got.Context.Employee == nil ||
					got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/report/money/plotseries?momentFrom=2018-09-01&momentTo=2018-09-04&interval=day" {
					t.Fatalf("%s: report = %+v", path, got)
				}
				points(t, got.Series, [3]float64{100, 50, 50}, [3]float64{100, 50, 50}, [3]float64{100, 50, 50}, [3]float64{100, 50, 50})
			}
			var filtered money.PlotSeries
			if err := codec.Unmarshal(read(t, "custom/money/plotseries_filter.json"), &filtered); err != nil {
				t.Fatal(err)
			}
			if filtered.Credit != 60 || filtered.Debit != 20 {
				t.Fatalf("filtered = %+v", filtered)
			}
			points(t, filtered.Series, [3]float64{50, 0, 50}, [3]float64{10, 20, -10}, [3]float64{}, [3]float64{})
		})
	}
}

// Деньги — копейки как есть: 0, 0.0, дробь и большое значение.
func TestPlotSeriesMoneyJSON(t *testing.T) {
	data := []byte(`{"credit":0,"debit":0.0,"series":[{"date":"2018-09-01 00:00:00","credit":150.5,"debit":999999999999.99,"balance":-999999999849.49}],"future":1}`)
	for _, codec := range documenttest.Codecs {
		var got money.PlotSeries
		if err := codec.Unmarshal(data, &got); err != nil || got.Credit != 0 || got.Debit != 0 || got.Context != nil {
			t.Fatalf("%s: report = %+v, err = %v", codec.Name, got, err)
		}
		points(t, got.Series, [3]float64{150.5, 999999999999.99, -999999999849.49})
	}
}

func TestPlotSeriesHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("organization", "https://api.moysklad.ru/api/remap/1.2/entity/organization/00cd5a99-6897-11e7-7a6c-d2a9000c4fc0"),
		query.Eq("project", "https://api.moysklad.ru/api/remap/1.2/entity/project/02e64f51-6897-11e7-7a34-5acf000c8448"))
	if err != nil {
		t.Fatal(err)
	}
	period := money.PlotSeriesOptions{MomentFrom: documenttest.DateTime("2018-09-01 00:00:00"), MomentTo: documenttest.DateTime("2018-09-04 00:00:00"), Interval: money.IntervalDay}
	filtered := period
	filtered.Filter = filter
	hour := period
	hour.Interval, hour.Extra = "future", url.Values{"x": {"1"}}
	const moments = "momentFrom=2018-09-01+00%3A00%3A00&momentTo=2018-09-04+00%3A00%3A00"
	for _, tc := range []struct {
		name    string
		options money.PlotSeriesOptions
		query   string
		body    string
	}{
		{"period", period, "interval=day&" + moments, "custom/money/plotseries.json"},
		{"filter", filtered, "filter=organization%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Forganization%2F00cd5a99-6897-11e7-7a6c-d2a9000c4fc0%3Bproject%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproject%2F02e64f51-6897-11e7-7a34-5acf000c8448&interval=day&" + moments, "custom/money/plotseries_filter.json"},
		{"open interval", hour, "interval=future&" + moments + "&x=1", "custom/money/plotseries.json"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*money.PlotSeries, error) {
			return service(c).PlotSeries(ctx, tc.options)
		}
		for _, rc := range documenttest.ResponseCases(string(read(t, tc.body))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/money/plotseries?"+tc.query, nil, invoke, func(got *money.PlotSeries) {
					if len(got.Series) != 4 || got.Meta.Type != "moneyplotseries" || got.Credit == 0 {
						t.Fatalf("report = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Зона расчёта уходит заголовком X-Lognex-Accept-Timezone, ответная — в ResponseInfo.
func TestPlotSeriesTimezone(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Lognex-Accept-Timezone") != "Wed, 16 Aug 2017 23:07:01 +0700" {
			t.Errorf("timezone = %q", r.Header.Get("X-Lognex-Accept-Timezone"))
		}
		w.Header().Set("X-Lognex-Content-Timezone", "Wed, 16 Aug 2017 23:07:01 +0700")
		w.Write([]byte(`{"meta":{"href":"h"},"credit":0,"debit":0,"series":[]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	now := time.Date(2017, 8, 16, 23, 7, 1, 0, time.FixedZone("", 7*3600))
	options := money.PlotSeriesOptions{MomentFrom: model.DateTime(now), MomentTo: model.DateTime(now), Interval: money.IntervalHour}
	got, err := service(c).PlotSeries(t.Context(), options, transport.WithAcceptTimezone(now), transport.WithResponseInfo(&info))
	if err != nil || got.Series == nil || len(got.Series) != 0 {
		t.Fatalf("report = %+v, err = %v", got, err)
	}
	if tz, ok := info.ContentTimezone(); !ok || !tz.Equal(now) {
		t.Fatalf("content timezone = %v, %v", tz, ok)
	}
}

// Обязательные momentFrom, momentTo и interval проверяются до HTTP; ResponseInfo сброшен.
func TestPlotSeriesValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	from, to := documenttest.DateTime("2018-09-01 00:00:00"), documenttest.DateTime("2018-09-04 00:00:00")
	for want, options := range map[string]money.PlotSeriesOptions{
		"momentFrom is required": {MomentTo: to, Interval: money.IntervalDay},
		"momentTo is required":   {MomentFrom: from, Interval: money.IntervalDay},
		"interval is required":   {MomentFrom: from, MomentTo: to},
		"interval is set twice":  {MomentFrom: from, MomentTo: to, Interval: money.IntervalDay, Extra: url.Values{"interval": {"hour"}}},
		"filter is set twice":    {MomentFrom: from, MomentTo: to, Interval: money.IntervalDay, Filter: "project=p", Extra: url.Values{"filter": {"x"}}},
		"async is not supported": {MomentFrom: from, MomentTo: to, Interval: money.IntervalDay, Extra: url.Values{"async": {"true"}}},
	} {
		info := transport.ResponseInfo{StatusCode: 299}
		_, err := service(c).PlotSeries(t.Context(), options, transport.WithResponseInfo(&info))
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get money plotseries") ||
			!strings.Contains(strings.ReplaceAll(err.Error(), "parameter ", ""), want) || info.StatusCode != 0 {
			t.Errorf("%s: err = %v, info = %+v", want, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := money.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
