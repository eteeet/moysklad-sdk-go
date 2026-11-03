package counterpartyadjustment_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/counterpartyadjustment"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	otherID   = "13131313-1313-1313-1313-131313131313"
	fixtureID = "1bb899bc-0b17-11ec-ac16-000a00000006"
)

const minimal = `{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
	`"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/33333333-3333-3333-3333-333333333333","type":"employee"}}}`

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

func service(c *transport.Client) *counterpartyadjustment.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := counterpartyadjustment.NewService(c)
	return s
}

var (
	organization = ref("organization/11111111-1111-1111-1111-111111111111", "organization")
	agent        = ref("employee/33333333-3333-3333-3333-333333333333", "employee")
	fullCreate   = counterpartyadjustment.CreateRequest{
		Organization: organization, Agent: agent,
		Name: model.Value("00001"), Description: model.Value(""), ExternalCode: model.Value("ext-1"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Owner:      model.Value(ref("employee/44444444-4444-4444-4444-444444444444", "employee")),
		Group:      model.Null[model.Reference](),
		Attributes: model.Value([]model.Attribute{}),
	}
	fullUpdate = counterpartyadjustment.UpdateRequest{
		Agent:       model.Value(ref("counterparty/99999999-9999-9999-9999-999999999999", "counterparty")),
		Description: model.Value(""),
		Applicable:  model.Value(true),
		Owner:       model.Null[model.Reference](),
	}
)

func TestCounterpartyAdjustmentFixture(t *testing.T) {
	for _, path := range []string{"custom/counterpartyadjustment/read.json", "upstream/counterparty_adjustment.json"} {
		data := read(t, path)
		for _, codec := range documenttest.Codecs {
			t.Run(path+"/"+codec.Name, func(t *testing.T) {
				var got counterpartyadjustment.CounterpartyAdjustment
				if err := codec.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if got.Meta.Type != "counterpartyadjustment" || got.Meta.UUIDHref == "" || !got.Applicable || got.Printed || got.Published ||
					got.Organization == nil || got.Agent == nil || got.Agent.Meta.Type != "counterparty" || got.Owner == nil || got.Group == nil ||
					got.Moment.IsZero() || got.Updated.IsZero() || got.Created.IsZero() || got.Files == nil || got.Files.Rows != nil || got.ExternalCode == "" {
					t.Fatalf("document = %+v", got)
				}
				// Суммы примеров — копейки как есть: 12.5 в ответе Get, 12.0 в upstream.
				if path == "upstream/counterparty_adjustment.json" && (got.ID != fixtureID || got.Sum != 12 || got.Description != "description" || len(got.Attributes) != 1) ||
					path != "upstream/counterparty_adjustment.json" && (got.Sum != 12.5 || got.Name != "00004") {
					t.Fatalf("document = %+v", got)
				}
			})
		}
	}
	// Знак суммы источник не описывает: отрицательное значение читается без изменений.
	for _, codec := range documenttest.Codecs {
		var got counterpartyadjustment.CounterpartyAdjustment
		if err := codec.Unmarshal([]byte(`{"id":"x","sum":-1250,"deleted":"2021-09-02 10:00:00.000","future":1}`), &got); err != nil ||
			got.Sum != -1250 || got.Deleted.IsZero() || got.Organization != nil || got.Files != nil || got.Attributes != nil {
			t.Fatalf("%s: document = %+v, err = %v", codec.Name, got, err)
		}
	}
}

func TestCounterpartyAdjustmentRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/counterpartyadjustment/create-full.golden.json")},
			{"create minimal", counterpartyadjustment.CreateRequest{Organization: organization, Agent: agent}, []byte(minimal)},
			{"update", fullUpdate, read(t, "custom/counterpartyadjustment/update.golden.json")},
			{"update applicable false", counterpartyadjustment.UpdateRequest{Applicable: model.Value(false)}, []byte(`{"applicable":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[counterpartyadjustment.CreateRequest](), reflect.TypeFor[counterpartyadjustment.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "Files", "Project", "Store", "Positions", "Code", "SyncID", "Rate", "State"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestCounterpartyAdjustmentHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":-12.5}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*counterpartyadjustment.CounterpartyAdjustment, error)
	}{
		{"get", "GET", "entity/counterpartyadjustment/" + testID, nil, func(ctx context.Context, c *transport.Client) (*counterpartyadjustment.CounterpartyAdjustment, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/counterpartyadjustment", read(t, "custom/counterpartyadjustment/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*counterpartyadjustment.CounterpartyAdjustment, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/counterpartyadjustment/" + testID, read(t, "custom/counterpartyadjustment/update.golden.json"), func(ctx context.Context, c *transport.Client) (*counterpartyadjustment.CounterpartyAdjustment, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *counterpartyadjustment.CounterpartyAdjustment) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != -12.5 {
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

// Delete, Trash и массовые операции: успех, ошибки ответа с контекстом операции, один запрос.
func TestCounterpartyAdjustmentWrites(t *testing.T) {
	metaOf := func(id string) string {
		return `{"meta":{"href":"{base}/entity/counterpartyadjustment/` + id + `","metadataHref":"{base}/entity/counterpartyadjustment/metadata","type":"counterpartyadjustment","mediaType":"application/json"}`
	}
	var (
		created []batch.Result[counterpartyadjustment.CounterpartyAdjustment]
		deleted []batch.Deleted
	)
	for _, op := range []struct {
		name, method, path, body, success, context string
		invoke                                     func(context.Context, *transport.Client) error
		check                                      func() bool
	}{
		{"delete", "DELETE", "entity/counterpartyadjustment/" + testID, "", "", "delete counterpartyadjustment", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}, func() bool { return true }},
		{"trash", "POST", "entity/counterpartyadjustment/" + testID + "/trash", "", "", "trash counterpartyadjustment", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}, func() bool { return true }},
		{"create many", "POST", "entity/counterpartyadjustment", "[" + minimal + "," + minimal + "]", `[{"id":"a"},{"errors":[{"error":"bad","code":3000}]}]`, "create many counterpartyadjustment",
			func(ctx context.Context, c *transport.Client) (err error) {
				created, err = service(c).CreateMany(ctx, []counterpartyadjustment.CreateRequest{{Organization: organization, Agent: agent}, {Organization: organization, Agent: agent}})
				return err
			}, func() bool {
				return len(created) == 2 && created[0].Value.ID == "a" && created[1].Value == nil && created[1].Errors[0].Code == 3000
			}},
		{"update many", "POST", "entity/counterpartyadjustment", "[" + metaOf(testID) + `,"name":"a"},` + metaOf(otherID) + `,"description":""}]`, `[{"id":"a"},{"id":"b"}]`, "update many counterpartyadjustment",
			func(ctx context.Context, c *transport.Client) (err error) {
				created, err = service(c).UpdateMany(ctx, []batch.Update[counterpartyadjustment.UpdateRequest]{
					{ID: testID, Request: counterpartyadjustment.UpdateRequest{Name: model.Value("a")}},
					{ID: otherID, Request: counterpartyadjustment.UpdateRequest{Description: model.Value("")}},
				})
				return err
			}, func() bool { return len(created) == 2 && created[1].Value.ID == "b" }},
		{"delete many", "POST", "entity/counterpartyadjustment/delete", "[" + metaOf(testID) + "}," + metaOf(otherID) + "}]", `[{"info":"ok"},{"info":"ok"}]`, "delete many counterpartyadjustment",
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
func TestCounterpartyAdjustmentValidation(t *testing.T) {
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
			_, err := s.Create(ctx, counterpartyadjustment.CreateRequest{Agent: agent}, o...)
			return err
		},
		"create no agent": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, counterpartyadjustment.CreateRequest{Organization: organization}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, counterpartyadjustment.UpdateRequest{}, o...)
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
			_, err := s.UpdateMany(ctx, []batch.Update[counterpartyadjustment.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[counterpartyadjustment.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{testID, "bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "counterpartyadjustment") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := counterpartyadjustment.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestCounterpartyAdjustmentList(t *testing.T) {
	listtest.Run(t, "counterpartyadjustment", read(t, "upstream/counterparty_adjustment.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[counterpartyadjustment.CounterpartyAdjustment] {
		return service(c).List
	}, func(t *testing.T, got counterpartyadjustment.CounterpartyAdjustment) {
		if got.ID != fixtureID || got.Sum != 12 || got.Agent == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
