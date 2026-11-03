package stock_test

import (
	"context"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/stock"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const orderID = "34efe2ee-015e-11e6-9464-e4de0000006b"

func TestOperationRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got stock.OperationRow
			if err := codec.Unmarshal(read(t, "upstream/stock_by_operation.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got.Meta.Type != "customerorder" || len(got.Positions) != 3 || got.Positions[1].Meta.Type != "variant" ||
				got.Positions[2] != (stock.OperationPosition{Meta: got.Positions[2].Meta, Name: "ТоварМногоМодификаций", Stock: 3, Reserve: 1, Quantity: 2}) {
				t.Fatalf("row = %+v", got)
			}
			var all model.Collection[stock.OperationRow]
			if err := codec.Unmarshal(read(t, "custom/stock/byoperation.json"), &all); err != nil {
				t.Fatal(err)
			}
			if len(all.Rows) != 1 || all.Meta.Type != "stockbyoperation" || all.Context == nil || len(all.Rows[0].Positions) != 6 ||
				all.Rows[0].Positions[0].Stock != -2 || all.Rows[0].Meta.Href != base+"entity/customerorder/"+orderID {
				t.Fatalf("page = %+v", all)
			}
			var future stock.OperationRow
			if err := codec.Unmarshal([]byte(`{"meta":{"href":"h","type":"futuredocument"},"positions":[{"meta":{"href":"b","type":"bundle"},"name":"B","stock":2.5,"cost":10,"future":1}]}`), &future); err != nil ||
				future.Meta.Type != "futuredocument" || future.Positions[0] != (stock.OperationPosition{Meta: model.Meta{Href: "b", Type: "bundle"}, Name: "B", Stock: 2.5, Cost: 10}) {
				t.Fatalf("future = %+v, err = %v", future, err)
			}
		})
	}
}

// Себестоимость позиции — копейки float64: 0, 0.0, дробь и большое значение без потерь.
func TestOperationMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789012.34": 123456789012.34} {
			var got stock.OperationPosition
			if err := codec.Unmarshal([]byte(`{"cost":`+raw+`}`), &got); err != nil || got.Cost != want {
				t.Fatalf("%s: cost=%s: got %v, err = %v", codec.Name, raw, got.Cost, err)
			}
		}
	}
}

func TestByOperationHTTP(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[stock.OperationRow], error) {
		return service(c).ByOperation(ctx, orderID)
	}
	for _, rc := range documenttest.ResponseCases(page("report/stock/byoperation?operation.id="+orderID, read(t, "upstream/stock_by_operation.json"))) {
		t.Run(rc.Name, func(t *testing.T) {
			documenttest.HTTP(t, rc, "GET", "report/stock/byoperation?operation.id="+orderID, nil, invoke, func(got *model.Collection[stock.OperationRow]) {
				if len(got.Rows) != 1 || len(got.Rows[0].Positions) != 3 || got.Rows[0].Positions[0].Name != "Клюшка хоккейная" {
					t.Fatalf("page = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// operation.id обязателен и должен быть UUID; без него запрос не уходит.
func TestByOperationValidation(t *testing.T) {
	c, calls := noCalls(t)
	for _, id := range []string{"", "34efe2ee", orderID + "&x=1", "34efe2ee-015e-11e6-9464-e4de0000006g"} {
		localError(t, id, "get stock byoperation", func(o transport.RequestOption) error {
			_, err := service(c).ByOperation(t.Context(), id, o)
			return err
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
