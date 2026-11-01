package positions_test

import (
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

func ref(href, typ string) *model.Reference {
	return &model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + href, Type: model.MetaType(typ)}}
}

// D01: поля строк заказа покупателя, оприходования, списания, возвратов и розницы читаются
// в общую форму; вложенные коды маркировки и остатки — своими типами.
func TestPositionExtensionsJSON(t *testing.T) {
	const input = `{
		"id":"row","quantity":2,"price":100,"reserve":1.5,"taxSystem":"GENERAL_TAX_SYSTEM","reason":"брак",
		"slot":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/s/slots/1","type":"slot"}},
		"pack":{"id":"pack-1","quantity":6,"uom":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/uom/u","type":"uom"}},"barcodes":[{"ean13":"2000000000015"}]},
		"country":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/country/c","type":"country"}},
		"gtd":{"name":"12345678/121217/1212321"},
		"things":["s1","s2"],
		"trackingCodes":[{"cis":"box","type":"transportpack","trackingCodes":[{"cis":"item","type":"trackingcode"}]}],
		"trackingCodes_1162":[{"cis_1162":"1162","type":"consumerpack","trackingCodes_1162":[]}],
		"declaration":[{"gtd":"10130000/010120/0000001","rnpt":"","country":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/country/c","type":"country"}},"quantity":2}],
		"stock":{"cost":5000,"quantity":999,"reserve":90,"intransit":9,"available":918},
		"future":1
	}`
	cost := 5000.0
	want := positions.Position{
		Item:    positions.Item{Core: positions.Core{ID: "row"}, Quantity: 2},
		Price:   100,
		Reserve: 1.5, TaxSystem: "GENERAL_TAX_SYSTEM", Reason: "брак",
		Slot:    ref("entity/store/s/slots/1", "slot"),
		Pack:    &positions.Pack{ID: "pack-1", Quantity: 6, UOM: ref("entity/uom/u", "uom"), Barcodes: []model.Barcode{{"ean13": "2000000000015"}}},
		Country: ref("entity/country/c", "country"),
		GTD:     &positions.GTD{Name: "12345678/121217/1212321"},
		Things:  []string{"s1", "s2"},
		TrackingCodes: []positions.TrackingCode{{CIS: "box", Type: "transportpack",
			TrackingCodes: []positions.TrackingCode{{CIS: "item", Type: "trackingcode"}}}},
		TrackingCodes1162: []positions.TrackingCode1162{{CIS1162: "1162", Type: "consumerpack", TrackingCodes1162: []positions.TrackingCode1162{}}},
		Declaration:       []positions.Declaration{{GTD: "10130000/010120/0000001", Country: ref("entity/country/c", "country"), Quantity: 2}},
		Stock:             &positions.Stock{Cost: &cost, Quantity: 999, Reserve: 90, InTransit: 9, Available: 918},
	}
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			var got positions.Position
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("position = %+v\nwant %+v", got, want)
			}
			var noCost positions.Position
			if err := codec.unmarshal([]byte(`{"stock":{"quantity":1,"reserve":0,"intransit":0,"available":1}}`), &noCost); err != nil ||
				noCost.Stock == nil || noCost.Stock.Cost != nil || noCost.Stock.Available != 1 {
				t.Fatalf("stock without cost = %+v, err = %v", noCost.Stock, err)
			}
		})
	}
}

func TestPositionExtensionsRequestJSON(t *testing.T) {
	input := positions.PositionRequest{
		Reserve: model.Value(0.0), TaxSystem: model.Value("PATENT_BASED"), Reason: model.Value(""),
		Slot:    model.Null[model.Reference](),
		Pack:    model.Value(positions.PackRequest{ID: "pack-1"}),
		Country: model.Value(*ref("entity/country/c", "country")),
		GTD:     model.Value(positions.GTD{Name: "12345678/121217/1212321"}),
		Things:  model.Value([]string{}),
		TrackingCodes: model.Value([]positions.TrackingCode{{CIS: "box", Type: "transportpack",
			TrackingCodes: []positions.TrackingCode{{CIS: "item", Type: "trackingcode"}}}}),
	}
	const want = `{"reserve":0,"taxSystem":"PATENT_BASED","reason":"","slot":null,"pack":{"id":"pack-1"},
		"country":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/country/c","type":"country"}},
		"gtd":{"name":"12345678/121217/1212321"},"things":[],
		"trackingCodes":[{"cis":"box","type":"transportpack","trackingCodes":[{"cis":"item","type":"trackingcode"}]}]}`
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, data, want)
		})
	}
	for _, name := range []string{"Declaration", "Stock", "TrackingCodes1162", "Shipped", "Overhead"} {
		if _, ok := reflect.TypeFor[positions.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains read-only %s", name)
		}
	}
}

// fields у позиций — только документированные значения для кода документа, одним значением.
func TestServiceGetPageFields(t *testing.T) {
	for _, tc := range []struct {
		document, fields string
		ok               bool
	}{
		{"demand", "stock", true},
		{"customerorder", "stock", true},
		{"salesreturn", "stock", true},
		{"loss", "declaration", true},
		{"retaildemand", "declaration", true},
		{"loss", "stock", false},
		{"move", "stock", false},
		{"enter", "stock", false},
		{"payroll", "stock", false},
		{"retaildemand", "stock,declaration", false},
	} {
		t.Run(tc.document+"/"+tc.fields, func(t *testing.T) {
			var calls atomic.Int32
			service := readerNewService(t, tc.document, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.URL.RawQuery != "fields="+tc.fields+"&limit=1000&offset=0" {
					t.Errorf("query = %s", r.URL.RawQuery)
				}
				readerWriteResponse(t, w, `{"meta":{"href":"h"},"rows":[]}`, false)
			}))
			_, err := service.GetPage(t.Context(), readerDocumentID, positions.PageOptions{Fields: tc.fields})
			want := int32(0)
			if tc.ok {
				want = 1
			}
			if tc.ok != (err == nil) || calls.Load() != want || err != nil && strings.Count(err.Error(), "moysklad:") != 1 {
				t.Fatalf("err = %v, calls = %d", err, calls.Load())
			}
		})
	}
}
