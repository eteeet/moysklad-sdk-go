package webhook_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/webhook"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID      = "aec51463-bbd2-11e6-8a84-bae500000003"
	webhookPath = "entity/webhook/" + testID
	appHref     = "https://api.moysklad.ru/api/remap/1.2/entity/application/9e1ad712-3e45-4679-8896-7159973a8ef5"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *webhook.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := webhook.NewService(c)
	return s
}

var (
	full = webhook.CreateRequest{
		URL: "https://example.com/hooks/moysklad", Action: "UPDATE", EntityType: "supply",
		DiffType: model.Value("FIELDS"), Enabled: model.Value(false), Method: model.Value("POST"),
	}
	update = webhook.UpdateRequest{
		URL: model.Value("http://www.example.com"), Action: model.Value("DELETE"), EntityType: model.Value("demand"),
		DiffType: model.Value("NONE"), Enabled: model.Value(false),
	}
)

func TestWebhookFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			// Upstream SDK: authorApplication с полями решения — читается только meta.
			var got webhook.Webhook
			if err := codec.Unmarshal(read(t, "upstream/webhook.json"), &got); err != nil {
				t.Fatal(err)
			}
			wantMeta := model.Meta{
				Href:         "https://api.moysklad.ru/api/remap/1.2/entity/webhook/" + testID,
				MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/webhook/metadata",
				Type:         "webhook", MediaType: "application/json",
			}
			if got.Meta != wantMeta || got.ID != testID || got.AccountID != "b8b74698-9128-11e6-8a84-bae500000001" ||
				got.EntityType != "demand" || got.URL != "http://www.example.com" || got.Method != "POST" ||
				!got.Enabled || got.Action != "CREATE" || got.DiffType != "" {
				t.Fatalf("webhook = %+v", got)
			}
			if got.AuthorApplication == nil || got.AuthorApplication.Meta.Href != appHref || got.AuthorApplication.Meta.Type != "application" {
				t.Fatalf("authorApplication = %+v", got.AuthorApplication)
			}

			var page model.Collection[webhook.Webhook]
			if err := codec.Unmarshal(read(t, "custom/webhook/list.json"), &page); err != nil {
				t.Fatal(err)
			}
			if page.Meta.Size != 3 || len(page.Rows) != 3 {
				t.Fatalf("page = %+v", page)
			}
			first, third := page.Rows[0], page.Rows[2]
			if first.AuthorApplication == nil || first.AuthorApplication.Meta.Href != appHref || page.Rows[1].AuthorApplication != nil ||
				page.Rows[1].EntityType != "supply" || third.ID != "d946c7ff-bbd2-11e6-8a84-bae500000005" ||
				third.EntityType != "cashin" || third.Action != "UPDATE" || third.DiffType != "NONE" || !third.Enabled {
				t.Fatalf("rows = %+v", page.Rows)
			}

			var fields webhook.Webhook
			if err := codec.Unmarshal(read(t, "custom/webhook/create-fields.json"), &fields); err != nil ||
				fields.ID != "d08f9217-bbd2-11e6-8a84-bae500000004" || fields.Action != "UPDATE" || fields.DiffType != "FIELDS" {
				t.Fatalf("create fields = %+v, err = %v", fields, err)
			}

			// Явный false приходит в ответе «Отключить вебхук».
			disabled := webhook.Webhook{Enabled: true}
			if err := codec.Unmarshal(read(t, "custom/webhook/disabled.json"), &disabled); err != nil ||
				disabled.ID != testID || disabled.Enabled || disabled.Action != "DELETE" {
				t.Fatalf("disabled = %+v, err = %v", disabled, err)
			}

			// Открытые enum и неизвестные поля не отвергаются.
			var open webhook.Webhook
			if err := codec.Unmarshal([]byte(`{"id":"x","action":"ARCHIVE","entityType":"futureentity","diffType":"FULL","method":"PUT","futureField":{"a":1}}`), &open); err != nil ||
				open.Action != "ARCHIVE" || open.EntityType != "futureentity" || open.DiffType != "FULL" || open.Method != "PUT" ||
				open.AuthorApplication != nil || open.Enabled || open.URL != "" {
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
			{"create full", full, read(t, "custom/webhook/create-full.golden.json")},
			{"update", update, read(t, "custom/webhook/update.golden.json")},
			{"create minimal", webhook.CreateRequest{URL: "http://www.example.com", Action: "CREATE", EntityType: "supply"},
				[]byte(`{"url":"http://www.example.com","action":"CREATE","entityType":"supply"}`)},
			{"create enabled true", webhook.CreateRequest{URL: "u", Action: "PROCESSED", EntityType: "e", Enabled: model.Value(true)},
				[]byte(`{"url":"u","action":"PROCESSED","entityType":"e","enabled":true}`)},
			{"update disable", webhook.UpdateRequest{Enabled: model.Value(false)}, []byte(`{"enabled":false}`)},
			{"update unset", webhook.UpdateRequest{}, []byte(`{}`)},
			{"update nulls", webhook.UpdateRequest{DiffType: model.Null[string](), Enabled: model.Null[bool]()}, []byte(`{"diffType":null,"enabled":null}`)},
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
	// method без маркера readonly пишется; authorApplication в запросах не нужен (его ставит сервер).
	for _, typ := range []reflect.Type{reflect.TypeFor[webhook.CreateRequest](), reflect.TypeFor[webhook.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "AuthorApplication"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestServiceHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса.
	const response = `{"id":"server-id","entityType":"cashin","url":"https://server.example","method":"POST","enabled":false,"action":"UPDATE","diffType":"NONE"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*webhook.Webhook, error)
	}{
		{"get", "GET", webhookPath, nil, func(ctx context.Context, c *transport.Client) (*webhook.Webhook, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/webhook", read(t, "custom/webhook/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*webhook.Webhook, error) {
			return service(c).Create(ctx, full)
		}},
		// Значения вне перечислений и url без схемы уходят как есть: проверяет сервер.
		{"create open values", "POST", "entity/webhook", []byte(`{"url":"not a url","action":"ARCHIVE","entityType":"futureentity","diffType":"FULL"}`),
			func(ctx context.Context, c *transport.Client) (*webhook.Webhook, error) {
				return service(c).Create(ctx, webhook.CreateRequest{URL: "not a url", Action: "ARCHIVE", EntityType: "futureentity", DiffType: model.Value("FULL")})
			}},
		{"update", "PUT", webhookPath, read(t, "custom/webhook/update.golden.json"), func(ctx context.Context, c *transport.Client) (*webhook.Webhook, error) {
			return service(c).Update(ctx, testID, update)
		}},
		{"update disable", "PUT", webhookPath, []byte(`{"enabled":false}`), func(ctx context.Context, c *transport.Client) (*webhook.Webhook, error) {
			return service(c).Update(ctx, testID, webhook.UpdateRequest{Enabled: model.Value(false)})
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *webhook.Webhook) {
						if got.ID != "server-id" || got.EntityType != "cashin" || got.URL != "https://server.example" ||
							got.Enabled || got.Action != "UPDATE" || got.DiffType != "NONE" || got.Method != "POST" {
							t.Fatalf("webhook = %+v", got)
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
	const base = "{base}/entity/webhook"
	webhookMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"webhook","mediaType":"application/json"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", webhookPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"create many", "POST", "entity/webhook", `[{"url":"http://www.example.com","action":"CREATE","entityType":"supply"},{"url":"u","action":"UPDATE","entityType":"demand","enabled":false}]`,
			`[{"id":"x"},{"errors":[{"error":"bad","code":3000}]}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []webhook.CreateRequest{
					{URL: "http://www.example.com", Action: "CREATE", EntityType: "supply"},
					{URL: "u", Action: "UPDATE", EntityType: "demand", Enabled: model.Value(false)},
				})
				if err == nil && (len(r) != 2 || r[0].Value == nil || r[0].Value.ID != "x" || r[1].Value != nil || len(r[1].Errors) != 1 || r[1].Errors[0].Code != 3000) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/webhook", `[{"meta":` + webhookMeta + `,"url":"http://www.example.com","action":"DELETE"}]`, `[{"id":"` + testID + `","action":"DELETE"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[webhook.UpdateRequest]{{ID: testID, Request: webhook.UpdateRequest{URL: model.Value("http://www.example.com"), Action: model.Value("DELETE")}}})
				if err == nil && (len(r) != 1 || r[0].Value == nil || r[0].Value.ID != testID || r[0].Value.Action != "DELETE") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/webhook/delete", `[{"meta":` + webhookMeta + `}]`, `[{"info":"ok"}]`, false,
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "webhook") ||
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
	if s, err := webhook.NewService(nil); s != nil || err == nil {
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
	create := func(in webhook.CreateRequest) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error { _, err := s.Create(ctx, in, o...); return err }
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":         func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":           func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no url":      create(webhook.CreateRequest{Action: "CREATE", EntityType: "supply"}),
		"create no action":   create(webhook.CreateRequest{URL: "u", EntityType: "supply"}),
		"create no entity":   create(webhook.CreateRequest{URL: "u", Action: "CREATE"}),
		"create only fields": create(webhook.CreateRequest{DiffType: model.Value("FIELDS"), Enabled: model.Value(true)}),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, webhook.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[webhook.UpdateRequest]{{ID: "bad", Request: update}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[webhook.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "webhook") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "webhook", read(t, "upstream/webhook.json"), "", "", func(c *transport.Client) listtest.ListFunc[webhook.Webhook] {
		return service(c).List
	}, func(t *testing.T, got webhook.Webhook) {
		if got.ID != testID || got.EntityType != "demand" || got.Action != "CREATE" || !got.Enabled || got.AuthorApplication == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
