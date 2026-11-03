package positions_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	mutDoc  = "7944ef04-f831-11e5-7a69-971500188b19"
	mutRow  = "34f6344f-015e-11e6-9464-e4de0000006c"
	mutPath = "entity/customerorder/" + mutDoc + "/positions"
	mutJSON = `{"meta":{"type":"customerorderposition"},"id":"` + mutRow + `","quantity":49,"price":12345.0,"reserve":19}`
)

func rowService(c *transport.Client) *positions.Service[positions.Position] {
	s, _ := positions.NewService[positions.Position](c, "customerorder")
	return s
}

// Get и Update строки: общий HTTP-контракт; Update принимает и массив из одного объекта.
func TestPositionHTTP(t *testing.T) {
	id := func(p *positions.Position, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &p.ID, nil
	}
	for _, o := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}{
		{"get", "GET", mutPath + "/" + mutRow, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(rowService(c).GetPosition(ctx, mutDoc, mutRow))
		}},
		{"update", "PUT", mutPath + "/" + mutRow, `{"quantity":0,"reserve":5}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(rowService(c).UpdatePosition(ctx, mutDoc, mutRow, positions.PositionRequest{
				ItemRequest: positions.ItemRequest{Quantity: model.Value(0.0)}, Reserve: model.Value(5.0)}))
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			var body []byte
			if o.body != "" {
				body = []byte(o.body)
			}
			for _, tc := range documenttest.ResponseCases(mutJSON) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, body, o.invoke, func(got *string) {
						if *got != mutRow {
							t.Fatalf("id = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
}

// Создание строк массивом и удаление: путь, тело, meta с типом строки документа; массовое удаление —
// 200 с пустым телом (так в источнике).
func TestPositionWrites(t *testing.T) {
	var got []string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.Path+" "+string(b))
		if r.Method == http.MethodPost && !strings.HasSuffix(r.URL.Path, "/delete") {
			io.WriteString(w, `[`+mutJSON+`]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	ctx := t.Context()
	product := model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/66ccbc9f-f7d2-11e5-8a84-bae500000076", Type: "product"}}
	r, err := rowService(c).CreatePositions(ctx, mutDoc, []positions.PositionRequest{{
		ItemRequest: positions.ItemRequest{Assortment: model.Value(product), Quantity: model.Value(49.0)},
		Price:       model.Value(12345.0), Discount: model.Value(0.0), VAT: model.Value(18), Reserve: model.Value(19.0)}})
	if err != nil || len(r) != 1 || r[0].Value.ID != mutRow {
		t.Fatalf("create = %+v, err = %v", r, err)
	}
	if err := rowService(c).DeletePosition(ctx, mutDoc, mutRow); err != nil {
		t.Fatal(err)
	}
	if err := rowService(c).DeleteManyPositions(ctx, mutDoc, []string{mutRow}); err != nil {
		t.Fatal(err)
	}
	invoice, _ := positions.NewService[positions.Position](c, "invoiceout")
	if err := invoice.DeleteManyPositions(ctx, mutDoc, []string{mutRow}); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"POST /" + mutPath + ` [{"assortment":{"meta":{"href":"` + product.Meta.Href + `","type":"product"}},"quantity":49,"price":12345,"discount":0,"vat":18,"reserve":19}]`,
		"DELETE /" + mutPath + "/" + mutRow + " ",
		"POST /" + mutPath + `/delete [{"meta":{"href":"` + server.URL + "/" + mutPath + "/" + mutRow + `","type":"customerorderposition","mediaType":"application/json"}}]`,
		"POST /entity/invoiceout/" + mutDoc + `/positions/delete [{"meta":{"href":"` + server.URL + "/entity/invoiceout/" + mutDoc + "/positions/" + mutRow + `","type":"invoiceposition","mediaType":"application/json"}}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestDeletePositionsResponses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		ok     bool
	}{
		{"empty 200", 200, "", true},
		{"204", 204, "", true},
		{"200 body ignored", 200, `[{"info":"ok"}]`, true},
		{"api error", 400, `{"errors":[{"error":"bad","code":3000}]}`, false},
		{"202", 202, "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			err = rowService(c).DeleteManyPositions(t.Context(), mutDoc, []string{mutRow})
			if (err == nil) != tc.ok || calls.Load() != 1 {
				t.Fatalf("err = %v, calls = %d", err, calls.Load())
			}
			if err != nil && strings.Count(err.Error(), "moysklad:") != 1 {
				t.Errorf("error = %v", err)
			}
		})
	}
}

func TestPositionValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := rowService(c), t.Context()
	many := make([]string, 1001)
	for i := range many {
		many[i] = mutRow
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad doc": func(o ...transport.RequestOption) error {
			_, err := s.GetPosition(ctx, "bad", mutRow, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, mutDoc, mutRow, positions.PositionRequest{}, o...)
			return err
		},
		"create empty": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, mutDoc, nil, o...)
			return err
		},
		"delete bad row":    func(o ...transport.RequestOption) error { return s.DeletePosition(ctx, mutDoc, "x", o...) },
		"delete many empty": func(o ...transport.RequestOption) error { return s.DeleteManyPositions(ctx, mutDoc, nil, o...) },
		"delete many 1001":  func(o ...transport.RequestOption) error { return s.DeleteManyPositions(ctx, mutDoc, many, o...) },
		"delete many bad": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, mutDoc, []string{"x"}, o...)
		},
		"delete many doc": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, "../x", []string{mutRow}, o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
