package processingplan_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/processingplan"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const productID = "120b4591-b0bd-11e6-5bed-427b00000001"

// Пример ответа «Получить продукт» источника — массив из одного объекта: GetProduct его принимает.
func TestProductFixture(t *testing.T) {
	data := read(t, "custom/processingplan/product.json")
	check := func(got *processingplan.Product) {
		if got.ID != productID || got.AccountID != "d55cbfba-91f1-11e6-5bed-427b00000000" || got.Meta.Type != "processingplanresult" || got.Quantity != 1 ||
			got.Assortment == nil || got.Product == nil || got.Product.Meta.Type != "product" {
			t.Fatalf("product = %+v", got)
		}
	}
	for _, codec := range documenttest.Codecs {
		var rows []processingplan.Product
		if err := codec.Unmarshal(data, &rows); err != nil || len(rows) != 1 {
			t.Fatalf("%s: rows = %+v, err = %v", codec.Name, rows, err)
		}
		check(&rows[0])
	}
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(data) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	got, err := service(client).GetProduct(t.Context(), testID, productID)
	if err != nil {
		t.Fatal(err)
	}
	check(got)
}

func TestProductsHTTP(t *testing.T) {
	const row = `{"meta":{"type":"processingplanresult"},"id":"` + productID + `","quantity":5,"assortment":{"meta":{"type":"variant"}},"product":{"meta":{"type":"product"}}}`
	check := func(got *processingplan.Product) {
		if got.ID != productID || got.Quantity != 5 || got.Assortment == nil || got.Assortment.Meta.Type != "variant" || got.Product == nil {
			t.Fatalf("product = %+v", got)
		}
	}
	t.Run("get products", func(t *testing.T) {
		pageCases(t, planPath+"/products", "limit=1000&offset=40", `{"meta":{"href":"h","type":"processingplanresult","size":1},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingplan.Product], error) {
				return service(c).GetProducts(ctx, testID, query.PageOptions{Offset: 40})
			}, func(got *model.Collection[processingplan.Product]) { check(&got.Rows[0]) })
	})
	for name, response := range map[string]string{"get product object": row, "get product array": `[` + row + `]`} {
		t.Run(name, func(t *testing.T) {
			httpCases(t, "GET", planPath+"/products/"+productID, "", response, func(ctx context.Context, c *transport.Client) (*processingplan.Product, error) {
				return service(c).GetProduct(ctx, testID, productID)
			}, check)
		})
	}
	// Пример источника: {"quantity": 5}, ответ — массив из одного объекта.
	t.Run("update product", func(t *testing.T) {
		httpCases(t, "PUT", planPath+"/products/"+productID, `{"quantity":5}`, `[`+row+`]`, func(ctx context.Context, c *transport.Client) (*processingplan.Product, error) {
			return service(c).UpdateProduct(ctx, testID, productID, processingplan.ProductRequest{Quantity: model.Value(5.0)})
		}, check)
	})
}

// GetProduct: массив другого размера или с битым элементом — ошибка с контекстом операции.
func TestGetProductArray(t *testing.T) {
	for name, body := range map[string]string{
		"empty array": `[]`,
		"two objects": `[{"id":"` + productID + `"},{"id":"` + rowID + `"}]`,
		"no id":       `[{"quantity":1}]`,
		"bad element": `[1]`,
	} {
		server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, body) }))
		client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
		if err != nil {
			t.Fatal(err)
		}
		got, err := service(client).GetProduct(t.Context(), testID, productID)
		server.Close()
		if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get processingplan "+testID+" product "+productID) {
			t.Errorf("%s: got = %+v, err = %v", name, got, err)
		}
	}
}

func TestProductsWrites(t *testing.T) {
	variant := ref("entity/variant/0de151c1-acdc-11e6-5bed-427b00000081", "variant")
	runWrites(t, []writeCase{
		// Пример источника: два продукта.
		{"create products", "POST", planPath + "/products",
			`[{"assortment":` + product + `,"quantity":1},{"assortment":{"meta":{"href":"` + variant.Meta.Href + `","type":"variant"}},"quantity":2}]`,
			`[{"id":"` + productID + `","quantity":1},{"id":"` + rowID + `","quantity":2}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateProducts(ctx, testID, []processingplan.ProductRequest{
					{Assortment: model.Value(productRef), Quantity: model.Value(1.0)},
					{Assortment: model.Value(variant), Quantity: model.Value(2.0)},
				})
				if err == nil && (r[0].Value.ID != productID || r[1].Value.Quantity != 2) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete product", "DELETE", planPath + "/products/" + productID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteProduct(ctx, testID, productID)
		}},
		{"delete many products", "POST", planPath + "/products/delete", `[` + rowMeta("products", "processingplanresult", productID) + `,` + rowMeta("products", "processingplanresult", rowID) + `]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyProducts(ctx, testID, []string{productID, rowID})
			}},
	})
}

func TestProductsValidation(t *testing.T) {
	noHTTP(t, func(s *processingplan.Service, ctx context.Context) map[string]func(...transport.RequestOption) error {
		return map[string]func(...transport.RequestOption) error{
			"list bad id": func(o ...transport.RequestOption) error {
				_, err := s.GetProducts(ctx, "bad", query.PageOptions{}, o...)
				return err
			},
			"get bad plan": func(o ...transport.RequestOption) error {
				_, err := s.GetProduct(ctx, "bad", productID, o...)
				return err
			},
			"get bad product": func(o ...transport.RequestOption) error {
				_, err := s.GetProduct(ctx, testID, productID+"/x", o...)
				return err
			},
			"create empty": func(o ...transport.RequestOption) error {
				_, err := s.CreateProducts(ctx, testID, []processingplan.ProductRequest{}, o...)
				return err
			},
			"update empty": func(o ...transport.RequestOption) error {
				_, err := s.UpdateProduct(ctx, testID, productID, processingplan.ProductRequest{}, o...)
				return err
			},
			"delete bad": func(o ...transport.RequestOption) error { return s.DeleteProduct(ctx, testID, "bad", o...) },
			"delete many bad row": func(o ...transport.RequestOption) error {
				return s.DeleteManyProducts(ctx, testID, []string{"bad"}, o...)
			},
		}
	})
}
