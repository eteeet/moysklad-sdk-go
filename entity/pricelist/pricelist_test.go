package pricelist_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/pricelist"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	positionID = "9560e3e3-9609-11e6-8af5-581e00000008"
	fixtureID  = "312cc207-afe9-11e6-8af5-581e0000002c"
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

func service(c *transport.Client) *pricelist.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := pricelist.NewService(c)
	return s
}

var (
	product    = ref("entity/product/88888888-8888-8888-8888-888888888888", "product")
	columns    = []pricelist.ColumnRequest{{Name: "without_discount"}}
	fullCreate = pricelist.CreateRequest{
		Columns: []pricelist.ColumnRequest{{Name: "without_discount"}, {Name: "with_discount", PercentageDiscount: model.Value(1000)},
			{Name: "zero", PercentageDiscount: model.Value(0)}},
		Name: model.Value("DOC-1"), Description: model.Value(""), Code: model.Value("P-1"), SyncID: model.Value("22222222-2222-2222-2222-222222222222"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Organization: model.Value(ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")),
		Attributes:   model.Value([]model.Attribute{}),
		PriceType:    model.Value(model.PriceType{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/672559f1-cbf3-11e1-9eb9-889ffa6f49fd", Type: "pricetype"}}),
		Positions: model.Value(pricelist.PositionRows{Rows: []pricelist.PositionRequest{{
			Assortment: model.Value(product),
			Cells:      model.Value([]pricelist.Cell{{Column: "without_discount", Sum: 90000}, {Column: "with_discount", Sum: 1000.5}}),
			Pack:       model.Value(positions.PackRequest{ID: "77777777-7777-7777-7777-777777777777"}),
		}}}),
	}
	fullUpdate = pricelist.UpdateRequest{
		Name: model.Value("Новое название"), Organization: model.Null[model.Reference](), Applicable: model.Value(true),
		Positions: model.Value(pricelist.PositionRows{Rows: []pricelist.PositionRequest{
			{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), Cells: model.Value([]pricelist.Cell{})},
			{Assortment: model.Value(product)},
		}}),
	}
)

func TestPriceListFixture(t *testing.T) {
	data, row := read(t, "custom/pricelist/read.json"), read(t, "custom/pricelist/position.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got pricelist.PriceList
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.PriceType == nil || got.PriceType.Name != "Цена продажи" ||
				len(got.Columns) != 2 || got.Columns[0] != (pricelist.Column{Name: "a"}) || got.Columns[1] != (pricelist.Column{Name: "b", PercentageDiscount: 1000}) ||
				got.Positions == nil || got.Positions.Meta.Type != "pricelistrow" || got.Positions.Meta.Size != 2 || got.Positions.Rows != nil {
				t.Fatalf("document = %+v", got)
			}
			var position pricelist.Position
			if err := codec.Unmarshal(row, &position); err != nil {
				t.Fatal(err)
			}
			if position.ID != positionID || position.Meta.Type != "pricelistrow" || position.Assortment == nil ||
				len(position.Cells) != 1 || position.Cells[0] != (pricelist.Cell{Column: "123", Sum: 100000}) || position.Pack != nil {
				t.Fatalf("position = %+v", position)
			}
			var sparse pricelist.PriceList
			if err := codec.Unmarshal([]byte(`{"id":"x","positions":{"meta":{"href":"h","size":2}},"future":1}`), &sparse); err != nil ||
				sparse.Positions.Rows != nil || sparse.PriceType != nil || sparse.Columns != nil || !sparse.Created.IsZero() {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestPriceListRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/pricelist/create-full.golden.json")},
			{"create minimal", pricelist.CreateRequest{Columns: columns}, []byte(`{"columns":[{"name":"without_discount"}]}`)},
			{"update", fullUpdate, read(t, "custom/pricelist/update.golden.json")},
			{"update empty positions", pricelist.UpdateRequest{Positions: model.Value(pricelist.PositionRows{Rows: []pricelist.PositionRequest{}})}, []byte(`{"positions":{"rows":[]}}`)},
			{"position", pricelist.PositionRequest{Assortment: model.Value(product), Cells: model.Value([]pricelist.Cell{{Column: "123", Sum: 100000}})},
				[]byte(`{"assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/88888888-8888-8888-8888-888888888888","type":"product"}},"cells":[{"column":"123","sum":100000}]}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[pricelist.CreateRequest](), reflect.TypeFor[pricelist.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Files", "Sum", "Store", "Rate"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
	// Столбцы после создания не меняются.
	if _, ok := reflect.TypeFor[pricelist.UpdateRequest]().FieldByName("Columns"); ok {
		t.Error("UpdateRequest contains immutable Columns")
	}
	// priceType (readonly) пишется только в Create — исключение по примеру Create, утверждено
	// пользователем; пример Update его не передаёт.
	if _, ok := reflect.TypeFor[pricelist.UpdateRequest]().FieldByName("PriceType"); ok {
		t.Error("UpdateRequest contains read-only PriceType")
	}
}

func TestPriceListHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и столбцы — только из ответа.
	const response = `{"id":"server-id","name":"server-name","columns":[{"name":"s"}],"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*pricelist.PriceList, error)
	}{
		{"get", "GET", "entity/pricelist/" + testID, nil, func(ctx context.Context, c *transport.Client) (*pricelist.PriceList, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/pricelist", read(t, "custom/pricelist/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*pricelist.PriceList, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/pricelist/" + testID, read(t, "custom/pricelist/update.golden.json"), func(ctx context.Context, c *transport.Client) (*pricelist.PriceList, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *pricelist.PriceList) {
						if got.ID != "server-id" || got.Name != "server-name" || len(got.Columns) != 1 || got.Columns[0].Name != "s" || got.Positions.Rows != nil {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestPriceListPositionHTTP(t *testing.T) {
	row := string(read(t, "custom/pricelist/position.json"))
	path := "entity/pricelist/" + testID + "/positions"
	for _, op := range []struct {
		name, method, path, body, response string
		invoke                             func(context.Context, *transport.Client) (*pricelist.Position, error)
	}{
		{"get", "GET", path + "/" + positionID, "", row, func(ctx context.Context, c *transport.Client) (*pricelist.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}},
		// Пример источника: PUT {"cells":[…]}; ответ — объект или массив из одного объекта.
		{"update", "PUT", path + "/" + positionID, `{"cells":[{"column":"123","sum":100002}]}`, row, func(ctx context.Context, c *transport.Client) (*pricelist.Position, error) {
			return service(c).UpdatePosition(ctx, testID, positionID, pricelist.PositionRequest{Cells: model.Value([]pricelist.Cell{{Column: "123", Sum: 100002}})})
		}},
		{"update array", "PUT", path + "/" + positionID, `{"cells":[]}`, "[" + row + "]", func(ctx context.Context, c *transport.Client) (*pricelist.Position, error) {
			return service(c).UpdatePosition(ctx, testID, positionID, pricelist.PositionRequest{Cells: model.Value([]pricelist.Cell{})})
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			var body []byte
			if op.body != "" {
				body = []byte(op.body)
			}
			for _, tc := range documenttest.ResponseCases(op.response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, body, op.invoke, func(got *pricelist.Position) {
						if got.ID != positionID || len(got.Cells) != 1 || got.Cells[0].Sum != 100000 {
							t.Fatalf("position = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7}} {
		query := "limit=1000&offset=0"
		if page.Limit == 1 {
			query = "limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[pricelist.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		for _, tc := range documenttest.ResponseCases(`{"meta":{"href":"https://example.test/positions","type":"pricelistrow"},"rows":[` + row + `]}`) {
			t.Run("list/"+tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", path+"?"+query, nil, invoke, func(got *model.Collection[pricelist.Position]) {
					if len(got.Rows) != 1 || got.Rows[0].Cells[0].Column != "123" {
						t.Fatalf("rows = %+v", got.Rows)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// exchange — один вызов против сервера со статусом status и телом response: запрос, тело (адрес
// сервера заменён на BASE), число обращений, ошибка.
func exchange(t *testing.T, status int, response string, invoke func(context.Context, *transport.Client) error) (string, string, int32, error) {
	t.Helper()
	var calls atomic.Int32
	var request, body string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

// Delete, Trash, массовые операции и операции строк: маршрут, тело, успех, ошибка без повтора.
func TestPriceListWrites(t *testing.T) {
	meta := func(id string) string {
		return `{"meta":{"href":"BASE/entity/pricelist/` + id + `","metadataHref":"BASE/entity/pricelist/metadata","type":"pricelist","mediaType":"application/json"}`
	}
	row := string(read(t, "custom/pricelist/position.json"))
	positionsPath := "/entity/pricelist/" + testID + "/positions"
	for _, op := range []struct {
		name, request, body, response, context string
		invoke                                 func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE /entity/pricelist/" + testID, "", "", "delete pricelist", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}},
		{"trash", "POST /entity/pricelist/" + testID + "/trash", "", "", "trash pricelist", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}},
		{"create many", "POST /entity/pricelist", `[{"columns":[{"name":"without_discount"}]}]`, `[{"id":"a"}]`, "create many pricelist", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).CreateMany(ctx, []pricelist.CreateRequest{{Columns: columns}})
			if err == nil && got[0].Value.ID != "a" {
				t.Errorf("result = %+v", got)
			}
			return err
		}},
		{"update many", "POST /entity/pricelist", `[` + meta(testID) + `,"shared":true,"name":"Новое название"}]`, `[{"id":"a"}]`, "update many pricelist", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).UpdateMany(ctx, []batch.Update[pricelist.UpdateRequest]{{ID: testID, Request: pricelist.UpdateRequest{Shared: model.Value(true), Name: model.Value("Новое название")}}})
			return err
		}},
		{"delete many", "POST /entity/pricelist/delete", `[` + meta(testID) + `}]`, `[{"info":"ok"}]`, "delete many pricelist", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).DeleteMany(ctx, []string{testID})
			if err == nil && got[0].Info != "ok" {
				t.Errorf("result = %+v", got)
			}
			return err
		}},
		{"create positions", "POST " + positionsPath, `[{"assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/88888888-8888-8888-8888-888888888888","type":"product"}},"cells":[{"column":"123","sum":100000}]}]`,
			"[" + row + "]", "create many pricelist " + testID + " positions", func(ctx context.Context, c *transport.Client) error {
				got, err := service(c).CreatePositions(ctx, testID, []pricelist.PositionRequest{{Assortment: model.Value(product), Cells: model.Value([]pricelist.Cell{{Column: "123", Sum: 100000}})}})
				if err == nil && (got[0].Value.ID != positionID || got[0].Value.Cells[0].Sum != 100000) {
					t.Errorf("result = %+v", got)
				}
				return err
			}},
		{"delete position", "DELETE " + positionsPath + "/" + positionID, "", "", "delete pricelist " + testID + " position", func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, positionID)
		}},
		{"delete many positions", "POST " + positionsPath + "/delete", `[{"meta":{"href":"BASE` + positionsPath + "/" + positionID + `","type":"pricelistrow","mediaType":"application/json"}}]`, "",
			"delete pricelist " + testID + " positions", func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyPositions(ctx, testID, []string{positionID})
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
	for name, invoke := range map[string]func(context.Context, *transport.Client) error{
		"trash": func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) },
		"delete many positions": func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteManyPositions(ctx, testID, []string{positionID})
		},
	} {
		if _, _, calls, err := exchange(t, 204, "", invoke); err != nil || calls != 1 {
			t.Errorf("%s 204: calls = %d, err = %v", name, calls, err)
		}
		if _, _, calls, err := exchange(t, 202, "", invoke); err == nil || calls != 1 {
			t.Errorf("%s 202: calls = %d, err = %v", name, calls, err)
		}
	}
}

// Ровно 1000 строк уходят; пустые столбцы, 1001 строка, неверные id и пустые запросы
// отклоняются до запроса; ResponseInfo сбрасывается.
func TestPriceListValidation(t *testing.T) {
	rows := func(n int) model.Field[pricelist.PositionRows] {
		return model.Value(pricelist.PositionRows{Rows: make([]pricelist.PositionRequest, n)})
	}
	for name, invoke := range map[string]func(context.Context, *transport.Client) error{
		"create 1000": func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).Create(ctx, pricelist.CreateRequest{Columns: columns, Positions: rows(1000)})
			return err
		},
		"update 1000": func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).Update(ctx, testID, pricelist.UpdateRequest{Positions: rows(1000)})
			return err
		},
	} {
		if _, _, calls, err := exchange(t, 200, `{"id":"x"}`, invoke); err != nil || calls != 1 {
			t.Errorf("%s: calls = %d, err = %v", name, calls, err)
		}
	}
	for name, invoke := range map[string]func(context.Context, *transport.Client, ...transport.RequestOption) error{
		"no columns": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Create(ctx, pricelist.CreateRequest{Name: model.Value("x")}, o...)
			return err
		},
		"empty column name": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Create(ctx, pricelist.CreateRequest{Columns: []pricelist.ColumnRequest{{Name: "a"}, {PercentageDiscount: model.Value(5)}}}, o...)
			return err
		},
		"create 1001": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Create(ctx, pricelist.CreateRequest{Columns: columns, Positions: rows(1001)}, o...)
			return err
		},
		"update 1001": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Update(ctx, testID, pricelist.UpdateRequest{Positions: rows(1001)}, o...)
			return err
		},
		"empty update": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).Update(ctx, testID, pricelist.UpdateRequest{}, o...)
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
			_, err := service(c).UpdateMany(ctx, []batch.Update[pricelist.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"bad delete many id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).DeleteMany(ctx, []string{"bad"}, o...)
			return err
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
		"bad position id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).GetPosition(ctx, testID, "bad", o...)
			return err
		},
		"empty create positions": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).CreatePositions(ctx, testID, nil, o...)
			return err
		},
		"empty update position": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			_, err := service(c).UpdatePosition(ctx, testID, positionID, pricelist.PositionRequest{}, o...)
			return err
		},
		"bad delete position id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).DeletePosition(ctx, "bad", positionID, o...)
		},
		"empty delete many positions": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).DeleteManyPositions(ctx, testID, nil, o...)
		},
		"bad delete many positions id": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).DeleteManyPositions(ctx, testID, []string{positionID, "bad"}, o...)
		},
		"bad delete many positions document": func(ctx context.Context, c *transport.Client, o ...transport.RequestOption) error {
			return service(c).DeleteManyPositions(ctx, "bad", []string{positionID}, o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		_, _, calls, err := exchange(t, 200, `{"id":"x"}`, func(ctx context.Context, c *transport.Client) error {
			return invoke(ctx, c, transport.WithResponseInfo(&info))
		})
		if err == nil || calls != 0 || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "pricelist") {
			t.Errorf("%s: calls = %d, info = %+v, err = %v", name, calls, info, err)
		}
	}
	if _, err := pricelist.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestPriceListList(t *testing.T) {
	listtest.Run(t, "pricelist", read(t, "custom/pricelist/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[pricelist.PriceList] {
		return service(c).List
	}, func(t *testing.T, got pricelist.PriceList) {
		if got.ID != fixtureID || len(got.Columns) != 2 || got.PriceType == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
