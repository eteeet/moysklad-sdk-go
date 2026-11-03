package turnover_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestOperationRowJSON(t *testing.T) {
	data := read(t, "upstream/report_turnover_by_operation.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got turnover.OperationRow
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.Operation.Meta.Type != "supply" || got.Operation.Name != "00002" || got.Operation.Moment.IsZero() || got.Operation.Agent == nil ||
				got.Store.Name != "Основной склад" || got.Quantity != 3 || got.Cost != 10000 || got.Sum != 30000 {
				t.Fatalf("row = %+v", got)
			}
			var page model.Collection[turnover.OperationRow]
			if err := codec.Unmarshal(read(t, "custom/turnover/byoperations.json"), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Rows) != 4 || page.Context == nil || page.Rows[0].Operation.Meta.Href == "" || page.Rows[0].Operation.Moment.IsZero() {
				t.Fatalf("page = %+v", page)
			}
			var future turnover.OperationRow
			if err := codec.Unmarshal([]byte(`{"operation":{"meta":{"href":"h","type":"futuredocument"},"name":"1","moment":"2026-11-01 10:00:00.000"},"quantity":-1}`), &future); err != nil ||
				future.Operation.Meta.Type != "futuredocument" || future.Operation.Agent != nil || future.Quantity != -1 {
				t.Fatalf("future = %+v, err = %v", future, err)
			}
		})
	}
}

func TestByOperationsHTTP(t *testing.T) {
	typ, err := query.Filter(query.Eq("type", "supply"))
	if err != nil {
		t.Fatal(err)
	}
	options := turnover.DetailOptions{Assortment: ref("entity/product/p1", "product"), Filter: typ,
		MomentFrom: documenttest.DateTime("2026-10-01 00:00:00.000"), MomentTo: documenttest.DateTime("2026-11-01 00:00:00.000")}
	const q = "filter=product%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2Fp1%3Btype%3Dsupply&momentFrom=2026-10-01+00%3A00%3A00&momentTo=2026-11-01+00%3A00%3A00"
	for _, tc := range []struct {
		name    string
		options turnover.DetailOptions
		query   string
	}{
		{"explicit period", options, q},
		{"server period", turnover.DetailOptions{Assortment: ref("entity/variant/v1", "variant")},
			"filter=variant%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fvariant%2Fv1"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[turnover.OperationRow], error) {
			return service(c).ByOperations(ctx, tc.options)
		}
		for _, rc := range documenttest.ResponseCases(page(read(t, "upstream/report_turnover_by_operation.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/turnover/byoperations?"+tc.query, nil, invoke, func(got *model.Collection[turnover.OperationRow]) {
					if len(got.Rows) != 1 || got.Rows[0].Sum != 30000 {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestByOperationsValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for name, options := range map[string]turnover.DetailOptions{
		"no assortment": {},
		"bundle":        {Assortment: ref("entity/bundle/b1", "bundle")},
	} {
		if _, err := service(c).ByOperations(t.Context(), options); err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "turnover byoperations") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
