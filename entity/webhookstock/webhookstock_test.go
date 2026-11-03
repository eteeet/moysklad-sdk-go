package webhookstock_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/webhookstock"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "7944ef04-f831-11e5-7a69-971500188b19"
	stockPath = "entity/webhookstock/" + testID
	appHref   = "https://api.moysklad.ru/api/remap/1.2/entity/application/9e1ad712-3e45-4679-8896-7159973a8ef5"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *webhookstock.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := webhookstock.NewService(c)
	return s
}

var (
	full = webhookstock.CreateRequest{
		URL: "https://example.com/hooks/stock", StockType: "stock", ReportType: "bystore", Enabled: model.Value(false),
	}
	update = webhookstock.UpdateRequest{
		URL: model.Value("https://www.example.com"), StockType: model.Value("stock"), ReportType: model.Value("all"), Enabled: model.Value(true),
	}
)

func TestWebhookStockFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			// Upstream SDK: authorApplication с полями решения — читается только meta.
			var got webhookstock.WebhookStock
			if err := codec.Unmarshal(read(t, "upstream/webhookstock.json"), &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/webhookstock/" + testID,
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/webhookstock/metadata",
				Type:         "webhookstock", MediaType: "application/json",
			}
			if got.Meta != wantMeta || got.ID != testID || got.AccountID != "b8b74698-9128-11e6-8a84-bae500000001" ||
				got.StockType != "stock" || got.ReportType != "bystore" || got.URL != "http://www.example.com" || !got.Enabled {
				t.Fatalf("webhookstock = %+v", got)
			}
			if got.AuthorApplication == nil || got.AuthorApplication.Meta.Href != appHref || got.AuthorApplication.Meta.Type != "application" {
				t.Fatalf("authorApplication = %+v", got.AuthorApplication)
			}

			var page model.Collection[webhookstock.WebhookStock]
			if err := codec.Unmarshal(read(t, "custom/webhookstock/list.json"), &page); err != nil {
				t.Fatal(err)
			}
			if page.Meta.Size != 3 || len(page.Rows) != 3 {
				t.Fatalf("page = %+v", page)
			}
			first, second := page.Rows[0], page.Rows[1]
			if first.AuthorApplication == nil || first.AuthorApplication.Meta.Href != appHref || second.AuthorApplication != nil ||
				second.ID != "d08f9217-bbd2-11e6-8a84-bae500000004" || second.AccountID != "8afc8c88-38a5-11ed-ac14-000f00000001" ||
				second.StockType != "stock" || second.ReportType != "all" || second.URL != "https://www.example.com" || !second.Enabled {
				t.Fatalf("rows = %+v", page.Rows)
			}

			// Явный false приходит в ответе «Отключить вебхук на изменение остатков».
			disabled := webhookstock.WebhookStock{Enabled: true}
			if err := codec.Unmarshal(read(t, "custom/webhookstock/disabled.json"), &disabled); err != nil ||
				disabled.ID != testID || disabled.Enabled || disabled.ReportType != "bystore" {
				t.Fatalf("disabled = %+v, err = %v", disabled, err)
			}

			// Открытые enum и неизвестные поля не отвергаются.
			var open webhookstock.WebhookStock
			if err := codec.Unmarshal([]byte(`{"id":"x","stockType":"reserve","reportType":"byslot","futureField":{"a":1}}`), &open); err != nil ||
				open.StockType != "reserve" || open.ReportType != "byslot" || open.AuthorApplication != nil || open.Enabled || open.URL != "" {
				t.Fatalf("open = %+v, err = %v", open, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/webhookstock/create-full.golden.json")},
			{"update", update, read(t, "custom/webhookstock/update.golden.json")},
			{"create minimal", webhookstock.CreateRequest{URL: "https://www.example.com", StockType: "stock", ReportType: "all"},
				[]byte(`{"url":"https://www.example.com","stockType":"stock","reportType":"all"}`)},
			{"create enabled true", webhookstock.CreateRequest{URL: "u", StockType: "s", ReportType: "r", Enabled: model.Value(true)},
				[]byte(`{"url":"u","stockType":"s","reportType":"r","enabled":true}`)},
			{"update disable", webhookstock.UpdateRequest{Enabled: model.Value(false)}, []byte(`{"enabled":false}`)},
			{"update unset", webhookstock.UpdateRequest{}, []byte(`{}`)},
			{"update nulls", webhookstock.UpdateRequest{ReportType: model.Null[string](), Enabled: model.Null[bool]()}, []byte(`{"reportType":null,"enabled":null}`)},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, tc.want)
			})
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[webhookstock.CreateRequest](), reflect.TypeFor[webhookstock.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "AuthorApplication"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestServiceHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса.
	const response = `{"id":"server-id","stockType":"stock","reportType":"bystore","url":"https://server.example","enabled":false}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*webhookstock.WebhookStock, error)
	}{
		{"get", "GET", stockPath, nil, func(ctx context.Context, c *transport.Client) (*webhookstock.WebhookStock, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/webhookstock", read(t, "custom/webhookstock/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*webhookstock.WebhookStock, error) {
			return service(c).Create(ctx, full)
		}},
		// Значения вне перечислений и url без схемы уходят как есть: проверяет сервер.
		{"create open values", "POST", "entity/webhookstock", []byte(`{"url":"not a url","stockType":"reserve","reportType":"byslot"}`),
			func(ctx context.Context, c *transport.Client) (*webhookstock.WebhookStock, error) {
				return service(c).Create(ctx, webhookstock.CreateRequest{URL: "not a url", StockType: "reserve", ReportType: "byslot"})
			}},
		{"update", "PUT", stockPath, read(t, "custom/webhookstock/update.golden.json"), func(ctx context.Context, c *transport.Client) (*webhookstock.WebhookStock, error) {
			return service(c).Update(ctx, testID, update)
		}},
		{"update disable", "PUT", stockPath, []byte(`{"enabled":false}`), func(ctx context.Context, c *transport.Client) (*webhookstock.WebhookStock, error) {
			return service(c).Update(ctx, testID, webhookstock.UpdateRequest{Enabled: model.Value(false)})
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *webhookstock.WebhookStock) {
						if got.ID != "server-id" || got.StockType != "stock" || got.ReportType != "bystore" ||
							got.URL != "https://server.example" || got.Enabled {
							t.Fatalf("webhookstock = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/webhookstock"
	stockMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"webhookstock","mediaType":"application/json"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", stockPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"create many", "POST", "entity/webhookstock", `[{"url":"https://www.example.com","stockType":"stock","reportType":"all"},{"url":"u","stockType":"stock","reportType":"bystore","enabled":false}]`,
			`[{"id":"x"},{"errors":[{"error":"bad","code":3000}]}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []webhookstock.CreateRequest{
					{URL: "https://www.example.com", StockType: "stock", ReportType: "all"},
					{URL: "u", StockType: "stock", ReportType: "bystore", Enabled: model.Value(false)},
				})
				if err == nil && (len(r) != 2 || r[0].Value == nil || r[0].Value.ID != "x" || r[1].Value != nil || len(r[1].Errors) != 1 || r[1].Errors[0].Code != 3000) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/webhookstock", `[{"meta":` + stockMeta + `,"url":"https://www.example.com","stockType":"stock","reportType":"bystore"}]`,
			`[{"id":"` + testID + `","reportType":"bystore"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[webhookstock.UpdateRequest]{{ID: testID, Request: webhookstock.UpdateRequest{
					URL: model.Value("https://www.example.com"), StockType: model.Value("stock"), ReportType: model.Value("bystore")}}})
				if err == nil && (len(r) != 1 || r[0].Value == nil || r[0].Value.ID != testID || r[0].Value.ReportType != "bystore") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/webhookstock/delete", `[{"meta":` + stockMeta + `}]`, `[{"info":"ok"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).DeleteMany(ctx, []string{testID})
				if err == nil && (len(r) != 1 || r[0].Info != "ok") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, st := range []struct {
				status int
				ok     bool
			}{{200, true}, {204, tc.empty}, {201, false}, {403, false}, {503, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != tc.method || r.URL.Path != "/api/remap/1.2/"+tc.path || r.URL.RawQuery != "" {
						t.Errorf("request = %s %s", r.Method, r.URL)
					}
					if tc.body == "" {
						if len(b) != 0 {
							t.Errorf("body = %s", b)
						}
					} else {
						documenttest.JSONEqual(t, b, []byte(strings.ReplaceAll(tc.body, "{base}", server.URL+"/api/remap/1.2")))
					}
					w.WriteHeader(st.status)
					switch st.status {
					case 200:
						io.WriteString(w, tc.response)
					case 403:
						io.WriteString(w, `{"errors":[{"error":"forbidden","code":1016}]}`)
					}
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.invoke(t.Context(), client)
				server.Close()
				if (err == nil) != st.ok || calls.Load() != 1 {
					t.Fatalf("status %d: err = %v, calls = %d", st.status, err, calls.Load())
				}
				var apiErr *sdkerrors.APIError
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "webhookstock") ||
					st.status == 403 && !sdkerrors.HasCode(err, 1016) || errors.As(err, &apiErr) != (st.status >= 400)) {
					t.Errorf("status %d: error = %v", st.status, err)
				}
			}
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
				if err := tc.invoke(ctx, c); err != nil {
					return nil, err
				}
				return &struct{}{}, nil
			})
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestServiceValidation(t *testing.T) {
	if s, err := webhookstock.NewService(nil); s != nil || err == nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	create := func(in webhookstock.CreateRequest) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error { _, err := s.Create(ctx, in, o...); return err }
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":          func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":            func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no url":       create(webhookstock.CreateRequest{StockType: "stock", ReportType: "all"}),
		"create no stock":     create(webhookstock.CreateRequest{URL: "u", ReportType: "all"}),
		"create no report":    create(webhookstock.CreateRequest{URL: "u", StockType: "stock"}),
		"create only enabled": create(webhookstock.CreateRequest{Enabled: model.Value(true)}),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, webhookstock.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[webhookstock.UpdateRequest]{{ID: "bad", Request: update}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[webhookstock.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "webhookstock") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "webhookstock", read(t, "upstream/webhookstock.json"), "", "", func(c *transport.Client) listtest.ListFunc[webhookstock.WebhookStock] {
		return service(c).List
	}, func(t *testing.T, got webhookstock.WebhookStock) {
		if got.ID != testID || got.StockType != "stock" || got.ReportType != "bystore" || !got.Enabled || got.AuthorApplication == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
