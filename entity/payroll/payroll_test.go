package payroll_test

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
	"time"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/payroll"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	fixtureID  = "946f8595-17c3-11f1-4095-ec9f00000001"
	positionID = "defb003f-2bf3-11f1-37a8-e68f0000000d"
	docPath    = "entity/payroll/" + testID
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

func service(c *transport.Client) *payroll.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := payroll.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	employee     = ref("entity/employee/77777777-7777-7777-7777-777777777777", "employee")
	periodStart  = documenttest.DateTime("2026-03-24 00:00:00")
	periodEnd    = documenttest.DateTime("2026-03-25 23:59:59")
	minimal      = payroll.CreateRequest{Organization: organization, StartPayrollPeriod: periodStart, EndPayrollPeriod: periodEnd}
	fullCreate   = payroll.CreateRequest{
		Organization: organization, StartPayrollPeriod: periodStart, EndPayrollPeriod: periodEnd,
		Name: model.Value("00006"), Description: model.Value(""), ExternalCode: model.Value("ext"), SyncID: model.Value("22222222-2222-2222-2222-222222222222"),
		Moment: model.Value(documenttest.DateTime("2026-03-25 19:03:00")), Applicable: model.Value(false), Shared: model.Value(true),
		Owner:      model.Value(ref("entity/employee/55555555-5555-5555-5555-555555555555", "employee")),
		Group:      model.Value(ref("entity/group/66666666-6666-6666-6666-666666666666", "group")),
		State:      model.Null[model.Reference](),
		Attributes: model.Value([]model.Attribute{}),
		Positions:  model.Value([]payroll.PositionRequest{{Employee: model.Value(employee), BaseSalary: model.Value(202000.0), PieceworkSalary: model.Value(0.0)}}),
	}
	fullUpdate = payroll.UpdateRequest{
		EndPayrollPeriod: model.Value(documenttest.DateTime("2026-03-26 23:59:59")),
		Name:             model.Value("putSecond"),
		Group:            model.Null[model.Reference](),
		Positions:        model.Value([]payroll.PositionRequest{{Employee: model.Value(employee), BaseSalary: model.Value(200.0), PieceworkSalary: model.Value(300.0)}}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/payroll/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got payroll.Payroll
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Name != "00001" || got.ExternalCode != "EX0JefwKhGFRrJ7uS95aI0" || got.Sum != 60000 || !got.Applicable || got.Shared || got.Printed || got.Published {
				t.Fatalf("document = %+v", got)
			}
			if !time.Time(got.StartPayrollPeriod).Equal(time.Time(documenttest.DateTime("2026-03-03 00:00:00"))) ||
				!time.Time(got.EndPayrollPeriod).Equal(time.Time(documenttest.DateTime("2026-03-04 23:59:59"))) || got.Created.IsZero() || got.Updated.IsZero() || got.Moment.IsZero() {
				t.Fatalf("dates = %+v", got)
			}
			if got.Organization == nil || got.Owner == nil || got.Group == nil || got.State == nil || got.Files == nil ||
				got.Positions == nil || got.Positions.Meta.Type != "payrollposition" || got.Positions.Meta.Size != 2 {
				t.Fatalf("references = %+v", got)
			}
			if len(got.Payments) != 2 || got.Payments[0].Meta.Type != "paymentout" || got.Payments[1].Meta.Type != "cashout" {
				t.Fatalf("payments = %+v", got.Payments)
			}
			var sparse payroll.Payroll
			if err := codec.Unmarshal([]byte(`{"id":"x","syncId":"s","description":"d","attributes":[],`+
				`"positions":{"meta":{"href":"h"},"rows":[{"meta":{"type":"payrollposition"},"id":"p","employee":{"meta":{"type":"employee"}},"baseSalary":2020,"pieceworkSalary":2020}]},"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			row := sparse.Positions.Rows[0]
			if sparse.SyncID != "s" || sparse.Description != "d" || sparse.Attributes == nil || row.Employee == nil || row.BaseSalary != 2020 || row.PieceworkSalary != 2020 ||
				sparse.Organization != nil || !sparse.Created.IsZero() || sparse.Payments != nil {
				t.Fatalf("sparse = %+v", sparse)
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
			{"create full", fullCreate, read(t, "custom/payroll/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
				`"startPayrollPeriod":"2026-03-24 00:00:00.000","endPayrollPeriod":"2026-03-25 23:59:59.000"}`)},
			{"update", fullUpdate, read(t, "custom/payroll/update.golden.json")},
			{"update empty positions", payroll.UpdateRequest{Positions: model.Value([]payroll.PositionRequest{})}, []byte(`{"positions":[]}`)},
			{"position update", payroll.PositionRequest{BaseSalary: model.Value(20.0)}, []byte(`{"baseSalary":20}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[payroll.CreateRequest](), reflect.TypeFor[payroll.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Printed", "Published", "Sum", "Files", "Payments"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"ID", "Meta", "AccountID"} {
		if _, ok := reflect.TypeFor[payroll.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains read-only %s", name)
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":4040,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*payroll.Payroll, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*payroll.Payroll, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/payroll", read(t, "custom/payroll/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*payroll.Payroll, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/payroll/update.golden.json"), func(ctx context.Context, c *transport.Client) (*payroll.Payroll, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *payroll.Payroll) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 4040 || got.Positions.Rows != nil {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func httpCases[T any](t *testing.T, method, path, body, response string, invoke func(context.Context, *transport.Client) (*T, error), check func(*T)) {
	t.Helper()
	var sent []byte
	if body != "" {
		sent = []byte(body)
	}
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, method, path, sent, invoke, check)
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestPositionsHTTP(t *testing.T) {
	const row = `{"meta":{"type":"payrollposition"},"id":"` + positionID + `","employee":{"meta":{"type":"employee"}},"baseSalary":20,"pieceworkSalary":5000}`
	checkRow := func(got *payroll.Position) {
		if got.ID != positionID || got.BaseSalary != 20 || got.PieceworkSalary != 5000 || got.Employee == nil {
			t.Fatalf("row = %+v", got)
		}
	}
	t.Run("get positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=2&offset=7", "", `{"meta":{"href":"h"},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[payroll.Position], error) {
				return service(c).GetPositions(ctx, testID, positions.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[payroll.Position]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions/"+positionID, "", row, func(ctx context.Context, c *transport.Client) (*payroll.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}, checkRow)
	})
	t.Run("update position", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/positions/"+positionID, `{"baseSalary":20,"pieceworkSalary":5000}`, row, func(ctx context.Context, c *transport.Client) (*payroll.Position, error) {
			return service(c).UpdatePosition(ctx, testID, positionID, payroll.PositionRequest{BaseSalary: model.Value(20.0), PieceworkSalary: model.Value(5000.0)})
		}, checkRow)
	})
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/payroll"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"payroll","mediaType":"application/json"}`
	employeeJSON := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/77777777-7777-7777-7777-777777777777","type":"employee"}}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"create many", "POST", "entity/payroll", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
			`"startPayrollPeriod":"2026-03-24 00:00:00.000","endPayrollPeriod":"2026-03-25 23:59:59.000"}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []payroll.CreateRequest{minimal})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/payroll", `[{"meta":` + docMeta + `,"applicable":false}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).UpdateMany(ctx, []batch.Update[payroll.UpdateRequest]{{ID: testID, Request: payroll.UpdateRequest{Applicable: model.Value(false)}}})
				return err
			}},
		{"delete many", "POST", "entity/payroll/delete", `[{"meta":` + docMeta + `}]`, `[{"info":"ok"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).DeleteMany(ctx, []string{testID})
				return err
			}},
		{"create positions", "POST", docPath + "/positions", `[{"employee":` + employeeJSON + `,"baseSalary":2000,"pieceworkSalary":3000}]`, `[{"id":"` + positionID + `"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreatePositions(ctx, testID, []payroll.PositionRequest{{Employee: model.Value(employee), BaseSalary: model.Value(2000.0), PieceworkSalary: model.Value(3000.0)}})
				if err == nil && r[0].Value.ID != positionID {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete position", "DELETE", docPath + "/positions/" + positionID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, positionID)
		}},
		{"delete many positions", "POST", docPath + "/positions/delete", `[{"meta":{"href":"` + base + `/` + testID + `/positions/` + positionID + `","type":"payrollposition","mediaType":"application/json"}}]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyPositions(ctx, testID, []string{positionID})
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, st := range []struct {
				status int
				ok     bool
			}{{200, true}, {204, tc.empty}, {201, false}, {403, false}, {503, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != tc.method || r.URL.Path != "/api/remap/1.2/"+tc.path {
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "payroll") ||
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

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 1000 строк уходят.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		io.WriteString(w, `{"id":"server-id"}`)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	rows := func(n int) model.Field[[]payroll.PositionRequest] {
		return model.Value(make([]payroll.PositionRequest, n))
	}
	limit := minimal
	limit.Positions = rows(1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, payroll.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	create := func(change func(*payroll.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	getPositions := func(id string, p positions.PageOptions) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			_, err := s.GetPositions(ctx, id, p, o...)
			return err
		}
	}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":      func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":        func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no org":   create(func(r *payroll.CreateRequest) { r.Organization = model.Reference{} }),
		"create no start": create(func(r *payroll.CreateRequest) { r.StartPayrollPeriod = model.DateTime{} }),
		"create no end":   create(func(r *payroll.CreateRequest) { r.EndPayrollPeriod = model.DateTime{} }),
		"create 1001":     create(func(r *payroll.CreateRequest) { r.Positions = rows(1001) }),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, payroll.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update 1001": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, payroll.UpdateRequest{Positions: rows(1001)}, o...)
			return err
		},
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":  func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[payroll.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[payroll.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"positions bad id": getPositions("bad", positions.PageOptions{}),
		"positions fields": getPositions(testID, positions.PageOptions{Fields: "stock"}),
		"positions limit":  getPositions(testID, positions.PageOptions{Limit: 1001}),
		"get position bad": func(o ...transport.RequestOption) error {
			_, err := s.GetPosition(ctx, testID, "bad", o...)
			return err
		},
		"create positions empty": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, testID, nil, o...)
			return err
		},
		"create positions bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, "bad", []payroll.PositionRequest{{}}, o...)
			return err
		},
		"update position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, payroll.PositionRequest{}, o...)
			return err
		},
		"delete position bad": func(o ...transport.RequestOption) error { return s.DeletePosition(ctx, "bad", positionID, o...) },
		"delete many positions empty": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, nil, o...)
		},
		"delete many positions bad row": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, []string{"bad"}, o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "payroll") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := payroll.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "payroll", read(t, "custom/payroll/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[payroll.Payroll] {
		return service(c).List
	}, func(t *testing.T, got payroll.Payroll) {
		if got.ID != fixtureID || len(got.Payments) != 2 {
			t.Fatalf("row = %+v", got)
		}
	})
}
