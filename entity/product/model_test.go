package product_test

import (
	"os"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
)

func TestProductFixtures(t *testing.T) {
	for _, file := range []string{"product.json", "product_full.json"} {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile("../../testdata/upstream/" + file)
			if err != nil {
				t.Fatal(err)
			}
			for _, codec := range codecs {
				t.Run(codec.name, func(t *testing.T) {
					var got product.Product
					if err := codec.unmarshal(data, &got); err != nil {
						t.Fatal(err)
					}
					if !time.Time(got.Updated).Equal(time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC)) {
						t.Fatalf("updated = %s", time.Time(got.Updated))
					}
					if file == "product.json" {
						if got.BuyPrice != nil || got.MinPrice != nil || len(got.SalePrices) != 0 || len(got.Attributes) != 0 {
							t.Fatal("absent prices or attributes were invented")
						}
						return
					}
					if got.BuyPrice == nil || got.BuyPrice.Value != 50000 || got.MinPrice == nil || got.MinPrice.Value != 1000 {
						t.Fatalf("prices = %+v, %+v", got.BuyPrice, got.MinPrice)
					}
					if len(got.SalePrices) != 1 || got.SalePrices[0].Value != 100000 || got.SalePrices[0].PriceType.Name != "Цена продажи" {
						t.Fatalf("salePrices = %+v", got.SalePrices)
					}
					if len(got.Attributes) != 3 {
						t.Fatalf("attributes = %d", len(got.Attributes))
					}
					if v, e := got.Attributes[0].AsDouble(); e != nil || v != 0.75 {
						t.Fatalf("double = %v, %v", v, e)
					}
					if v, e := got.Attributes[1].AsTime(); e != nil || !time.Time(v).Equal(time.Date(2016, 6, 7, 9, 52, 33, 0, time.UTC)) {
						t.Fatalf("time = %v, %v", v, e)
					}
					if v, e := got.Attributes[2].AsBool(); e != nil || v {
						t.Fatalf("bool = %v, %v", v, e)
					}
				})
			}
		})
	}
}
