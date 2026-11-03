package product_test

import (
	"os"
	"strings"
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
			// Пример изменения товара источника: существующая упаковка по id.
			{"update existing pack", product.UpdateRequest{Packs: model.Value([]product.PackRequest{{ID: "354ba98c-2cb9-11e6-8a84-bae5000000e3", UOM: uom, Quantity: 2}})},
				`{"packs":[{"id":"354ba98c-2cb9-11e6-8a84-bae5000000e3","uom":{"meta":{"href":"` + uom.Meta.Href + `","type":"uom"}},"quantity":2}]}`},
		} {
			data, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, data, []byte(tc.want))
		}
	}
}

func TestProductAttributes(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		got := fullProduct(t, codec)
		if got.Alcoholic == nil || *got.Alcoholic != (product.Alcoholic{Excise: true, Type: 10, Strength: 15.5, Volume: 15.5}) ||
			!got.PartialDisposal || got.PaymentItemType != "GOOD" || got.PPEType != "2400001226405" || got.TrackingType != "MILK" ||
			len(got.Things) != 1 || got.Things[0] != "1" || got.VariantsCount != 10 || got.Weighed || got.OnTap || got.NctKZ != nil ||
			got.Images == nil || got.Files == nil {
			t.Fatalf("%s: product = %+v", codec.Name, got)
		}
		// weighed и onTap выводятся только при true; mod__nct__kz — объект Казахстана.
		var p product.Product
		if err := codec.Unmarshal([]byte(`{"weighed":true,"onTap":true,"minimumBalance":2.5,"mod__nct__kz":{"ntin":"0001"},"mod__marking__uz":{"trackingType":"WATER"}}`), &p); err != nil ||
			!p.Weighed || !p.OnTap || p.MinimumBalance != 2.5 || p.NctKZ == nil || p.NctKZ.NTIN != "0001" || p.MarkingUZ == nil || p.MarkingUZ.TrackingType != "WATER" {
			t.Fatalf("%s: product = %+v, err = %v", codec.Name, p, err)
		}
		for _, tc := range []struct {
			name  string
			input any
			want  string
		}{
			{"create", product.CreateRequest{Name: "Вино", Alcoholic: model.Value(product.AlcoholicRequest{Excise: model.Value(true), Type: model.Value(200), Strength: model.Value(12.5), Volume: model.Value(0.75)}),
				TrackingType: model.Value("NOT_TRACKED"), PaymentItemType: model.Value("GOOD"), TaxSystem: model.Value("GENERAL_TAX_SYSTEM"), TNVED: model.Value("2204"),
				PartialDisposal: model.Value(false), IsSerialTrackable: model.Value(false), MinimumBalance: model.Value(0.0)},
				`{"name":"Вино","alcoholic":{"excise":true,"type":200,"strength":12.5,"volume":0.75},"trackingType":"NOT_TRACKED","paymentItemType":"GOOD",` +
					`"taxSystem":"GENERAL_TAX_SYSTEM","tnved":"2204","partialDisposal":false,"isSerialTrackable":false,"minimumBalance":0}`},
			{"update", product.UpdateRequest{Weighed: model.Value(true), OnTap: model.Value(false), Things: model.Value([]string{"SN-1"}), PPEType: model.Null[string](),
				Alcoholic: model.Null[product.AlcoholicRequest](), NctKZ: model.Value(model.NctKZ{NTIN: "0001"}), TasnifUZ: model.Value(model.TasnifUZ{IKPU: "1"})},
				`{"weighed":true,"onTap":false,"things":["SN-1"],"ppeType":null,"alcoholic":null,"mod__nct__kz":{"ntin":"0001"},"mod__tasnif__uz":{"ikpu":"1"}}`},
		} {
			data, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, data, []byte(tc.want))
		}
	}
}

func TestMinimumStockField(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		got := fullProduct(t, codec)
		if got.MinimumStock == nil || got.MinimumStock.Type != model.MinimumStockWarehouseVaried || got.MinimumStock.StoreBalances == nil ||
			!strings.HasSuffix(got.MinimumStock.StoreBalances.Meta.Href, "/storebalances") {
			t.Fatalf("%s: minimumStock = %+v", codec.Name, got.MinimumStock)
		}
		data, err := codec.Marshal(product.CreateRequest{Name: "Картофель", MinimumStock: model.Value(model.MinimumStockRequest{Type: model.MinimumStockAllWarehouseSum, Quantity: model.Value(10.0)})})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`{"name":"Картофель","minimumStock":{"type":"ALL_WAREHOUSE_SUM","quantity":10}}`))
	}
}
