package audit_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/audit"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/upstream/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestContextFixture(t *testing.T) {
	data := readFixture(t, "audit.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got audit.Context
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != "83732f4b-792a-11e7-9464-d04800000005" || got.UID != "admin@11" || got.Source != "app" ||
				got.EventType != "create" || got.EntityType != "customerorder" || got.ObjectCount != 1 ||
				time.Time(got.Moment).IsZero() || got.Events == nil || got.Events.Meta.Size != 1 ||
				!strings.Contains(string(got.Info), "loginInfo") {
				t.Fatalf("context = %+v", got)
			}
		})
	}
}

// Неизвестные source/entityType и info произвольной формы сохраняются.
func TestContextOpenValues(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got audit.Context
			data := `{"id":"x","source":"newsource","entityType":"newentity","objectType":"x","info":{"n":12345678901234567890},"supportAccess":true}`
			if err := codec.Unmarshal([]byte(data), &got); err != nil {
				t.Fatal(err)
			}
			if got.Source != "newsource" || got.EntityType != "newentity" || !got.SupportAccess ||
				string(got.Info) != `{"n":12345678901234567890}` || got.Events != nil {
				t.Fatalf("%s: context = %+v", codec.Name, got)
			}
			// Info строкой, как в примере списка источника.
			var str audit.Context
			if err := codec.Unmarshal([]byte(`{"id":"y","info":"Вход в МойСклад (127.0.0.1)"}`), &str); err != nil || string(str.Info) != `"Вход в МойСклад (127.0.0.1)"` {
				t.Fatalf("string info = %s, err = %v", str.Info, err)
			}
			var nulls audit.Context
			if err := codec.Unmarshal([]byte(`{"id":"z","events":null,"info":null}`), &nulls); err != nil || nulls.Events != nil {
				t.Fatalf("nulls = %+v, err = %v", nulls, err)
			}
		})
	}
}

// Ключи каталога в форме таблицы (camelCase) читаются так же, как в форме примера.
func TestFiltersKeyCasing(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, data := range []string{
			`{"source":["jsonapi"],"entitytype":["demand"],"eventtype":["create"]}`,
			`{"source":["jsonapi"],"entityType":["demand"],"eventType":["create"]}`,
		} {
			var got audit.FilterCatalog
			if err := codec.Unmarshal([]byte(data), &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Source) != 1 || len(got.EntityType) != 1 || got.EntityType[0] != "demand" || len(got.EventType) != 1 {
				t.Errorf("%s %s: filters = %+v", codec.Name, data, got)
			}
		}
	}
}

func TestFiltersFixture(t *testing.T) {
	data := readFixture(t, "audit_filters.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got audit.FilterCatalog
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Source) != 18 || len(got.EntityType) != 89 || len(got.EventType) != 14 {
				t.Fatalf("filters = %d/%d/%d", len(got.Source), len(got.EntityType), len(got.EventType))
			}
		})
	}
}

func newService(t *testing.T, url string) *audit.Service {
	t.Helper()
	client, err := transport.NewClient(transport.Config{BaseURL: url})
	if err != nil {
		t.Fatal(err)
	}
	s, err := audit.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestListContexts(t *testing.T) {
	row := readFixture(t, "audit.json")
	page := []byte(`{"context":{"employee":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/employee"}}},` +
		`"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/audit/","type":"audit","size":1,"limit":25,"offset":0},"rows":[` + string(row) + `]}`)
	from := model.DateTime(time.Date(2026, 10, 1, 0, 0, 0, 0, time.FixedZone("MSK", 3*60*60)))
	filter, err := query.Filter(query.GreaterOrEq("moment", from.QueryValue()), query.Eq("source", "jsonapi"),
		query.Eq("application", "https://api.moysklad.ru/api/remap/1.2/entity/application/46ea8005-2965-11e9-9ff4-34e80009ac49"))
	if err != nil {
		t.Fatal(err)
	}
	otherFilter, err := query.Filter(query.LessOrEq("moment", "2026-10-31 23:59:59"), query.Eq("uid", "admin@shop"),
		query.Eq("employee", "https://api.moysklad.ru/api/remap/1.2/entity/employee/33333333-3333-3333-3333-333333333333"),
		query.Eq("eventType", "update"), query.Eq("entityType", "https://api.moysklad.ru/api/remap/1.2/entity/customentity/eaacabaf-2655-11e6-8a84-bae500000045"))
	if err != nil {
		t.Fatal(err)
	}
	otherQuery := "filter=" + url.QueryEscape(otherFilter)
	for _, tc := range []struct {
		name    string
		options audit.ContextsOptions
		query   string
		status  int
		body    []byte
		gzip    bool
		rows    int
		wantErr bool
	}{
		{"defaults", audit.ContextsOptions{}, "", 200, page, false, 1, false},
		{"max limit and filter", audit.ContextsOptions{Limit: 100, Offset: 25, Filter: filter},
			"filter=moment%3E%3D2026-10-01+00%3A00%3A00%3Bsource%3Djsonapi%3Bapplication%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fapplication%2F46ea8005-2965-11e9-9ff4-34e80009ac49&limit=100&offset=25",
			200, page, true, 1, false},
		{"other filter fields", audit.ContextsOptions{Filter: otherFilter}, otherQuery, 200, page, false, 1, false},
		// Страница без context — например, при токене решения: SDK его не требует.
		{"empty without context", audit.ContextsOptions{}, "", 200, []byte(`{"meta":{"href":"h","size":0},"rows":[]}`), false, 0, false},
		{"forbidden", audit.ContextsOptions{}, "", 403, []byte(`{"errors":[{"error":"no audit access","code":1016}]}`), false, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.Path != "/audit" || r.URL.RawQuery != tc.query {
					t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
				}
				body := tc.body
				if tc.gzip {
					var buf bytes.Buffer
					zw := gzip.NewWriter(&buf)
					zw.Write(body)
					zw.Close()
					body = buf.Bytes()
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.WriteHeader(tc.status)
				w.Write(body)
			}))
			defer server.Close()
			var info transport.ResponseInfo
			got, err := newService(t, server.URL).ListContexts(t.Context(), tc.options, transport.WithResponseInfo(&info))
			if calls.Load() != 1 || info.StatusCode != tc.status {
				t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
			}
			if tc.wantErr {
				if got != nil || sdkerrors.StatusCode(err) != tc.status || strings.Count(err.Error(), "moysklad:") != 1 {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				return
			}
			if err != nil || len(got.Rows) != tc.rows {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
			if tc.rows > 0 && (got.Rows[0].ID != "83732f4b-792a-11e7-9464-d04800000005" || got.Context == nil) {
				t.Fatalf("row = %+v", got.Rows[0])
			}
		})
	}
}

// Лимит контекстов — 100, не 1000; неверные опции не отправляются.
func TestListContextsRejects(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	s := newService(t, server.URL)
	for _, options := range []audit.ContextsOptions{{Limit: 101}, {Limit: -1}, {Offset: -1}} {
		if got, err := s.ListContexts(t.Context(), options); got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%+v: got = %+v, err = %v", options, got, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*model.Collection[audit.Context], error) {
		s, err := audit.NewService(c)
		if err != nil {
			return nil, err
		}
		return s.ListContexts(ctx, audit.ContextsOptions{})
	})
}

func TestGetFilters(t *testing.T) {
	fixture := readFixture(t, "audit_filters.json")
	for _, tc := range []struct {
		status int
		body   []byte
	}{{200, fixture}, {403, []byte(`{"errors":[{"error":"denied","code":1016}]}`)}, {200, []byte(`null`)}, {200, []byte(`{}`)}} {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			if r.Method != http.MethodGet || r.URL.RequestURI() != "/audit/metadata/filters" {
				t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
			}
			w.WriteHeader(tc.status)
			w.Write(tc.body)
		}))
		got, err := newService(t, server.URL).GetFilters(t.Context())
		server.Close()
		ok := tc.status == 200 && len(tc.body) > 4
		if calls.Load() != 1 || ok && (err != nil || len(got.EventType) != 14) ||
			!ok && (got != nil || err == nil || tc.status != 200 && sdkerrors.StatusCode(err) != tc.status) {
			t.Fatalf("status %d: got = %+v, err = %v", tc.status, got, err)
		}
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*audit.FilterCatalog, error) {
		s, err := audit.NewService(c)
		if err != nil {
			return nil, err
		}
		return s.GetFilters(ctx)
	})
	if _, err := audit.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}
