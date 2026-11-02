package assortment_test

import (
	"context"
	"encoding/json/jsontext"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/assortment"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
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

func service(c *transport.Client) *assortment.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := assortment.NewService(c)
	return s
}

// Каждая строка типизирована по meta.type, неизвестный тип остаётся только в Raw,
// а запись Row возвращает строку байт в байт.
func TestRows(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var rows model.Collection[assortment.Row]
			if err := codec.Unmarshal(read(t, "custom/assortment/list.json"), &rows); err != nil {
				t.Fatal(err)
			}
			if len(rows.Rows) != 6 {
				t.Fatalf("rows = %d", len(rows.Rows))
			}
			p, s, b, v := rows.Rows[0], rows.Rows[1], rows.Rows[2], rows.Rows[4]
			if p.Type != "product" || p.Product == nil || p.Product.Name != "product" || p.Service != nil ||
				s.Type != "service" || s.Service == nil || s.Service.ID != "437f2d67-36e7-11e7-8a7f-40d0000000df" ||
				b.Type != "bundle" || b.Bundle == nil || b.Bundle.Name != "bundle" ||
				v.Type != "variant" || v.Variant == nil || v.Variant.Product.Meta.Href == "" || len(v.Variant.Characteristics) == 0 {
				t.Fatalf("rows = %+v", rows.Rows)
			}
			for _, row := range rows.Rows {
				if row.Stock != nil || row.Quantity != nil || row.Reserve != nil || row.InTransit != nil {
					t.Fatalf("stock in new endpoint row: %+v", row)
				}
				var raw map[string]jsontext.Value
				if err := codec.Unmarshal(row.Raw, &raw); err != nil || raw["id"] == nil {
					t.Fatalf("raw = %s, err = %v", row.Raw, err)
				}
				data, err := codec.Marshal(row)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, data, row.Raw)
			}

			var consignments model.Collection[assortment.Row]
			if err := codec.Unmarshal(read(t, "custom/assortment/consignment.json"), &consignments); err != nil {
				t.Fatal(err)
			}
			c := consignments.Rows[1].Consignment
			if consignments.Rows[1].Type != "consignment" || c == nil || c.Name != "продукт / consignment" || c.Label == "" ||
				c.Assortment.Meta.Type != "product" {
				t.Fatalf("consignment = %+v", consignments.Rows[1])
			}

			var unknown assortment.Row
			future := []byte(`{"meta":{"href":"h","type":"future"},"id":"x","extra":[1,{"a":null}]}`)
			if err := codec.Unmarshal(future, &unknown); err != nil || unknown.Type != "future" || unknown.Product != nil ||
				unknown.Service != nil || unknown.Bundle != nil || unknown.Variant != nil || unknown.Consignment != nil {
				t.Fatalf("unknown = %+v, err = %v", unknown, err)
			}
			if data, err := codec.Marshal(unknown); err != nil || string(data) != string(future) {
				t.Fatalf("marshal = %s, err = %v", data, err)
			}
			if err := codec.Unmarshal([]byte(`{"meta":{"type":"product"},"name":1}`), &unknown); err == nil {
				t.Fatal("malformed product row accepted")
			}
		})
	}
}

// Строка legacy endpoint несёт остатки.
func TestLegacyStock(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var row assortment.Row
			if err := codec.Unmarshal(read(t, "upstream/assortment.json"), &row); err != nil || row.Product == nil ||
				row.Stock == nil || row.Reserve == nil || row.InTransit == nil || row.Quantity == nil {
				t.Fatalf("row = %+v, err = %v", row, err)
			}
		})
	}
}

func lists(s *assortment.Service) map[string]listtest.ListFunc[assortment.Row] {
	return map[string]listtest.ListFunc[assortment.Row]{
		"list": func(ctx context.Context, o query.ListOptions, opts ...transport.RequestOption) (*model.Collection[assortment.Row], error) {
			return s.List(ctx, assortment.Options{ListOptions: o}, opts...)
		},
		"legacy": func(ctx context.Context, o query.ListOptions, opts ...transport.RequestOption) (*model.Collection[assortment.Row], error) {
			return s.ListLegacy(ctx, assortment.Options{ListOptions: o}, opts...)
		},
	}
}

// Строка нового endpoint — из примера источника без остатков, legacy — upstream-фикстура с остатками.
func TestListContract(t *testing.T) {
	var example struct {
		Rows []jsontext.Value `json:"rows"`
	}
	if err := documenttest.Codecs[0].Unmarshal(read(t, "custom/assortment/list.json"), &example); err != nil {
		t.Fatal(err)
	}
	rows := map[string][]byte{"list": example.Rows[0], "legacy": read(t, "upstream/assortment.json")}
	for _, name := range []string{"list", "legacy"} {
		t.Run(name, func(t *testing.T) {
			listtest.Run(t, "assortment", rows[name], "", "archived", func(c *transport.Client) listtest.ListFunc[assortment.Row] {
				return lists(service(c))[name]
			}, func(t *testing.T, got assortment.Row) {
				if got.Type != "product" || got.Product == nil {
					t.Fatalf("row = %+v", got)
				}
			})
		})
	}
}

// List всегда включает assortmentWithoutStock, поверх чужого beta-тега вызова; ListLegacy
// заголовок не передаёт. groupBy уходит как есть, скрытых запросов нет.
func TestListModes(t *testing.T) {
	for _, tc := range []struct {
		name  string
		call  func(*assortment.Service, context.Context) error
		beta  string
		query string
	}{
		{"list", func(s *assortment.Service, ctx context.Context) error {
			_, err := s.List(ctx, assortment.Options{GroupBy: assortment.GroupByConsignment}, transport.WithBetaFeature("other"))
			return err
		}, assortment.BetaFeature, "groupBy=consignment"},
		{"legacy", func(s *assortment.Service, ctx context.Context) error {
			_, err := s.ListLegacy(ctx, assortment.Options{GroupBy: assortment.GroupByProduct})
			return err
		}, "", "groupBy=product"},
		{"list default", func(s *assortment.Service, ctx context.Context) error {
			_, err := s.List(ctx, assortment.Options{})
			return err
		}, assortment.BetaFeature, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if got := r.Header.Values(transport.HeaderBetaFeature); tc.beta == "" && len(got) != 0 || tc.beta != "" && (len(got) != 1 || got[0] != tc.beta) {
					t.Errorf("beta = %q", got)
				}
				if r.URL.Path != "/entity/assortment" || r.URL.RawQuery != tc.query {
					t.Errorf("request = %s", r.URL)
				}
				w.Write(read(t, "custom/assortment/list.json"))
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			if err := tc.call(service(c), t.Context()); err != nil || calls.Load() != 1 {
				t.Fatalf("err = %v, calls = %d", err, calls.Load())
			}
		})
	}
}

// All: каждая страница — с assortmentWithoutStock поверх чужого beta-тега вызова, следующая — GET
// по nextHref; ResponseInfo — последней страницы. Ошибка первой страницы — один элемент, обход
// не начинается; break после первой страницы не делает второй запрос.
func TestAll(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if got := r.Header.Values(transport.HeaderBetaFeature); len(got) != 1 || got[0] != assortment.BetaFeature {
			t.Errorf("beta = %q", got)
		}
		href := "http://" + r.Host + "/entity/assortment?limit=1&search=tea&offset="
		switch r.URL.RawQuery {
		case "limit=1&search=tea":
			w.Header().Set("X-Page", "1")
			fmt.Fprintf(w, `{"meta":{"href":%q,"nextHref":%q},"rows":[{"meta":{"type":"product"},"id":"p1"}]}`, href+"0", href+"1")
		case "limit=1&search=tea&offset=1":
			w.Header().Set("X-Page", "2")
			fmt.Fprintf(w, `{"meta":{"href":%q},"rows":[{"meta":{"type":"service"},"id":"s1"}]}`, href+"1")
		default:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":[{"error":"auth","code":1056}]}`))
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	options := assortment.Options{ListOptions: query.ListOptions{Limit: 1, Search: "tea"}}

	var ids []string
	var info transport.ResponseInfo
	for page, err := range s.All(t.Context(), options, transport.WithBetaFeature("other"), transport.WithResponseInfo(&info)) {
		if err != nil {
			t.Fatal(err)
		}
		for _, row := range page.Rows {
			ids = append(ids, string(row.Type))
		}
	}
	if fmt.Sprint(ids) != "[product service]" || calls.Load() != 2 || info.Header.Get("X-Page") != "2" {
		t.Fatalf("ids = %v, calls = %d, info = %+v", ids, calls.Load(), info)
	}

	calls.Store(0)
	for range s.All(t.Context(), options) {
		break
	}
	if calls.Load() != 1 {
		t.Fatalf("break: calls = %d", calls.Load())
	}

	calls.Store(0)
	var errs int
	for page, err := range s.All(t.Context(), assortment.Options{}) {
		if page != nil || sdkerrors.StatusCode(err) != http.StatusUnauthorized {
			t.Fatalf("page = %v, err = %v", page, err)
		}
		errs++
	}
	if errs != 1 || calls.Load() != 1 {
		t.Fatalf("first page error: errs = %d, calls = %d", errs, calls.Load())
	}
}

// Регрессия ревью 2b91a89: beta-тег нового режима из настроек клиента или опций вызова
// не переключает ListLegacy молча — ошибка до запроса; другой тег проходит.
func TestListLegacyRejectsBetaMode(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Header.Get(transport.HeaderBetaFeature) == assortment.BetaFeature {
			t.Errorf("legacy request with %s", assortment.BetaFeature)
		}
		w.Write([]byte(`{"meta":{"href":"h"},"rows":[]}`))
	}))
	defer server.Close()
	for name, tc := range map[string]struct {
		config http.Header
		opts   []transport.RequestOption
		calls  int32
	}{
		"config header":   {http.Header{transport.HeaderBetaFeature: {assortment.BetaFeature}}, nil, 0},
		"config csv":      {http.Header{transport.HeaderBetaFeature: {"other, " + assortment.BetaFeature}}, nil, 0},
		"WithBetaFeature": {nil, []transport.RequestOption{transport.WithBetaFeature(assortment.BetaFeature)}, 0},
		"WithHeader":      {nil, []transport.RequestOption{transport.WithHeader(transport.HeaderBetaFeature, assortment.BetaFeature)}, 0},
		"call replaces":   {http.Header{transport.HeaderBetaFeature: {assortment.BetaFeature}}, []transport.RequestOption{transport.WithBetaFeature("other")}, 1},
		"other tag":       {http.Header{transport.HeaderBetaFeature: {"other"}}, nil, 1},
	} {
		t.Run(name, func(t *testing.T) {
			calls.Store(0)
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Header: tc.config})
			if err != nil {
				t.Fatal(err)
			}
			info := transport.ResponseInfo{StatusCode: 503}
			_, err = service(c).ListLegacy(t.Context(), assortment.Options{}, append(tc.opts, transport.WithResponseInfo(&info))...)
			if calls.Load() != tc.calls || (tc.calls == 0) != (err != nil) || tc.calls == 0 && info.StatusCode != 0 {
				t.Fatalf("calls = %d, err = %v, info = %+v", calls.Load(), err, info)
			}
		})
	}
}

func TestListRejects(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	for name, options := range map[string]assortment.Options{
		"fields":        {ListOptions: query.ListOptions{Fields: "stock"}},
		"extra groupBy": {ListOptions: query.ListOptions{Extra: url.Values{"groupBy": {"product"}}}},
		"limit":         {ListOptions: query.ListOptions{Limit: 1001}},
	} {
		for mode, list := range lists(s) {
			if _, err := list(t.Context(), options.ListOptions); err == nil {
				t.Errorf("%s %s accepted", mode, name)
			}
		}
		if _, err := s.List(t.Context(), options); err == nil {
			t.Errorf("list %s accepted", name)
		}
	}
	if _, err := assortment.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
