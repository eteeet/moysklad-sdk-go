package product_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func fullProduct(t *testing.T, codec documenttest.Codec) product.Product {
	t.Helper()
	data, err := os.ReadFile("../../testdata/upstream/product_full.json")
	if err != nil {
		t.Fatal(err)
	}
	var got product.Product
	if err := codec.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	return got
}

func TestBarcodesAndPacks(t *testing.T) {
	uom := model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/uom/89f7fa51-2271-4b87-baf8-c3fc2a3f5119", Type: "uom"}}
	for _, codec := range documenttest.Codecs {
		got := fullProduct(t, codec)
		if len(got.Barcodes) != 2 || got.Barcodes[0]["ean13"] != "4600000000001" || got.Barcodes[1]["code128"] != "CODE128-001" {
			t.Fatalf("%s: barcodes = %+v", codec.Name, got.Barcodes)
		}
		if len(got.Packs) != 1 || got.Packs[0].ID != "4511ceca-0347-11f1-0a81-072700000203" || got.Packs[0].Quantity != 1 ||
			got.Packs[0].UOM == nil || got.Packs[0].UOM.Meta.Type != "uom" || len(got.Packs[0].Barcodes) != 1 || got.Packs[0].Barcodes[0]["ean13"] != "2000000000039" {
			t.Fatalf("%s: packs = %+v", codec.Name, got.Packs)
		}
		for _, tc := range []struct {
			name  string
			input any
			want  string
		}{
			{"create", product.CreateRequest{Name: "P", Barcodes: model.Value([]model.Barcode{{"ean13": "4600000000001"}, {"gtin": "04600000000001"}}),
				Packs: model.Value([]product.PackRequest{{UOM: uom, Quantity: 10, Barcodes: model.Value([]model.Barcode{{"ean13": "2000000000039"}})}})},
				`{"name":"P","barcodes":[{"ean13":"4600000000001"},{"gtin":"04600000000001"}],"packs":[{"uom":{"meta":{"href":"` + uom.Meta.Href + `","type":"uom"}},"quantity":10,"barcodes":[{"ean13":"2000000000039"}]}]}`},
			{"create no barcodes", product.CreateRequest{Name: "P", Barcodes: model.Value([]model.Barcode{{"ean13": ""}})}, `{"name":"P","barcodes":[{"ean13":""}]}`},
			{"update clear", product.UpdateRequest{Barcodes: model.Value([]model.Barcode{}), Packs: model.Value([]product.PackRequest{})}, `{"barcodes":[],"packs":[]}`},
			{"update pack without barcode", product.UpdateRequest{Packs: model.Value([]product.PackRequest{{UOM: uom, Quantity: 0.5}})},
				`{"packs":[{"uom":{"meta":{"href":"` + uom.Meta.Href + `","type":"uom"}},"quantity":0.5}]}`},
		} {
			data, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, data, []byte(tc.want))
		}
	}
}
