package prepayment_test

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

	"github.com/eteeet/moysklad-sdk-go/entity/prepayment"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	positionID = "34f6344f-015e-11e6-9464-e4de0000006c"
	fixtureID  = "7944ef04-f831-11e5-7a69-971500188b19"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *prepayment.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := prepayment.NewService(c)
	return s
}

func TestPrepaymentFixture(t *testing.T) {
	data, row := read(t, "custom/prepayment/read.json"), read(t, "custom/prepayment/position.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got prepayment.Prepayment
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Sum != 30000 || got.CashSum != 15000 || got.NoCashSum != 15000 || got.QRSum != 0 || got.VATSum != 5000 ||
				got.RetailShift == nil || got.RetailShift.Meta.Type != "retailshift" || got.RetailStore == nil || got.CustomerOrder == nil ||
				got.Agent == nil || got.OrganizationBranch == nil || !got.VATEnabled || !got.VATIncluded || got.Positions == nil ||
				got.Positions.Meta.Type != "prepaymentposition" || got.Positions.Rows != nil {
				t.Fatalf("document = %+v", got)
			}
			var position prepayment.Position
			if err := codec.Unmarshal(row, &position); err != nil {
				t.Fatal(err)
			}
			if position.ID != positionID || position.Meta.Type != "prepaymentposition" || position.Quantity != 2 || position.Price != 15000 ||
				position.VAT != 20 || !position.VATEnabled || position.Assortment.Meta.Type != "product" {
				t.Fatalf("position = %+v", position)
			}
			var sparse prepayment.Prepayment
			if err := codec.Unmarshal([]byte(`{"id":"x","returns":[{"meta":{"href":"h","type":"prepaymentreturn"}}],"taxSystem":"PATENT_BASED","future":1}`), &sparse); err != nil ||
				len(sparse.Returns) != 1 || sparse.TaxSystem != "PATENT_BASED" || sparse.RetailShift != nil || sparse.Positions != nil || !sparse.Created.IsZero() {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

// Записи нет: сервис не предлагает Create, Update и массовые операции.
func TestPrepaymentReadOnly(t *testing.T) {
	typ := reflect.TypeFor[*prepayment.Service]()
	for _, name := range []string{"Create", "Update", "CreateMany", "UpdateMany", "DeleteMany", "New", "CreatePositions", "UpdatePosition", "DeletePosition", "DeleteManyPositions"} {
		if _, ok := typ.MethodByName(name); ok {
			t.Errorf("Service has %s", name)
		}
	}
}

func TestPrepaymentHTTP(t *testing.T) {
	const document = `{"id":"server-id","name":"server-name","sum":30000,"cashSum":10000,"noCashSum":20000,"qrSum":0,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	invoke := func(ctx context.Context, c *transport.Client) (*prepayment.Prepayment, error) {
		return service(c).Get(ctx, testID)
	}
	for _, tc := range documenttest.ResponseCases(document) {
		t.Run("get/"+tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/prepayment/"+testID, nil, invoke, func(got *prepayment.Prepayment) {
				if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 30000 || got.CashSum != 10000 || got.NoCashSum != 20000 || got.Positions.Rows != nil {
					t.Fatalf("document = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)

	position := func(ctx context.Context, c *transport.Client) (*prepayment.Position, error) {
		return service(c).GetPosition(ctx, testID, positionID)
	}
	for _, tc := range documenttest.ResponseCases(string(read(t, "custom/prepayment/position.json"))) {
		t.Run("get position/"+tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/prepayment/"+testID+"/positions/"+positionID, nil, position, func(got *prepayment.Position) {
				if got.ID != positionID || got.Price != 15000 {
					t.Fatalf("position = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, position)
}

func TestPrepaymentGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7}} {
		query := "limit=1000&offset=0"
		if page.Limit == 1 {
			query = "limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[prepayment.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		body := `{"meta":{"href":"https://example.test/positions","type":"prepaymentposition"},"rows":[` + string(read(t, "custom/prepayment/position.json")) + `]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/prepayment/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[prepayment.Position]) {
					if len(got.Rows) != 1 || got.Rows[0].ID != positionID || got.Rows[0].Discount != 0 || got.Rows[0].VAT != 20 {
						t.Fatalf("rows = %+v", got.Rows)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// exchange — один вызов против сервера со статусом status: запрос, число обращений, ошибка.
func exchange(t *testing.T, status int, response string, invoke func(context.Context, *transport.Client) error) (string, int32, error) {
	t.Helper()
	var calls atomic.Int32
	var request string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		data, _ := io.ReadAll(r.Body)
		request = r.Method + " " + r.URL.RequestURI() + " " + string(data)
		w.WriteHeader(status)
		io.WriteString(w, response)
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
	if err != nil {
		t.Fatal(err)
	}
	err = invoke(t.Context(), c)
	return request, calls.Load(), err
}

// Delete и Trash: маршрут без тела, 200 и 204 — успех, ошибка и неожиданный статус — без повтора.
func TestPrepaymentDeleteAndTrash(t *testing.T) {
	for _, op := range []struct {
		name, request string
		invoke        func(context.Context, *transport.Client) error
	}{
		{"delete prepayment", "DELETE /entity/prepayment/" + testID + " ", func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash prepayment", "POST /entity/prepayment/" + testID + "/trash ", func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
	} {
		for _, c := range []struct {
			status int
			body   string
			ok     bool
		}{{200, "", true}, {204, "", true}, {403, `{"errors":[{"error":"forbidden","code":1016}]}`, false}, {500, `{"errors":[{"error":"server","code":1000}]}`, false}, {202, "", false}} {
			request, calls, err := exchange(t, c.status, c.body, op.invoke)
			if request != op.request || calls != 1 || (err == nil) != c.ok {
				t.Fatalf("%s %d: request = %q, calls = %d, err = %v", op.name, c.status, request, calls, err)
			}
			if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), op.name)) {
				t.Errorf("%s %d: err = %v", op.name, c.status, err)
			}
		}
	}
}

// Неверные id и fields отклоняются до запроса; ResponseInfo сбрасывается.
func TestPrepaymentValidation(t *testing.T) {
	for name, invoke := range map[string]func(context.Context, *transport.Client, ...transport.RequestOption) error{
		"get": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Get(ctx, "bad", o...)
			return err
		},
		"delete": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).Delete(ctx, testID+"/x", o...)
		},
		"trash": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).Trash(ctx, "bad", o...)
		},
		"positions id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).GetPositions(ctx, "bad", positions.PageOptions{}, o...)
			return err
		},
		"positions fields": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).GetPositions(ctx, testID, positions.PageOptions{Fields: "stock"}, o...)
			return err
		},
		"positions limit": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).GetPositions(ctx, testID, positions.PageOptions{Limit: 1001}, o...)
			return err
		},
		"position id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).GetPosition(ctx, testID, "bad", o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		_, calls, err := exchange(t, 200, `{"id":"x"}`, func(ctx context.Context, c *transport.Client) error {
			return invoke(ctx, c, transport.WithResponseInfo(&info))
		})
		if err == nil || calls != 0 || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "prepayment") {
			t.Errorf("%s: calls = %d, info = %+v, err = %v", name, calls, info, err)
		}
	}
	if _, err := prepayment.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestPrepaymentList(t *testing.T) {
	listtest.Run(t, "prepayment", read(t, "custom/prepayment/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[prepayment.Prepayment] {
		return service(c).List
	}, func(t *testing.T, got prepayment.Prepayment) {
		if got.ID != fixtureID || got.RetailShift == nil || got.CashSum != 15000 {
			t.Fatalf("row = %+v", got)
		}
	})
}
