package commissionreportin_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/commissionreportin"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	fixtureID  = "cc5beac8-acdb-11e6-8a84-bae500000000"
	positionID = "34f6344f-015e-11e6-9464-e4de0000006c"
	docPath    = "entity/commissionreportin/" + testID
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

func service(c *transport.Client) *commissionreportin.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := commissionreportin.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	agent        = ref("entity/counterparty/22222222-2222-2222-2222-222222222222", "counterparty")
	contract     = ref("entity/contract/33333333-3333-3333-3333-333333333333", "contract")
	periodStart  = documenttest.DateTime("2016-11-01 10:07:00")
	periodEnd    = documenttest.DateTime("2016-11-29 10:07:00")
	minimal      = commissionreportin.CreateRequest{Organization: organization, Agent: agent, Contract: contract, CommissionPeriodStart: periodStart, CommissionPeriodEnd: periodEnd}
	fullCreate   = commissionreportin.CreateRequest{
		Organization: organization, Agent: agent, Contract: contract, CommissionPeriodStart: periodStart, CommissionPeriodEnd: periodEnd,
		Name: model.Value("3335551"), Description: model.Value(""), Code: model.Value("c1"), ExternalCode: model.Value("extCod1"),
		Moment: model.Value(documenttest.DateTime("2016-12-16 11:22:33")), Applicable: model.Value(true), Shared: model.Value(false),
		OrganizationAccount: model.Value(ref("entity/organization/11111111-1111-1111-1111-111111111111/accounts/44444444-4444-4444-4444-444444444444", "account")),
		AgentAccount:        model.Null[model.Reference](),
		Owner:               model.Value(ref("entity/employee/55555555-5555-5555-5555-555555555555", "employee")),
		Group:               model.Value(ref("entity/group/66666666-6666-6666-6666-666666666666", "group")),
		Project:             model.Value(ref("entity/project/77777777-7777-7777-7777-777777777777", "project")),
		SalesChannel:        model.Value(ref("entity/saleschannel/88888888-8888-8888-8888-888888888888", "saleschannel")),
		State:               model.Null[model.Reference](),
		Attributes:          model.Value([]model.Attribute{}),
		Rate:                model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/99999999-9999-9999-9999-999999999999", "currency"))}),
		VATEnabled:          model.Value(true), VATIncluded: model.Value(false),
		RewardType: model.Value("PercentOfSales"), RewardPercent: model.Value(15),
		CommissionOverhead: model.Value(commissionreportin.CommissionOverhead{Sum: -500}),
		Positions: model.Value([]commissionreportin.PositionRequest{{
			Assortment: model.Value(ref("entity/variant/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "variant")),
			Quantity:   model.Value(2.0), Price: model.Value(13205.0), VAT: model.Value(10), VATEnabled: model.Value(true),
			Reward: model.Value(-123.0), Pack: model.Value(positions.PackRequest{ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"}),
		}}),
		ReturnToCommissionerPositions: model.Value([]commissionreportin.ReturnPositionRequest{{
			Assortment: model.Value(ref("entity/product/cccccccc-cccc-cccc-cccc-cccccccccccc", "product")),
			Quantity:   model.Value(1.0), Price: model.Value(9939.9), VAT: model.Value(0), Reward: model.Value(100.0),
		}}),
		Payments: model.Value([]model.Reference{ref("entity/paymentin/dddddddd-dddd-dddd-dddd-dddddddddddd", "paymentin")}),
	}
	fullUpdate = commissionreportin.UpdateRequest{
		Contract:            model.Value(ref("entity/contract/eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee", "contract")),
		CommissionPeriodEnd: model.Value(documenttest.DateTime("2016-12-29 10:07:00")),
		Project:             model.Null[model.Reference](),
		RewardType:          model.Value("None"), RewardPercent: model.Value(0),
		CommissionOverhead:            model.Value(commissionreportin.CommissionOverhead{}),
		Positions:                     model.Value([]commissionreportin.PositionRequest{}),
		ReturnToCommissionerPositions: model.Value([]commissionreportin.ReturnPositionRequest{{Quantity: model.Value(3.0)}}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/commissionreportin/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got commissionreportin.CommissionReportIn
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Sum != 1764 || got.VATSum != 24 || got.CommitentSum != 1499 || got.RewardType != "PercentOfSales" || got.RewardPercent != 15 ||
				!got.VATEnabled || !got.VATIncluded || !got.Printed || !got.Published || !got.Applicable || got.Shared || got.Name != "3335551" {
				t.Fatalf("document = %+v", got)
			}
			if !time.Time(got.CommissionPeriodStart).Equal(time.Time(periodStart)) || !time.Time(got.CommissionPeriodEnd).Equal(time.Time(periodEnd)) || got.Created.IsZero() || got.Updated.IsZero() || got.Moment.IsZero() {
				t.Fatalf("dates = %+v", got)
			}
			if got.Organization == nil || got.OrganizationAccount == nil || got.Agent == nil || got.Contract == nil || got.Owner == nil || got.Group == nil || got.Project == nil ||
				got.Rate == nil || got.Rate.Currency == nil || len(got.Attributes) != 2 {
				t.Fatalf("references = %+v", got)
			}
			if got.Positions == nil || got.Positions.Meta.Type != "commissionreportinposition" || got.Positions.Meta.Size != 2 ||
				got.ReturnToCommissionerPositions == nil || got.ReturnToCommissionerPositions.Meta.Type != "commissionreportinreturnedposition" {
				t.Fatalf("collections = %+v / %+v", got.Positions, got.ReturnToCommissionerPositions)
			}
			// Поля таблицы, которых нет в примере ответа, и строки при expand.
			var sparse commissionreportin.CommissionReportIn
			if err := codec.Unmarshal([]byte(`{"id":"x","commissionOverhead":{"sum":-150},"payedSum":100,"syncId":"s","salesChannel":{"meta":{"type":"saleschannel"}},`+
				`"payments":[{"meta":{"type":"paymentin"},"linkedSum":0}],"positions":{"meta":{"href":"h"},"rows":[{"id":"p","quantity":2,"price":13205,"vat":10,"vatEnabled":true,"reward":-123,"pack":{"id":"k"}}]},`+
				`"returnToCommissionerPositions":{"meta":{"href":"r"},"rows":[{"id":"r","quantity":1,"reward":100}]},"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			row, back := sparse.Positions.Rows[0], sparse.ReturnToCommissionerPositions.Rows[0]
			if sparse.CommissionOverhead == nil || sparse.CommissionOverhead.Sum != -150 || sparse.PayedSum != 100 || sparse.SyncID != "s" || sparse.SalesChannel == nil ||
				len(sparse.Payments) != 1 || row.Reward != -123 || row.VAT != 10 || !row.VATEnabled || row.Pack == nil || row.Pack.ID != "k" || back.Reward != 100 || back.Quantity != 1 ||
				sparse.Organization != nil || !sparse.Created.IsZero() || !sparse.Deleted.IsZero() {
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
			{"create full", fullCreate, read(t, "custom/commissionreportin/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
				`"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/22222222-2222-2222-2222-222222222222","type":"counterparty"}},` +
				`"contract":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/contract/33333333-3333-3333-3333-333333333333","type":"contract"}},` +
				`"commissionPeriodStart":"2016-11-01 10:07:00.000","commissionPeriodEnd":"2016-11-29 10:07:00.000"}`)},
			{"update", fullUpdate, read(t, "custom/commissionreportin/update.golden.json")},
			{"update existing rows", commissionreportin.UpdateRequest{Positions: model.Value([]commissionreportin.PositionRequest{{ID: model.Value("7944ef04-f831-11e5-7a69-971500188b20"), Price: model.Value(0.5)}}), ReturnToCommissionerPositions: model.Value([]commissionreportin.ReturnPositionRequest{{ID: model.Value("7944ef04-f831-11e5-7a69-971500188b20"), Reward: model.Value(0.0)}})}, []byte(`{"positions":[{"id":"7944ef04-f831-11e5-7a69-971500188b20","price":0.5}],"returnToCommissionerPositions":[{"id":"7944ef04-f831-11e5-7a69-971500188b20","reward":0}]}`)},
			{"position update", commissionreportin.PositionRequest{Price: model.Value(0.0), VATEnabled: model.Value(false)}, []byte(`{"price":0,"vatEnabled":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[commissionreportin.CreateRequest](), reflect.TypeFor[commissionreportin.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum", "PayedSum", "CommitentSum", "Files", "SyncID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[commissionreportin.PositionRequest](), reflect.TypeFor[commissionreportin.ReturnPositionRequest]()} {
		for _, name := range []string{"Meta", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[commissionreportin.ReturnPositionRequest]().FieldByName("Pack"); ok {
		t.Error("return position has no pack in the source table")
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и суммы — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":24691,"rewardPercent":7,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*commissionreportin.CommissionReportIn, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*commissionreportin.CommissionReportIn, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/commissionreportin", read(t, "custom/commissionreportin/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*commissionreportin.CommissionReportIn, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/commissionreportin/update.golden.json"), func(ctx context.Context, c *transport.Client) (*commissionreportin.CommissionReportIn, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *commissionreportin.CommissionReportIn) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 24691 || got.RewardPercent != 7 || got.Positions.Rows != nil {
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
	const (
		row      = `{"meta":{"type":"commissionreportinposition"},"id":"` + positionID + `","quantity":14,"price":1301,"vat":10,"vatEnabled":true,"reward":-5,"assortment":{"meta":{"type":"variant"}}}`
		returned = `{"meta":{"type":"commissionreportinreturnedposition"},"id":"` + positionID + `","quantity":1,"price":99,"vat":0,"vatEnabled":false,"reward":100}`
		update   = `{"quantity":14,"reward":0}`
	)
	request := commissionreportin.PositionRequest{Quantity: model.Value(14.0), Reward: model.Value(0.0)}
	returnRequest := commissionreportin.ReturnPositionRequest{Quantity: model.Value(14.0), Reward: model.Value(0.0)}
	checkRow := func(got *commissionreportin.Position) {
		if got.ID != positionID || got.Reward != -5 || got.Quantity != 14 || got.Assortment == nil {
			t.Fatalf("row = %+v", got)
		}
	}
	checkReturned := func(got *commissionreportin.ReturnPosition) {
		if got.ID != positionID || got.Reward != 100 || got.Price != 99 {
			t.Fatalf("row = %+v", got)
		}
	}
	t.Run("get positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=2&offset=7", "", `{"meta":{"href":"h"},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[commissionreportin.Position], error) {
				return service(c).GetPositions(ctx, testID, positions.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[commissionreportin.Position]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions/"+positionID, "", row, func(ctx context.Context, c *transport.Client) (*commissionreportin.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}, checkRow)
	})
	t.Run("update position", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/positions/"+positionID, update, row, func(ctx context.Context, c *transport.Client) (*commissionreportin.Position, error) {
			return service(c).UpdatePosition(ctx, testID, positionID, request)
		}, checkRow)
	})
	t.Run("get return positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/returntocommissionerpositions?limit=1000&offset=0", "", `{"meta":{"href":"h"},"rows":[`+returned+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[commissionreportin.ReturnPosition], error) {
				return service(c).GetReturnPositions(ctx, testID, positions.PageOptions{})
			}, func(got *model.Collection[commissionreportin.ReturnPosition]) { checkReturned(&got.Rows[0]) })
	})
	t.Run("get return position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/returntocommissionerpositions/"+positionID, "", returned, func(ctx context.Context, c *transport.Client) (*commissionreportin.ReturnPosition, error) {
			return service(c).GetReturnPosition(ctx, testID, positionID)
		}, checkReturned)
	})
	t.Run("update return position", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/returntocommissionerpositions/"+positionID, update, returned, func(ctx context.Context, c *transport.Client) (*commissionreportin.ReturnPosition, error) {
			return service(c).UpdateReturnPosition(ctx, testID, positionID, returnRequest)
		}, checkReturned)
	})
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/commissionreportin"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"commissionreportin","mediaType":"application/json"}`
	assortment := ref("entity/product/cccccccc-cccc-cccc-cccc-cccccccccccc", "product")
	assortmentJSON := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/cccccccc-cccc-cccc-cccc-cccccccccccc","type":"product"}}`
	results := func(err error) error { return err }
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"create many", "POST", "entity/commissionreportin", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
			`"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/22222222-2222-2222-2222-222222222222","type":"counterparty"}},` +
			`"contract":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/contract/33333333-3333-3333-3333-333333333333","type":"contract"}},` +
			`"commissionPeriodStart":"2016-11-01 10:07:00.000","commissionPeriodEnd":"2016-11-29 10:07:00.000"}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []commissionreportin.CreateRequest{minimal})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return results(err)
			}},
		{"update many", "POST", "entity/commissionreportin", `[{"meta":` + docMeta + `,"rewardPercent":0}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).UpdateMany(ctx, []batch.Update[commissionreportin.UpdateRequest]{{ID: testID, Request: commissionreportin.UpdateRequest{RewardPercent: model.Value(0)}}})
				return err
			}},
		{"delete many", "POST", "entity/commissionreportin/delete", `[{"meta":` + docMeta + `}]`, `[{"info":"ok"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).DeleteMany(ctx, []string{testID})
				return err
			}},
		{"create positions", "POST", docPath + "/positions", `[{"assortment":` + assortmentJSON + `,"quantity":2,"reward":12}]`, `[{"id":"` + positionID + `"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreatePositions(ctx, testID, []commissionreportin.PositionRequest{{Assortment: model.Value(assortment), Quantity: model.Value(2.0), Reward: model.Value(12.0)}})
				if err == nil && r[0].Value.ID != positionID {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete position", "DELETE", docPath + "/positions/" + positionID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, positionID)
		}},
		{"delete many positions", "POST", docPath + "/positions/delete", `[{"meta":{"href":"` + base + `/` + testID + `/positions/` + positionID + `","type":"commissionreportinposition","mediaType":"application/json"}}]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyPositions(ctx, testID, []string{positionID})
			}},
		{"create return positions", "POST", docPath + "/returntocommissionerpositions", `[{"assortment":` + assortmentJSON + `,"quantity":1}]`, `[{"id":"` + positionID + `"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateReturnPositions(ctx, testID, []commissionreportin.ReturnPositionRequest{{Assortment: model.Value(assortment), Quantity: model.Value(1.0)}})
				if err == nil && r[0].Value.ID != positionID {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete return position", "DELETE", docPath + "/returntocommissionerpositions/" + positionID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteReturnPosition(ctx, testID, positionID)
		}},
		{"delete many return positions", "POST", docPath + "/returntocommissionerpositions/delete", `[{"meta":{"href":"` + base + `/` + testID + `/returntocommissionerpositions/` + positionID + `","type":"commissionreportinreturnedposition","mediaType":"application/json"}}]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyReturnPositions(ctx, testID, []string{positionID})
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "commissionreportin") ||
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

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 1000 строк каждого набора уходят.
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
	sold := func(n int) model.Field[[]commissionreportin.PositionRequest] {
		return model.Value(make([]commissionreportin.PositionRequest, n))
	}
	returned := func(n int) model.Field[[]commissionreportin.ReturnPositionRequest] {
		return model.Value(make([]commissionreportin.ReturnPositionRequest, n))
	}
	limit := minimal
	limit.Positions, limit.ReturnToCommissionerPositions = sold(1000), returned(1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, commissionreportin.UpdateRequest{Positions: sold(1000), ReturnToCommissionerPositions: returned(1000)}); err != nil {
		t.Fatal(err)
	}
	create := func(change func(*commissionreportin.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	page := func(p positions.PageOptions, id string, get func(context.Context, string, positions.PageOptions, ...transport.RequestOption) error) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error { return get(ctx, id, p, o...) }
	}
	getPositions := func(ctx context.Context, id string, p positions.PageOptions, o ...transport.RequestOption) error {
		_, err := s.GetPositions(ctx, id, p, o...)
		return err
	}
	getReturns := func(ctx context.Context, id string, p positions.PageOptions, o ...transport.RequestOption) error {
		_, err := s.GetReturnPositions(ctx, id, p, o...)
		return err
	}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":          func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":            func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no org":       create(func(r *commissionreportin.CreateRequest) { r.Organization = model.Reference{} }),
		"create no agent":     create(func(r *commissionreportin.CreateRequest) { r.Agent = model.Reference{} }),
		"create no contract":  create(func(r *commissionreportin.CreateRequest) { r.Contract = model.Reference{} }),
		"create no start":     create(func(r *commissionreportin.CreateRequest) { r.CommissionPeriodStart = model.DateTime{} }),
		"create no end":       create(func(r *commissionreportin.CreateRequest) { r.CommissionPeriodEnd = model.DateTime{} }),
		"create 1001":         create(func(r *commissionreportin.CreateRequest) { r.Positions = sold(1001) }),
		"create 1001 returns": create(func(r *commissionreportin.CreateRequest) { r.ReturnToCommissionerPositions = returned(1001) }),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, commissionreportin.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update 1001 returns": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, commissionreportin.UpdateRequest{ReturnToCommissionerPositions: returned(1001)}, o...)
			return err
		},
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":  func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[commissionreportin.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[commissionreportin.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"positions bad id":       page(positions.PageOptions{}, "bad", getPositions),
		"positions fields":       page(positions.PageOptions{Fields: "stock"}, testID, getPositions),
		"positions limit":        page(positions.PageOptions{Limit: 1001}, testID, getPositions),
		"return positions bad":   page(positions.PageOptions{}, "bad", getReturns),
		"return positions field": page(positions.PageOptions{Fields: "stock"}, testID, getReturns),
		"return positions limit": page(positions.PageOptions{Offset: -1}, testID, getReturns),
		"get position bad": func(o ...transport.RequestOption) error {
			_, err := s.GetPosition(ctx, testID, "bad", o...)
			return err
		},
		"create positions empty": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, testID, nil, o...)
			return err
		},
		"create positions bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, "bad", []commissionreportin.PositionRequest{{}}, o...)
			return err
		},
		"update position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, commissionreportin.PositionRequest{}, o...)
			return err
		},
		"delete position bad": func(o ...transport.RequestOption) error { return s.DeletePosition(ctx, "bad", positionID, o...) },
		"delete many positions empty": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, nil, o...)
		},
		"get return position bad": func(o ...transport.RequestOption) error {
			_, err := s.GetReturnPosition(ctx, "bad", positionID, o...)
			return err
		},
		"create return positions empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateReturnPositions(ctx, testID, nil, o...)
			return err
		},
		"update return position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateReturnPosition(ctx, testID, positionID, commissionreportin.ReturnPositionRequest{}, o...)
			return err
		},
		"delete return position bad": func(o ...transport.RequestOption) error {
			return s.DeleteReturnPosition(ctx, testID, "bad", o...)
		},
		"delete many returns empty": func(o ...transport.RequestOption) error {
			return s.DeleteManyReturnPositions(ctx, testID, nil, o...)
		},
		"delete many returns 1001": func(o ...transport.RequestOption) error {
			return s.DeleteManyReturnPositions(ctx, testID, make([]string, 1001), o...)
		},
		"delete many returns bad doc": func(o ...transport.RequestOption) error {
			return s.DeleteManyReturnPositions(ctx, "bad", []string{positionID}, o...)
		},
		"delete many returns bad row": func(o ...transport.RequestOption) error {
			return s.DeleteManyReturnPositions(ctx, testID, []string{positionID, "bad"}, o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "commissionreportin") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := commissionreportin.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "commissionreportin", read(t, "custom/commissionreportin/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[commissionreportin.CommissionReportIn] {
		return service(c).List
	}, func(t *testing.T, got commissionreportin.CommissionReportIn) {
		if got.ID != fixtureID || got.ReturnToCommissionerPositions == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
