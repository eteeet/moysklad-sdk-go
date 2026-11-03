package facturein_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/facturein"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	otherID   = "13131313-1313-1313-1313-131313131313"
	fixtureID = "209fe91a-ab00-11e6-8af5-581e00000076"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/" + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *facturein.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := facturein.NewService(c)
	return s
}

var (
	organization = ref("organization/11111111-1111-1111-1111-111111111111", "organization")
	supply       = ref("supply/88888888-8888-8888-8888-888888888888", "supply")
	fullCreate   = facturein.CreateRequest{
		Organization: organization,
		Name:         model.Value("00001"), Description: model.Value(""), Code: model.Value("c-1"), ExternalCode: model.Value("ext-1"),
		SyncID:     model.Value("22222222-2222-2222-2222-222222222222"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Agent:          model.Value(ref("counterparty/33333333-3333-3333-3333-333333333333", "counterparty")),
		Contract:       model.Null[model.Reference](),
		Owner:          model.Value(ref("employee/44444444-4444-4444-4444-444444444444", "employee")),
		Group:          model.Value(ref("group/55555555-5555-5555-5555-555555555555", "group")),
		State:          model.Value(ref("facturein/metadata/states/66666666-6666-6666-6666-666666666666", "state")),
		Attributes:     model.Value([]model.Attribute{}),
		Rate:           model.Value(model.DocumentRateRequest{Currency: model.Value(ref("currency/77777777-7777-7777-7777-777777777777", "currency"))}),
		Supplies:       model.Value([]model.Reference{supply}),
		Payments:       model.Value([]model.Reference{}),
		IncomingNumber: model.Value("356124365"),
		IncomingDate:   model.Value(documenttest.DateTime("2017-04-06 00:00:00")),
	}
	fullUpdate = facturein.UpdateRequest{
		Name:           model.Value("FactureIn2"),
		Organization:   model.Value(ref("organization/99999999-9999-9999-9999-999999999999", "organization")),
		State:          model.Null[model.Reference](),
		Applicable:     model.Value(true),
		Payments:       model.Value([]model.Reference{ref("paymentout/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "paymentout")}),
		IncomingNumber: model.Value(""),
		IncomingDate:   model.Null[model.DateTime](),
	}
)

func TestFactureInFixture(t *testing.T) {
	for _, path := range []string{"custom/facturein/read.json", "upstream/facture_in.json"} {
		data := read(t, path)
		for _, codec := range documenttest.Codecs {
			t.Run(path+"/"+codec.Name, func(t *testing.T) {
				var got facturein.FactureIn
				if err := codec.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if got.ID != fixtureID || got.Meta.Type != "facturein" || got.Name != "00002" || got.Sum != 200000 || !got.Applicable || !got.Printed ||
					got.Organization == nil || got.Agent == nil || got.Owner == nil || got.Group == nil || got.Rate == nil || got.Rate.Currency == nil ||
					got.Moment.IsZero() || got.Updated.IsZero() || len(got.Supplies) != 1 || got.Supplies[0].Meta.Type != "supply" {
					t.Fatalf("document = %+v", got)
				}
				if path == "upstream/facture_in.json" && (got.IncomingNumber != "356124365" || got.IncomingDate.IsZero() || got.State == nil ||
					len(got.Attributes) != 1 || len(got.Payments) != 1 || got.Payments[0].Meta.Type != "paymentout") {
					t.Fatalf("document = %+v", got)
				}
			})
		}
	}
	// Поля, которых нет в примерах источника.
	const rest = `{"id":"x","code":"c","syncId":"s","description":"d","created":"2017-04-27 15:55:04","deleted":"2017-04-28 10:00:00",
		"contract":{"meta":{"href":"h","type":"contract"}},"files":{"meta":{"href":"f","size":0}},"future":1}`
	for _, codec := range documenttest.Codecs {
		var got facturein.FactureIn
		if err := codec.Unmarshal([]byte(rest), &got); err != nil || got.Code != "c" || got.SyncID != "s" || got.Description != "d" ||
			got.Created.IsZero() || got.Deleted.IsZero() || got.Contract == nil || got.Files == nil || got.Files.Rows != nil ||
			got.Organization != nil || got.Supplies != nil || got.Payments != nil || !got.IncomingDate.IsZero() {
			t.Fatalf("%s: document = %+v, err = %v", codec.Name, got, err)
		}
	}
}

func TestFactureInRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/facturein/create-full.golden.json")},
			{"create minimal", facturein.CreateRequest{Organization: organization}, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}}}`)},
			{"update", fullUpdate, read(t, "custom/facturein/update.golden.json")},
			{"update empty supplies", facturein.UpdateRequest{Supplies: model.Value([]model.Reference{})}, []byte(`{"supplies":[]}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[facturein.CreateRequest](), reflect.TypeFor[facturein.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "Files", "Project", "Store", "Positions"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestFactureInHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":4000,"incomingNumber":"356124365"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*facturein.FactureIn, error)
	}{
		{"get", "GET", "entity/facturein/" + testID, nil, func(ctx context.Context, c *transport.Client) (*facturein.FactureIn, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/facturein", read(t, "custom/facturein/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*facturein.FactureIn, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/facturein/" + testID, read(t, "custom/facturein/update.golden.json"), func(ctx context.Context, c *transport.Client) (*facturein.FactureIn, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *facturein.FactureIn) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 4000 || got.IncomingNumber != "356124365" {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// exchange выполняет один вызов против сервера с ответом status/response и проверяет запрос:
// метод, путь, тело (в want "{base}" — адрес сервера; "" — тела нет) и отсутствие повторов.
func exchange(t *testing.T, method, path, want string, status int, response string, invoke func(context.Context, *transport.Client) error) error {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		body, _ := io.ReadAll(r.Body)
		if r.Method != method || r.URL.Path != "/"+path || r.URL.RawQuery != "" {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		if want == "" {
			if len(body) != 0 {
				t.Errorf("body = %s", body)
			}
		} else {
			documenttest.JSONEqual(t, body, []byte(strings.ReplaceAll(want, "{base}", "http://"+r.Host)))
		}
		w.WriteHeader(status)
		io.WriteString(w, response)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
	if err != nil {
		t.Fatal(err)
	}
	err = invoke(t.Context(), client)
	if calls.Load() != 1 {
		t.Errorf("calls = %d, want 1", calls.Load())
	}
	return err
}

// noResult приводит вызов без модели ответа к форме documenttest.LocalErrors.
func noResult(invoke func(context.Context, *transport.Client) error) func(context.Context, *transport.Client) (*struct{}, error) {
	return func(ctx context.Context, c *transport.Client) (*struct{}, error) {
		if err := invoke(ctx, c); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	}
}

// Delete, Trash, New и массовые операции: успех, ошибки ответа с контекстом операции, один запрос.
func TestFactureInWrites(t *testing.T) {
	const minimal = `{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}}}`
	metaOf := func(id string) string {
		return `{"meta":{"href":"{base}/entity/facturein/` + id + `","metadataHref":"{base}/entity/facturein/metadata","type":"facturein","mediaType":"application/json"}`
	}
	var (
		created []batch.Result[facturein.FactureIn]
		deleted []batch.Deleted
		tmpl    *facturein.FactureIn
	)
	for _, op := range []struct {
		name, method, path, body, success, context string
		invoke                                     func(context.Context, *transport.Client) error
		check                                      func() bool
	}{
		{"delete", "DELETE", "entity/facturein/" + testID, "", "", "delete facturein", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}, func() bool { return true }},
		{"trash", "POST", "entity/facturein/" + testID + "/trash", "", "", "trash facturein", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}, func() bool { return true }},
		{"new empty", "PUT", "entity/facturein/new", "", `{"name":"00001","applicable":true}`, "get facturein template", func(ctx context.Context, c *transport.Client) (err error) {
			tmpl, err = service(c).New(ctx, nil)
			return err
		}, func() bool { return tmpl != nil && tmpl.ID == "" && tmpl.Name == "00001" }},
		{"new basis", "PUT", "entity/facturein/new", string(read(t, "custom/facturein/new-basis.golden.json")), `{"supplies":[{"meta":{"href":"s","type":"supply"}}]}`, "get facturein template",
			func(ctx context.Context, c *transport.Client) (err error) {
				tmpl, err = service(c).New(ctx, map[string]any{"supplies": []model.Reference{supply}})
				return err
			}, func() bool { return tmpl != nil && len(tmpl.Supplies) == 1 }},
		{"create many", "POST", "entity/facturein", "[" + minimal + "," + minimal + "]", `[{"id":"a"},{"errors":[{"error":"bad","code":3000}]}]`, "create many facturein",
			func(ctx context.Context, c *transport.Client) (err error) {
				created, err = service(c).CreateMany(ctx, []facturein.CreateRequest{{Organization: organization}, {Organization: organization}})
				return err
			}, func() bool {
				return len(created) == 2 && created[0].Value.ID == "a" && created[1].Value == nil && created[1].Errors[0].Code == 3000
			}},
		{"update many", "POST", "entity/facturein", "[" + metaOf(testID) + `,"name":"a"},` + metaOf(otherID) + `,"incomingNumber":""}]`, `[{"id":"a"},{"id":"b"}]`, "update many facturein",
			func(ctx context.Context, c *transport.Client) (err error) {
				created, err = service(c).UpdateMany(ctx, []batch.Update[facturein.UpdateRequest]{
					{ID: testID, Request: facturein.UpdateRequest{Name: model.Value("a")}},
					{ID: otherID, Request: facturein.UpdateRequest{IncomingNumber: model.Value("")}},
				})
				return err
			}, func() bool { return len(created) == 2 && created[1].Value.ID == "b" }},
		{"delete many", "POST", "entity/facturein/delete", "[" + metaOf(testID) + "}," + metaOf(otherID) + "}]", `[{"info":"ok"},{"info":"ok"}]`, "delete many facturein",
			func(ctx context.Context, c *transport.Client) (err error) {
				deleted, err = service(c).DeleteMany(ctx, []string{testID, otherID})
				return err
			}, func() bool { return len(deleted) == 2 && deleted[1].Info == "ok" }},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, c := range []struct {
				status   int
				response string
				ok       bool
			}{
				{200, op.success, true},
				{204, "", op.success == ""},
				{400, `{"errors":[{"error":"bad","code":3000}]}`, false},
				{500, `{"errors":[{"error":"server","code":1000}]}`, false},
				{202, op.success, false},
			} {
				err := exchange(t, op.method, op.path, op.body, c.status, c.response, op.invoke)
				if (err == nil) != c.ok || (c.ok && !op.check()) {
					t.Fatalf("status %d: err = %v", c.status, err)
				}
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), op.context)) {
					t.Errorf("status %d: error = %v", c.status, err)
				}
			}
			documenttest.LocalErrors(t, noResult(op.invoke))
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestFactureInValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no org": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, facturein.CreateRequest{Supplies: model.Value([]model.Reference{supply})}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, facturein.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":  func(o ...transport.RequestOption) error { return s.Trash(ctx, testID+"/x", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[facturein.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[facturein.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{testID, "bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "facturein") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := facturein.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestFactureInList(t *testing.T) {
	listtest.Run(t, "facturein", read(t, "upstream/facture_in.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[facturein.FactureIn] {
		return service(c).List
	}, func(t *testing.T, got facturein.FactureIn) {
		if got.ID != fixtureID || len(got.Supplies) != 1 || len(got.Payments) != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
