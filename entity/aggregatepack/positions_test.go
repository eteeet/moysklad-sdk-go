package aggregatepack_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/aggregatepack"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	positionID   = "e3a7eea7-bc0c-11f1-0a82-14dc000065c3"
	positionPath = "entity/aggregatepack/" + testID + "/positions"
	positionJSON = `{"meta":{"type":"aggregatepackposition"},"id":"` + positionID + `","assortment":{"meta":{"type":"product"}},"quantity":1.0,"things":["SN"]}`
	productHref  = "https://api.moysklad.ru/api/remap/1.2/entity/product/001e6a8f-b8dc-11f1-0a82-14dc00003d64"
)

func TestPositionHTTP(t *testing.T) {
	id := func(p *aggregatepack.Position, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &p.ID, nil
	}
	for _, o := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}{
		{"get", "GET", positionPath + "/" + positionID, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(service(c).GetPosition(ctx, testID, positionID))
		}},
		{"update", "PUT", positionPath + "/" + positionID, `{"quantity":2,"things":[]}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(service(c).UpdatePosition(ctx, testID, positionID, aggregatepack.PositionUpdateRequest{Quantity: model.Value(2.0), Things: model.Value([]string{})}))
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			var body []byte
			if o.body != "" {
				body = []byte(o.body)
			}
			for _, tc := range documenttest.ResponseCases(positionJSON) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, body, o.invoke, func(got *string) {
						if *got != positionID {
							t.Fatalf("id = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
	for _, codec := range documenttest.Codecs {
		var p aggregatepack.Position
		if err := codec.Unmarshal([]byte(positionJSON), &p); err != nil || p.Quantity != 1 || p.Assortment == nil || len(p.Things) != 1 {
			t.Fatalf("%s: position = %+v, err = %v", codec.Name, p, err)
		}
	}
}

// Пример источника: создание позиции {assortment, quantity: 1.0}; список, удаление и массовое удаление.
// Пример «насыщенного» создания: позиции в теле упаковки.
func TestPositionsInBody(t *testing.T) {
	product := model.Reference{Meta: model.Meta{Href: productHref, Type: "product"}}
	for _, codec := range documenttest.Codecs {
		data, err := codec.Marshal(aggregatepack.CreateRequest{Barcodes: []model.Barcode{{"ean8": "00000000"}},
			Positions: model.Value([]aggregatepack.PositionRequest{{Assortment: product, Quantity: 1, Things: model.Value([]string{"SN"})}})})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`{"barcodes":[{"ean8":"00000000"}],"positions":[{"assortment":{"meta":{"href":"`+productHref+`","type":"product"}},"quantity":1,"things":["SN"]}]}`))
	}
}

func TestPositionCollection(t *testing.T) {
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(b))
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[`+positionJSON+`]}`)
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		case r.Method == http.MethodPost:
			io.WriteString(w, `[`+positionJSON+`]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(c), t.Context()
	if l, err := s.ListPositions(ctx, testID, query.PageOptions{}); err != nil || len(l.Rows) != 1 {
		t.Fatalf("list = %+v, err = %v", l, err)
	}
	product := model.Reference{Meta: model.Meta{Href: productHref, Type: "product"}}
	if r, err := s.CreatePositions(ctx, testID, []aggregatepack.PositionRequest{{Assortment: product, Quantity: 1}}); err != nil || r[0].Value.ID != positionID {
		t.Fatalf("create = %+v, err = %v", r, err)
	}
	if err := s.DeletePosition(ctx, testID, positionID); err != nil {
		t.Fatal(err)
	}
	if r, err := s.DeleteManyPositions(ctx, testID, []string{positionID}); err != nil || r[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", r, err)
	}
	meta := `{"href":"` + server.URL + "/" + positionPath + "/" + positionID + `","type":"aggregatepackposition","mediaType":"application/json"}`
	want := []string{
		"GET /" + positionPath + "?limit=1000&offset=0 ",
		"POST /" + positionPath + ` [{"assortment":{"meta":{"href":"` + productHref + `","type":"product"}},"quantity":1}]`,
		"DELETE /" + positionPath + "/" + positionID + " ",
		"POST /" + positionPath + `/delete [{"meta":` + meta + `}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list bad": func(o ...transport.RequestOption) error {
			_, err := s.ListPositions(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"create empty": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, testID, nil, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, aggregatepack.PositionUpdateRequest{}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error { return s.DeletePosition(ctx, testID, "x", o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}
