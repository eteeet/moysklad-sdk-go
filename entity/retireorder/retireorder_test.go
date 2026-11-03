package retireorder_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/retireorder"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	fixtureID  = "b14bcb5e-3b17-4765-87cf-bc4569fc5f32"
	positionID = "d3040339-a8bd-11ef-ac15-000400000059"
	docPath    = "entity/retireorder/" + testID
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

func service(c *transport.Client) *retireorder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := retireorder.NewService(c)
	return s
}

func code(cis string) positions.TrackingCode {
	return positions.TrackingCode{CIS: cis, Type: "trackingcode"}
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	product      = ref("entity/product/77777777-7777-7777-7777-777777777777", "product")
	minimal      = retireorder.CreateRequest{Organization: organization, RetireOrderType: "RETAIL_SALE", TrackingType: "SHOES"}
	fullCreate   = retireorder.CreateRequest{
		Organization: organization, RetireOrderType: "EXPORT_INSIDE_EEU", TrackingType: "SHOES",
		Name: model.Value("00002"), Description: model.Value(""), ExternalCode: model.Value("ext"),
		SyncID: model.Value("22222222-2222-2222-2222-222222222222"),
		Moment: model.Value(documenttest.DateTime("2026-03-25 19:03:00")), Shared: model.Value(true),
		OrganizationBranch: model.Null[model.Reference](),
		Agent:              model.Value(ref("entity/counterparty/33333333-3333-3333-3333-333333333333", "counterparty")),
		DestinationCountry: model.Value(ref("entity/country/44444444-4444-4444-4444-444444444444", "country")),
		Owner:              model.Value(ref("entity/employee/55555555-5555-5555-5555-555555555555", "employee")),
		Group:              model.Value(ref("entity/group/66666666-6666-6666-6666-666666666666", "group")),
		State:              model.Null[model.Reference](),
		Rate:               model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/99999999-9999-9999-9999-999999999999", "currency"))}),
		Sum:                model.Value(0.0), VATEnabled: model.Value(false), VATIncluded: model.Value(true),
		ReasonDescription: model.Value("причина"), PrimaryDocumentName: model.Value("doc name"), StateContractID: model.Value("contract"),
		SupportingTransaction:       model.Value("OTHER"),
		SupportingTransactionDate:   model.Value(documenttest.DateTime("2026-03-24 00:00:00")),
		SupportingTransactionNumber: model.Value("123"),
		Positions: model.Value([]retireorder.PositionRequest{{
			ItemRequest: positions.ItemRequest{Assortment: model.Value(product), Quantity: model.Value(1.0)},
			Price:       model.Value(50000.5), VAT: model.Value(0), VATEnabled: model.Value(false),
			TrackingCodes: model.Value([]positions.TrackingCode{code("0114650057073880217891234567893")}),
		}}),
	}
	fullUpdate = retireorder.UpdateRequest{
		RetireOrderType: model.Value("OTHER_TYPE"),
		Name:            model.Value("обновленный"),
		Description:     model.Value("новое описание Вывода из оборота"),
		Agent:           model.Null[model.Reference](),
		Positions: model.Value([]retireorder.PositionRequest{
			{ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value(positionID)}},
				TrackingCodes: model.Value([]positions.TrackingCode{code("010465011724001921yM1QsEStTWZfG2406402"), code("110465011724001921yM1QsEStTWZfG2406402")})},
			{ItemRequest: positions.ItemRequest{Assortment: model.Value(product)}, Price: model.Value(123456789.01), TrackingCodes: model.Value([]positions.TrackingCode{})},
		}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/retireorder/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got retireorder.RetireOrder
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "dbb8cfc1-cbfa-11e1-6dfb-889ffa6f49fd" || got.Name != "00001" || got.Description != "some_some" ||
				got.ExternalCode != "4fXzFRjOjDC8upZ1iEK-n3" || got.Meta.Type != "retireorder" || got.Shared || got.Printed || got.Published {
				t.Fatalf("document = %+v", got)
			}
			if got.DocumentState != "CREATED" || got.RetireOrderType != "RETAIL_SALE" || got.TrackingType != "SHOES" || got.SupportingTransaction != "OTHER" ||
				got.SupportingTransactionNumber != "123" || got.PrimaryDocumentName != "doc name" {
				t.Fatalf("enums = %+v", got)
			}
			if got.Sum != 100000 || got.VATSum != 16667 || !got.VATEnabled || !got.VATIncluded {
				t.Fatalf("money = %v %v %v %v", got.Sum, got.VATSum, got.VATEnabled, got.VATIncluded)
			}
			if !time.Time(got.SupportingTransactionDate).Equal(time.Time(documenttest.DateTime("2024-11-22 13:37:51.335"))) ||
				!time.Time(got.Moment).Equal(time.Time(documenttest.DateTime("2024-11-26 13:20:51.272"))) || got.Created.IsZero() || got.Updated.IsZero() || !got.Deleted.IsZero() {
				t.Fatalf("dates = %+v", got)
			}
			if got.Organization == nil || got.OrganizationBranch == nil || got.OrganizationBranch.Meta.Type != "organizationbranch" || got.Owner == nil || got.Group == nil ||
				got.Rate == nil || got.Rate.Currency == nil || got.Rate.Value != nil || got.Agent != nil || got.State != nil {
				t.Fatalf("references = %+v", got)
			}
			if got.Positions == nil || got.Positions.Meta.Type != "retireorderposition" || got.Positions.Meta.Size != 1 || got.Positions.Rows != nil {
				t.Fatalf("positions = %+v", got.Positions)
			}
			var sparse retireorder.RetireOrder
			if err := codec.Unmarshal([]byte(`{"id":"x","syncId":"s","deleted":"2024-11-22 13:37:51.355","stateContractId":"c","reasonDescription":"r",`+
				`"agent":{"meta":{"type":"counterparty"}},"destinationCountry":{"meta":{"type":"country"}},"state":{"meta":{"type":"state"}},`+
				`"retireOrderType":"FUTURE_REASON","documentState":"FUTURE_STATE","trackingType":"FUTURE","future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			if sparse.SyncID != "s" || sparse.Deleted.IsZero() || sparse.StateContractID != "c" || sparse.ReasonDescription != "r" ||
				sparse.Agent == nil || sparse.DestinationCountry == nil || sparse.State == nil || sparse.RetireOrderType != "FUTURE_REASON" ||
				sparse.DocumentState != "FUTURE_STATE" || sparse.TrackingType != "FUTURE" || sparse.Organization != nil || sparse.Rate != nil {
				t.Fatalf("sparse = %+v", sparse)
			}
		})
	}
}

func TestPositionsFixture(t *testing.T) {
	data := read(t, "custom/retireorder/positions.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got model.Collection[retireorder.Position]
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.Meta.Size != 2 || len(got.Rows) != 2 {
				t.Fatalf("collection = %+v", got)
			}
			row := got.Rows[0]
			if row.ID != positionID || row.Meta.Type != "retireorderposition" || row.AccountID != "dbb8cfc1-cbfa-11e1-6dfb-889ffa6f49fd" || row.Quantity != 2 ||
				row.Price != 50000 || row.VAT != 20 || !row.VATEnabled || row.Assortment.Meta.Type != "product" {
				t.Fatalf("row = %+v", row)
			}
			want := []positions.TrackingCode{code("010463003759026521uHpIIf-111122"), code("010465011724001921yM1QsEStTWZfG2406402")}
			if !reflect.DeepEqual(row.TrackingCodes, want) || got.Rows[1].TrackingCodes[1].CIS != "010463003759026521uHpIIf-nXIH>1" {
				t.Fatalf("trackingCodes = %+v", got.Rows)
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
			{"create full", fullCreate, read(t, "custom/retireorder/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
				`"retireOrderType":"RETAIL_SALE","trackingType":"SHOES"}`)},
			{"update", fullUpdate, read(t, "custom/retireorder/update.golden.json")},
			{"update empty positions", retireorder.UpdateRequest{Positions: model.Value([]retireorder.PositionRequest{})}, []byte(`{"positions":[]}`)},
			{"position update", retireorder.PositionRequest{VAT: model.Value(0), VATEnabled: model.Value(false)}, []byte(`{"vat":0,"vatEnabled":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[retireorder.CreateRequest](), reflect.TypeFor[retireorder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "DocumentState", "VATSum"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"Meta", "AccountID"} {
		if _, ok := reflect.TypeFor[retireorder.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains read-only %s", name)
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":4040,"documentState":"SEND","positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*retireorder.RetireOrder, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*retireorder.RetireOrder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/retireorder", read(t, "custom/retireorder/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*retireorder.RetireOrder, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/retireorder/update.golden.json"), func(ctx context.Context, c *transport.Client) (*retireorder.RetireOrder, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *retireorder.RetireOrder) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 4040 || got.DocumentState != "SEND" || got.Positions.Rows != nil {
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
	const row = `{"meta":{"type":"retireorderposition"},"id":"` + positionID + `","quantity":1,"price":50000.5,"vat":20,"vatEnabled":true,` +
		`"assortment":{"meta":{"type":"product"}},"trackingCodes":[{"cis":"a","type":"trackingcode"}]}`
	checkRow := func(got *retireorder.Position) {
		if got.ID != positionID || got.Quantity != 1 || got.Price != 50000.5 || got.VAT != 20 || !got.VATEnabled ||
			!reflect.DeepEqual(got.TrackingCodes, []positions.TrackingCode{code("a")}) {
			t.Fatalf("row = %+v", got)
		}
	}
	t.Run("get positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=2&offset=7", "", `{"meta":{"href":"h"},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[retireorder.Position], error) {
				return service(c).GetPositions(ctx, testID, positions.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[retireorder.Position]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions/"+positionID, "", row, func(ctx context.Context, c *transport.Client) (*retireorder.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}, checkRow)
	})
	t.Run("update position", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/positions/"+positionID, `{"trackingCodes":[{"cis":"a","type":"trackingcode"}]}`, row,
			func(ctx context.Context, c *transport.Client) (*retireorder.Position, error) {
				return service(c).UpdatePosition(ctx, testID, positionID, retireorder.PositionRequest{TrackingCodes: model.Value([]positions.TrackingCode{code("a")})})
			}, checkRow)
	})
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/retireorder"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"retireorder","mediaType":"application/json"}`
	productJSON := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/77777777-7777-7777-7777-777777777777","type":"product"}}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"create many", "POST", "entity/retireorder", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
			`"retireOrderType":"RETAIL_SALE","trackingType":"SHOES"}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []retireorder.CreateRequest{minimal})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/retireorder", `[{"meta":` + docMeta + `,"stateContractId":"test"}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).UpdateMany(ctx, []batch.Update[retireorder.UpdateRequest]{{ID: testID, Request: retireorder.UpdateRequest{StateContractID: model.Value("test")}}})
				return err
			}},
		{"create positions", "POST", docPath + "/positions", `[{"assortment":` + productJSON + `,"trackingCodes":[{"cis":"c","type":"trackingcode"}]}]`, `[{"id":"` + positionID + `","price":50000}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreatePositions(ctx, testID, []retireorder.PositionRequest{{ItemRequest: positions.ItemRequest{Assortment: model.Value(product)},
					TrackingCodes: model.Value([]positions.TrackingCode{code("c")})}})
				if err == nil && (r[0].Value.ID != positionID || r[0].Value.Price != 50000) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete position", "DELETE", docPath + "/positions/" + positionID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, positionID)
		}},
		{"delete many positions", "POST", docPath + "/positions/delete", `[{"meta":{"href":"` + base + `/` + testID + `/positions/` + positionID + `","type":"retireorderposition","mediaType":"application/json"}}]`, "", true,
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
				server = memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retireorder") ||
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

// Источник не описывает удаление документа: Delete и DeleteMany нет, есть только Trash.
func TestNoDelete(t *testing.T) {
	typ := reflect.TypeFor[*retireorder.Service]()
	for _, name := range []string{"Delete", "DeleteMany"} {
		if _, ok := typ.MethodByName(name); ok {
			t.Errorf("Service has %s", name)
		}
	}
}

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 1000 строк уходят.
func TestValidation(t *testing.T) {
	// Корректная новая строка: ошибку в случаях ниже вызывает только изменённая часть.
	validRow := retireorder.PositionRequest{ItemRequest: positions.ItemRequest{Assortment: model.Value(product)},
		TrackingCodes: model.Value([]positions.TrackingCode{{CIS: "c", Type: "trackingcode"}})}
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		io.WriteString(w, `{"id":"server-id"}`)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	rows := func(n int) model.Field[[]retireorder.PositionRequest] {
		return model.Value(make([]retireorder.PositionRequest, n))
	}
	limit := minimal
	limit.Positions = rows(1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, retireorder.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	// Открытый enum причин вывода и documentState-правила Update — на сервере.
	if _, err := s.Create(ctx, retireorder.CreateRequest{Organization: organization, RetireOrderType: "FUTURE_REASON", TrackingType: "FUTURE"}); err != nil {
		t.Fatal(err)
	}
	create := func(change func(*retireorder.CreateRequest)) func(...transport.RequestOption) error {
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
		"get bad id":                func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":                  func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no org":             create(func(r *retireorder.CreateRequest) { r.Organization = model.Reference{} }),
		"create no retireOrderType": create(func(r *retireorder.CreateRequest) { r.RetireOrderType = "" }),
		"create no trackingType":    create(func(r *retireorder.CreateRequest) { r.TrackingType = "" }),
		"create 1001":               create(func(r *retireorder.CreateRequest) { r.Positions = rows(1001) }),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, retireorder.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update 1001": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, retireorder.UpdateRequest{Positions: rows(1001)}, o...)
			return err
		},
		"trash bad id": func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[retireorder.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[retireorder.UpdateRequest]{{ID: testID}}, o...)
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
			_, err := s.CreatePositions(ctx, "bad", []retireorder.PositionRequest{validRow}, o...)
			return err
		},
		"create positions no assortment": func(o ...transport.RequestOption) error {
			row := validRow
			row.Assortment = model.Field[model.Reference]{}
			_, err := s.CreatePositions(ctx, testID, []retireorder.PositionRequest{validRow, row}, o...)
			return err
		},
		"create positions no tracking codes": func(o ...transport.RequestOption) error {
			row := validRow
			row.TrackingCodes = model.Value([]positions.TrackingCode{})
			_, err := s.CreatePositions(ctx, testID, []retireorder.PositionRequest{row}, o...)
			return err
		},
		"update position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, retireorder.PositionRequest{}, o...)
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
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retireorder") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := retireorder.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "retireorder", read(t, "custom/retireorder/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[retireorder.RetireOrder] {
		return service(c).List
	}, func(t *testing.T, got retireorder.RetireOrder) {
		if got.ID != fixtureID || got.RetireOrderType != "RETAIL_SALE" || got.Sum != 100000 {
			t.Fatalf("row = %+v", got)
		}
	})
}
