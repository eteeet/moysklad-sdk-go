package query_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const attrHref = "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes/0c1d1b7b-1c1e-11e6-9464-e4de00000001"

func TestFilter(t *testing.T) {
	for _, tc := range []struct {
		name       string
		conditions []query.Condition
		want       string
	}{
		{"operators", []query.Condition{
			query.Eq("name", "a"), query.NotEq("code", "b"), query.Less("sum", "1"), query.LessOrEq("moment", "2026-01-01 00:00:00"),
			query.Greater("qty", "0"), query.GreaterOrEq("updated", "2026-01-01 00:00:00"),
			query.Contains("description", "c"), query.StartsWith("article", "d"), query.EndsWith("externalCode", "e"),
		}, "name=a;code!=b;sum<1;moment<=2026-01-01 00:00:00;qty>0;updated>=2026-01-01 00:00:00;description~c;article~=d;externalCode=~e"},
		{"OR by repeated =", []query.Condition{query.Eq("archived", "true"), query.Eq("archived", "false")}, "archived=true;archived=false"},
		{"deleted documents too", []query.Condition{query.Eq("isDeleted", "true"), query.Eq("isDeleted", "false")}, "isDeleted=true;isDeleted=false"},
		{"repeated like", []query.Condition{query.Contains("name", "a"), query.Contains("name", "b")}, "name~a;name~b"},
		{"inner backslash", []query.Condition{query.Eq("name", `a\b`)}, `name=a\b`},
		{"AND by repeated !=", []query.Condition{query.NotEq("name", "a"), query.NotEq("name", "b")}, "name!=a;name!=b"},
		{"range", []query.Condition{query.GreaterOrEq("moment", "1"), query.LessOrEq("moment", "2")}, "moment>=1;moment<=2"},
		{"empty", []query.Condition{query.IsEmpty("code"), query.NotEmpty("article")}, "code=;article!="},
		{"reference and attribute", []query.Condition{
			query.Eq("agent", "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/3a1d1b7b-1c1e-11e6-9464-e4de00000001"),
			query.Eq("id", "3a1d1b7b-1c1e-11e6-9464-e4de00000001"),
			query.Eq(attrHref, "red"),
		}, "agent=https://api.moysklad.ru/api/remap/1.2/entity/counterparty/3a1d1b7b-1c1e-11e6-9464-e4de00000001;id=3a1d1b7b-1c1e-11e6-9464-e4de00000001;" + attrHref + "=red"},
		{"escape semicolon only", []query.Condition{query.Eq("name", `a;b\c%&+ Ёж`)}, `name=a\;b\c%&+ Ёж`},
		{"raw", []query.Condition{query.RawCondition("product!~x"), query.RawCondition(`name=a\;b`), query.Eq("name", "a")}, `product!~x;name=a\;b;name=a`},
		{"none", nil, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := query.Filter(tc.conditions...)
			if err != nil || got != tc.want {
				t.Fatalf("Filter = %q, %v; want %q", got, err, tc.want)
			}
		})
	}
}

func TestFilterRejectsAmbiguous(t *testing.T) {
	for name, conditions := range map[string][]query.Condition{
		"= with > (docs example)": {query.Eq("sum", "100"), query.Greater("sum", "99")},
		"<= with =":               {query.LessOrEq("sum", "1"), query.Eq("sum", "2")},
		"= with !=":               {query.Eq("name", "a"), query.NotEq("name", "b")},
		"!= with =":               {query.NotEq("name", "a"), query.Eq("name", "b")},
		"repeated >":              {query.Greater("sum", "99"), query.Greater("sum", "100")},
		"repeated <=":             {query.LessOrEq("sum", "1"), query.LessOrEq("sum", "2")},
		"trailing backslash":      {query.Eq("name", `a\`), query.Eq("code", "b")},
		"empty field":             {query.Eq("", "a")},
		"field with ;":            {query.Eq("a;b", "c")},
		"empty raw":               {query.RawCondition("")},
		"raw with ;":              {query.RawCondition("a=1;b=2")},
	} {
		if got, err := query.Filter(conditions...); err == nil || !strings.HasPrefix(err.Error(), "moysklad: ") {
			t.Errorf("%s: Filter = %q, %v; want moysklad error", name, got, err)
		}
	}
}

func TestOrder(t *testing.T) {
	got, err := query.Order(query.Asc("name"), query.Desc("moment"), query.Asc(attrHref))
	if err != nil || got != "name,asc;moment,desc;"+attrHref+",asc" {
		t.Fatalf("Order = %q, %v", got, err)
	}
	if got, err := query.Order(query.Desc("updated")); err != nil || got != "updated,desc" {
		t.Fatalf("Order = %q, %v", got, err)
	}
	for _, sorts := range [][]query.Sort{{query.Asc("")}, {query.Asc("a,b")}, {query.Desc("a;b")}} {
		if got, err := query.Order(sorts...); err == nil || !strings.HasPrefix(err.Error(), "moysklad: ") {
			t.Errorf("Order(%v) = %q, want error", sorts, got)
		}
	}
}

// Значение filter кодируется в URL ровно один раз, и сервер получает исходную строку.
func TestFilterEncodedOnce(t *testing.T) {
	filter, err := query.Filter(query.Eq("name", `a;b\c%&+=# Ёж`), query.Eq("agent", "https://x.example/a?b=c"))
	if err != nil {
		t.Fatal(err)
	}
	order, err := query.Order(query.Desc("name"))
	if err != nil {
		t.Fatal(err)
	}
	var got url.Values
	var raw string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, raw = r.URL.Query(), r.URL.RawQuery
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	params := url.Values{"filter": {filter}, "order": {order}}
	if _, err := client.Send(t.Context(), http.MethodGet, "entity/product", params, nil, nil); err != nil {
		t.Fatal(err)
	}
	const want = "filter=name%3Da%5C%3Bb%5Cc%25%26%2B%3D%23+%D0%81%D0%B6%3Bagent%3Dhttps%3A%2F%2Fx.example%2Fa%3Fb%3Dc&order=name%2Cdesc"
	if got.Get("filter") != filter || got.Get("order") != order || raw != want {
		t.Fatalf("filter = %q, raw = %q", got.Get("filter"), raw)
	}
}
