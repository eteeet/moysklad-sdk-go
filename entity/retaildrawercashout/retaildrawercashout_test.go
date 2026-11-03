package retaildrawercashout_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashout"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "1c5afb86-9603-11e6-8a84-bae500000079"
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
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *retaildrawercashout.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := retaildrawercashout.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	agent        = ref("entity/employee/22222222-2222-2222-2222-222222222222", "employee")
	shift        = ref("entity/retailshift/33333333-3333-3333-3333-333333333333", "retailshift")
	minimal      = retaildrawercashout.CreateRequest{Organization: organization, Agent: agent, RetailShift: shift}
	fullCreate   = retaildrawercashout.CreateRequest{
		Organization: organization, Agent: agent, RetailShift: shift,
		OrganizationBranch: model.Value(ref("entity/organizationbranch/44444444-4444-4444-4444-444444444444", "organizationbranch")),
		Name:               model.Value("DOC-1"), Description: model.Value(""), Code: model.Value("C-1"), SyncID: model.Value("55555555-5555-5555-5555-555555555555"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Rate:       model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/66666666-6666-6666-6666-666666666666", "currency"))}),
		Attributes: model.Value([]model.Attribute{}),
		Sum:        model.Value(100500.0),
	}
	fullUpdate = retaildrawercashout.UpdateRequest{
		Name: model.Value("0000004"), Applicable: model.Value(false),
		RetailShift: model.Value(ref("entity/retailshift/99999999-9999-9999-9999-999999999999", "retailshift")),
		State:       model.Null[model.Reference](),
		Sum:         model.Value(0.0),
	}
)

func TestRetailDrawerCashOutFixture(t *testing.T) {
	data := read(t, "custom/retaildrawercashout/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got retaildrawercashout.RetailDrawerCashOut
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Sum != 12441200 || got.Organization == nil || got.OrganizationBranch == nil ||
				got.RetailShift == nil || got.RetailShift.Meta.Type != "retailshift" || got.Rate == nil || got.Code != "" || got.Created.IsZero() || !got.Printed {
				t.Fatalf("document = %+v", got)
			}
			var sparse retaildrawercashout.RetailDrawerCashOut
			if err := codec.Unmarshal([]byte(`{"id":"x","future":1}`), &sparse); err != nil ||
				sparse.RetailShift != nil || sparse.Agent != nil || sparse.Files != nil || !sparse.Deleted.IsZero() {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRetailDrawerCashOutRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/retaildrawercashout/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/22222222-2222-2222-2222-222222222222","type":"employee"}},"retailShift":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailshift/33333333-3333-3333-3333-333333333333","type":"retailshift"}}}`)},
			{"update", fullUpdate, read(t, "custom/retaildrawercashout/update.golden.json")},
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
	// sum помечен «Только для чтения», но пишется: исключение утверждено пользователем.
	for _, typ := range []reflect.Type{reflect.TypeFor[retaildrawercashout.CreateRequest](), reflect.TypeFor[retaildrawercashout.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Files", "Project"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestRetailDrawerCashOutHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":100500}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*retaildrawercashout.RetailDrawerCashOut, error)
	}{
		{"get", "GET", "entity/retaildrawercashout/" + testID, nil, func(ctx context.Context, c *transport.Client) (*retaildrawercashout.RetailDrawerCashOut, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/retaildrawercashout", read(t, "custom/retaildrawercashout/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*retaildrawercashout.RetailDrawerCashOut, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/retaildrawercashout/" + testID, read(t, "custom/retaildrawercashout/update.golden.json"), func(ctx context.Context, c *transport.Client) (*retaildrawercashout.RetailDrawerCashOut, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *retaildrawercashout.RetailDrawerCashOut) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 100500 {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// exchange — один вызов против сервера со статусом status и телом response: запрос, число
// обращений, ошибка. BASE в ожидаемых телах — адрес тестового сервера.
func exchange(t *testing.T, status int, response string, invoke func(context.Context, *transport.Client) error) (string, string, int32, error) {
	t.Helper()
	var calls atomic.Int32
	var request, body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		data, _ := io.ReadAll(r.Body)
		request, body = r.Method+" "+r.URL.RequestURI(), string(data)
		w.WriteHeader(status)
		io.WriteString(w, response)
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
	if err != nil {
		t.Fatal(err)
	}
	err = invoke(t.Context(), c)
	return request, strings.ReplaceAll(body, server.URL, "BASE"), calls.Load(), err
}

// Delete, Trash, New и массовые операции: маршрут, тело, успех, ошибка без повтора.
func TestRetailDrawerCashOutWrites(t *testing.T) {
	meta := func(id string) string {
		return `{"meta":{"href":"BASE/entity/retaildrawercashout/` + id + `","metadataHref":"BASE/entity/retaildrawercashout/metadata","type":"retaildrawercashout","mediaType":"application/json"}`
	}
	basis := map[string]any{"retailShift": shift}
	for _, op := range []struct {
		name, request, body, response, context string
		invoke                                 func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE /entity/retaildrawercashout/" + testID, "", "", "delete retaildrawercashout", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}},
		{"trash", "POST /entity/retaildrawercashout/" + testID + "/trash", "", "", "trash retaildrawercashout", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}},
		{"new", "PUT /entity/retaildrawercashout/new", "", `{"sum":0}`, "retaildrawercashout template", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).New(ctx, nil)
			if err == nil && got.ID != "" {
				t.Errorf("template = %+v", got)
			}
			return err
		}},
		{"new basis", "PUT /entity/retaildrawercashout/new", `{"retailShift":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailshift/33333333-3333-3333-3333-333333333333","type":"retailshift"}}}`,
			`{"retailShift":{"meta":{"href":"h","type":"retailshift"}}}`, "retaildrawercashout template", func(ctx context.Context, c *transport.Client) error {
				got, err := service(c).New(ctx, basis)
				if err == nil && got.RetailShift == nil {
					t.Errorf("template = %+v", got)
				}
				return err
			}},
		{"create many", "POST /entity/retaildrawercashout", "[" + string(read(t, "custom/retaildrawercashout/create-full.golden.json")) + "]", `[{"id":"a"}]`, "create many retaildrawercashout", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).CreateMany(ctx, []retaildrawercashout.CreateRequest{fullCreate})
			if err == nil && got[0].Value.ID != "a" {
				t.Errorf("result = %+v", got)
			}
			return err
		}},
		{"update many", "POST /entity/retaildrawercashout", `[` + meta(testID) + `,"name":"0000004","applicable":false}]`, `[{"id":"a"}]`, "update many retaildrawercashout", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).UpdateMany(ctx, []batch.Update[retaildrawercashout.UpdateRequest]{{ID: testID, Request: retaildrawercashout.UpdateRequest{Name: model.Value("0000004"), Applicable: model.Value(false)}}})
			return err
		}},
		{"delete many", "POST /entity/retaildrawercashout/delete", `[` + meta(testID) + `}]`, `[{"info":"ok"}]`, "delete many retaildrawercashout", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).DeleteMany(ctx, []string{testID})
			if err == nil && got[0].Info != "ok" {
				t.Errorf("result = %+v", got)
			}
			return err
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			request, body, calls, err := exchange(t, 200, op.response, op.invoke)
			if err != nil || calls != 1 || request != op.request {
				t.Fatalf("request = %s, calls = %d, err = %v", request, calls, err)
			}
			if op.body == "" && body != "" {
				t.Errorf("body = %s", body)
			} else if op.body != "" {
				documenttest.JSONEqual(t, []byte(body), []byte(op.body))
			}
			_, _, calls, err = exchange(t, 500, `{"errors":[{"error":"server","code":1000}]}`, op.invoke)
			if err == nil || calls != 1 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), op.context) {
				t.Errorf("server error: calls = %d, err = %v", calls, err)
			}
		})
	}
	if _, _, calls, err := exchange(t, 204, "", func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }); err != nil || calls != 1 {
		t.Errorf("trash 204: calls = %d, err = %v", calls, err)
	}
	if _, _, _, err := exchange(t, 200, "null", func(ctx context.Context, c *transport.Client) error {
		_, err := service(c).New(ctx, nil)
		return err
	}); err == nil {
		t.Error("null template accepted")
	}
}

// Пустые обязательные ссылки, неверные id и пустые массовые запросы отклоняются до запроса;
// ResponseInfo сбрасывается.
func TestRetailDrawerCashOutValidation(t *testing.T) {
	without := func(edit func(*retaildrawercashout.CreateRequest)) retaildrawercashout.CreateRequest {
		input := minimal
		edit(&input)
		return input
	}
	for name, invoke := range map[string]func(context.Context, *transport.Client, ...transport.RequestOption) error{
		"no organization": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Create(ctx, without(func(r *retaildrawercashout.CreateRequest) { r.Organization = model.Reference{} }), o...)
			return err
		},
		"no agent": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Create(ctx, without(func(r *retaildrawercashout.CreateRequest) { r.Agent = model.Reference{} }), o...)
			return err
		},
		"no retail shift": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Create(ctx, without(func(r *retaildrawercashout.CreateRequest) { r.RetailShift = model.Reference{} }), o...)
			return err
		},
		"empty update": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Update(ctx, testID, retaildrawercashout.UpdateRequest{}, o...)
			return err
		},
		"bad get id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Get(ctx, "bad", o...)
			return err
		},
		"bad update id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Update(ctx, "bad", fullUpdate, o...)
			return err
		},
		"bad delete id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).Delete(ctx, testID+"/x", o...)
		},
		"bad trash id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).Trash(ctx, "bad", o...)
		},
		"empty create many": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).CreateMany(ctx, nil, o...)
			return err
		},
		"empty update many item": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).UpdateMany(ctx, []batch.Update[retaildrawercashout.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"bad delete many id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		_, _, calls, err := exchange(t, 200, `{"id":"x"}`, func(ctx context.Context, c *transport.Client) error {
			return invoke(ctx, c, transport.WithResponseInfo(&info))
		})
		if err == nil || calls != 0 || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retaildrawercashout") {
			t.Errorf("%s: calls = %d, info = %+v, err = %v", name, calls, info, err)
		}
	}
	if _, err := retaildrawercashout.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestRetailDrawerCashOutList(t *testing.T) {
	listtest.Run(t, "retaildrawercashout", read(t, "custom/retaildrawercashout/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[retaildrawercashout.RetailDrawerCashOut] {
		return service(c).List
	}, func(t *testing.T, got retaildrawercashout.RetailDrawerCashOut) {
		if got.ID != fixtureID || got.RetailShift == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
