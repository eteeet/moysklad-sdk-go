package query_test

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
)

func TestListOptionsValues(t *testing.T) {
	values, err := query.ListOptions{}.Values(query.MaxLimit)
	if err != nil || len(values) != 0 {
		t.Fatalf("zero options = %v, %v; server defaults must not be sent", values, err)
	}

	extra := url.Values{query.ParamOwner: {"https://x/employee/1"}, "tag": {"a", "b"}}
	values, err = query.ListOptions{
		Limit: 1000, Offset: 20, Filter: "archived=true;archived=false", Order: "name,asc",
		Search: "кофе 100%", NamedFilter: "https://x/namedfilter/1", Fields: "stock", Extra: extra,
	}.Values(query.MaxLimit)
	if err != nil {
		t.Fatal(err)
	}
	const want = "fields=stock&filter=archived%3Dtrue%3Barchived%3Dfalse&limit=1000&namedfilter=https%3A%2F%2Fx%2Fnamedfilter%2F1" +
		"&offset=20&order=name%2Casc&owner=https%3A%2F%2Fx%2Femployee%2F1&search=%D0%BA%D0%BE%D1%84%D0%B5+100%25&tag=a&tag=b"
	if values.Encode() != want {
		t.Fatalf("query = %s", values.Encode())
	}
	values.Add("tag", "c")
	if len(extra["tag"]) != 2 {
		t.Fatal("Extra aliased into result")
	}

	if values, err := (query.ListOptions{Limit: 100}).Values(100); err != nil || values.Get("limit") != "100" {
		t.Fatalf("audit-like max: %v, %v", values, err)
	}
}

func TestListOptionsRejects(t *testing.T) {
	for name, tc := range map[string]struct {
		opts query.ListOptions
		max  int
	}{
		"negative limit":   {query.ListOptions{Limit: -1}, 1000},
		"limit over max":   {query.ListOptions{Limit: 1001}, 1000},
		"limit over audit": {query.ListOptions{Limit: 101}, 100},
		"negative offset":  {query.ListOptions{Offset: -1}, 1000},
		"several fields":   {query.ListOptions{Fields: "stock,declaration"}, 1000},
		"expand":           {query.ListOptions{Extra: url.Values{"expand": {"agent"}}}, 1000},
		"async":            {query.ListOptions{Extra: url.Values{"async": {"true"}}}, 1000},
		"limit twice":      {query.ListOptions{Limit: 10, Extra: url.Values{"limit": {"5"}}}, 1000},
		"filter in extra":  {query.ListOptions{Extra: url.Values{"filter": {"a=b"}}}, 1000},
	} {
		if values, err := tc.opts.Values(tc.max); err == nil || !strings.HasPrefix(err.Error(), "moysklad: ") {
			t.Errorf("%s: %v, %v", name, values, err)
		}
	}
}

func TestReportOptionsValues(t *testing.T) {
	moscow := time.FixedZone("MSK", 3*60*60)
	from := model.DateTime(time.Date(2026, 10, 1, 0, 0, 0, 0, moscow))
	to := model.DateTime(time.Date(2026, 10, 31, 20, 59, 59, 500_000_000, time.UTC)) // 23:59:59.500 по Москве
	values, err := query.ReportOptions{MomentFrom: from, MomentTo: to, GroupBy: "variant", Interval: "day",
		Filter: "store=https://x/store/1", Extra: url.Values{"limit": {"10"}, "search": {"кофе"}, "order": {"name"}}}.Values()
	if err != nil {
		t.Fatal(err)
	}
	const want = "filter=store%3Dhttps%3A%2F%2Fx%2Fstore%2F1&groupBy=variant&interval=day&limit=10" +
		"&momentFrom=2026-10-01+00%3A00%3A00&momentTo=2026-10-31+23%3A59%3A59.500&order=name&search=%D0%BA%D0%BE%D1%84%D0%B5"
	if values.Encode() != want {
		t.Fatalf("query = %s", values.Encode())
	}
	if values, err := (query.ReportOptions{}).Values(); err != nil || len(values) != 0 {
		t.Fatalf("zero report options = %v, %v", values, err)
	}
	for name, opts := range map[string]query.ReportOptions{
		"moment twice":  {MomentFrom: from, Extra: url.Values{"momentFrom": {"x"}}},
		"groupBy twice": {Extra: url.Values{"groupBy": {"product"}}, GroupBy: "variant"},
		"expand":        {Extra: url.Values{"expand": {"x"}}},
		"async":         {Extra: url.Values{"async": {"true"}}},
	} {
		if _, err := opts.Values(); err == nil || !strings.HasPrefix(err.Error(), "moysklad: ") {
			t.Errorf("%s: %v", name, err)
		}
	}
}
