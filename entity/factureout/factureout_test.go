package factureout_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/factureout"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	otherID   = "13131313-1313-1313-1313-131313131313"
	fixtureID = "02eb1458-aa6e-11e6-8a84-bc5200000076"
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

func service(c *transport.Client) *factureout.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := factureout.NewService(c)
	return s
}

var (
	organization = ref("organization/11111111-1111-1111-1111-111111111111", "organization")
	demand       = ref("demand/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", "demand")
	fullCreate   = factureout.CreateRequest{
		Organization: organization,
		Name:         model.Value("00001"), Description: model.Value(""), Code: model.Value("c-1"), ExternalCode: model.Value("ext-1"),
		SyncID:     model.Value("22222222-2222-2222-2222-222222222222"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Agent:             model.Value(ref("counterparty/33333333-3333-3333-3333-333333333333", "counterparty")),
		Contract:          model.Null[model.Reference](),
		Owner:             model.Value(ref("employee/44444444-4444-4444-4444-444444444444", "employee")),
		Group:             model.Value(ref("group/55555555-5555-5555-5555-555555555555", "group")),
		State:             model.Value(ref("factureout/metadata/states/66666666-6666-6666-6666-666666666666", "state")),
		Attributes:        model.Value([]model.Attribute{}),
		Rate:              model.Value(model.DocumentRateRequest{Currency: model.Value(ref("currency/77777777-7777-7777-7777-777777777777", "currency"))}),
		StateContractID:   model.Value("005674"),
		AdvancePaymentVat: model.Value(0),
		PaymentPurpose:    model.Value("Авансовый платеж по договору №123"),
		Payments:          model.Value([]model.Reference{ref("paymentin/88888888-8888-8888-8888-888888888888", "paymentin")}),
		Demands:           model.Value([]model.Reference{}),
		Consignee:         model.Value(organization),
		ConsigneeBranch:   model.Null[model.Reference](),
		PaymentNumber:     model.Value("платежный д-т"),
		PaymentDate:       model.Value(documenttest.DateTime("2016-11-14 00:00:00")),
	}
	fullUpdate = factureout.UpdateRequest{
		Name:              model.Value("FactureOut2"),
		Organization:      model.Value(ref("organization/99999999-9999-9999-9999-999999999999", "organization")),
		State:             model.Null[model.Reference](),
		Applicable:        model.Value(true),
		StateContractID:   model.Value(""),
		AdvancePaymentVat: model.Value(20),
		Returns:           model.Value([]model.Reference{ref("purchasereturn/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "purchasereturn")}),
		PaymentDate:       model.Null[model.DateTime](),
	}
)

func TestFactureOutFixture(t *testing.T) {
	for _, path := range []string{"custom/factureout/read.json", "upstream/facture_out.json"} {
		data := read(t, path)
		for _, codec := range documenttest.Codecs {
			t.Run(path+"/"+codec.Name, func(t *testing.T) {
				var got factureout.FactureOut
				if err := codec.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if got.Meta.Type != "factureout" || got.Name != "00001" || got.Sum != 0 || !got.Applicable || !got.Printed ||
					got.Organization == nil || got.Agent == nil || got.Owner == nil || got.Group == nil || got.State == nil || got.Rate == nil || got.Rate.Currency == nil ||
					got.Moment.IsZero() || got.Updated.IsZero() || len(got.Attributes) != 1 || got.StateContractID != "005674" || got.Consignee == nil ||
					len(got.Demands) != 1 || got.Demands[0].Meta.Type != "demand" || got.Payments != nil || got.Returns != nil {
					t.Fatalf("document = %+v", got)
				}
				if path == "upstream/facture_out.json" && (got.ID != fixtureID || got.Contract == nil || got.Rate.Value == nil || *got.Rate.Value != 60 ||
					got.PaymentNumber != "платежный д-т" || got.PaymentDate.IsZero() || got.Description != "комментарий") {
					t.Fatalf("document = %+v", got)
				}
			})
		}
	}
	// Поля, которых нет в примерах источника.
	const rest = `{"id":"x","code":"c","syncId":"s","description":"d","created":"2017-04-27 18:27:09","deleted":"2017-04-28 10:00:00",
		"files":{"meta":{"href":"f","size":0}},"advancePaymentVat":20,"paymentPurpose":"p","vatSum":1666.67,
		"payments":[{"meta":{"href":"p","type":"paymentin"}}],"returns":[{"meta":{"href":"r","type":"purchasereturn"}}],
		"consigneeBranch":{"meta":{"href":"b","type":"organizationbranch"}},"future":1}`
	for _, codec := range documenttest.Codecs {
		var got factureout.FactureOut
		if err := codec.Unmarshal([]byte(rest), &got); err != nil || got.Code != "c" || got.SyncID != "s" || got.Description != "d" ||
			got.Created.IsZero() || got.Deleted.IsZero() || got.Files == nil || got.Files.Rows != nil || got.AdvancePaymentVat != 20 ||
			got.PaymentPurpose != "p" || got.VATSum != 1666.67 || len(got.Payments) != 1 || len(got.Returns) != 1 || got.ConsigneeBranch == nil ||
			got.Organization != nil || got.Demands != nil || !got.PaymentDate.IsZero() {
			t.Fatalf("%s: document = %+v, err = %v", codec.Name, got, err)
		}
	}
}

func TestFactureOutRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/factureout/create-full.golden.json")},
			{"create minimal", factureout.CreateRequest{Organization: organization}, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}}}`)},
			{"update", fullUpdate, read(t, "custom/factureout/update.golden.json")},
			{"update empty demands", factureout.UpdateRequest{Demands: model.Value([]model.Reference{})}, []byte(`{"demands":[]}`)},
			{"update vat 0", factureout.UpdateRequest{AdvancePaymentVat: model.Value(0)}, []byte(`{"advancePaymentVat":0}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[factureout.CreateRequest](), reflect.TypeFor[factureout.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum", "Files", "Project", "Store", "Positions"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestFactureOutHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":10000,"vatSum":1666.67,"stateContractId":"005674"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*factureout.FactureOut, error)
	}{
		{"get", "GET", "entity/factureout/" + testID, nil, func(ctx context.Context, c *transport.Client) (*factureout.FactureOut, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/factureout", read(t, "custom/factureout/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*factureout.FactureOut, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/factureout/" + testID, read(t, "custom/factureout/update.golden.json"), func(ctx context.Context, c *transport.Client) (*factureout.FactureOut, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *factureout.FactureOut) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 10000 || got.VATSum != 1666.67 || got.StateContractID != "005674" {
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
func TestFactureOutWrites(t *testing.T) {
	const minimal = `{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}}}`
	metaOf := func(id string) string {
		return `{"meta":{"href":"{base}/entity/factureout/` + id + `","metadataHref":"{base}/entity/factureout/metadata","type":"factureout","mediaType":"application/json"}`
	}
	var (
		created []batch.Result[factureout.FactureOut]
		deleted []batch.Deleted
		tmpl    *factureout.FactureOut
	)
	for _, op := range []struct {
		name, method, path, body, success, context string
		invoke                                     func(context.Context, *transport.Client) error
		check                                      func() bool
	}{
		{"delete", "DELETE", "entity/factureout/" + testID, "", "", "delete factureout", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}, func() bool { return true }},
		{"trash", "POST", "entity/factureout/" + testID + "/trash", "", "", "trash factureout", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}, func() bool { return true }},
		{"new empty", "PUT", "entity/factureout/new", "", `{"name":"00001","applicable":true}`, "get factureout template", func(ctx context.Context, c *transport.Client) (err error) {
			tmpl, err = service(c).New(ctx, nil)
			return err
		}, func() bool { return tmpl != nil && tmpl.ID == "" && tmpl.Name == "00001" }},
		{"new basis", "PUT", "entity/factureout/new", string(read(t, "custom/factureout/new-basis.golden.json")), `{"demands":[{"meta":{"href":"d","type":"demand"}}]}`, "get factureout template",
			func(ctx context.Context, c *transport.Client) (err error) {
				tmpl, err = service(c).New(ctx, map[string]any{"demands": []model.Reference{demand}})
				return err
			}, func() bool { return tmpl != nil && len(tmpl.Demands) == 1 }},
		{"create many", "POST", "entity/factureout", "[" + minimal + "," + minimal + "]", `[{"id":"a"},{"errors":[{"error":"bad","code":3000}]}]`, "create many factureout",
			func(ctx context.Context, c *transport.Client) (err error) {
				created, err = service(c).CreateMany(ctx, []factureout.CreateRequest{{Organization: organization}, {Organization: organization}})
				return err
			}, func() bool {
				return len(created) == 2 && created[0].Value.ID == "a" && created[1].Value == nil && created[1].Errors[0].Code == 3000
			}},
		{"update many", "POST", "entity/factureout", "[" + metaOf(testID) + `,"name":"a"},` + metaOf(otherID) + `,"stateContractId":""}]`, `[{"id":"a"},{"id":"b"}]`, "update many factureout",
			func(ctx context.Context, c *transport.Client) (err error) {
				created, err = service(c).UpdateMany(ctx, []batch.Update[factureout.UpdateRequest]{
					{ID: testID, Request: factureout.UpdateRequest{Name: model.Value("a")}},
					{ID: otherID, Request: factureout.UpdateRequest{StateContractID: model.Value("")}},
				})
				return err
			}, func() bool { return len(created) == 2 && created[1].Value.ID == "b" }},
		{"delete many", "POST", "entity/factureout/delete", "[" + metaOf(testID) + "}," + metaOf(otherID) + "}]", `[{"info":"ok"},{"info":"ok"}]`, "delete many factureout",
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
func TestFactureOutValidation(t *testing.T) {
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
			_, err := s.Create(ctx, factureout.CreateRequest{Demands: model.Value([]model.Reference{demand})}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, factureout.UpdateRequest{}, o...)
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
			_, err := s.UpdateMany(ctx, []batch.Update[factureout.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[factureout.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{testID, "bad"}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "factureout") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := factureout.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestFactureOutList(t *testing.T) {
	listtest.Run(t, "factureout", read(t, "upstream/facture_out.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[factureout.FactureOut] {
		return service(c).List
	}, func(t *testing.T, got factureout.FactureOut) {
		if got.ID != fixtureID || len(got.Demands) != 1 || got.PaymentNumber == "" {
			t.Fatalf("row = %+v", got)
		}
	})
}
