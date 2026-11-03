package retireorder_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retireorder"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// Денежные поля — копейки float64: 0, 0.0, дробь и большое значение читаются без потерь обоими кодеками.
func TestMoneyJSON(t *testing.T) {
	values := map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01}
	for _, codec := range documenttest.Codecs {
		for raw, want := range values {
			var doc retireorder.RetireOrder
			if err := codec.Unmarshal([]byte(`{"sum":`+raw+`,"vatSum":`+raw+`}`), &doc); err != nil || doc.Sum != want || doc.VATSum != want {
				t.Fatalf("%s: sum/vatSum=%s: got %v/%v, err = %v", codec.Name, raw, doc.Sum, doc.VATSum, err)
			}
			var row retireorder.Position
			if err := codec.Unmarshal([]byte(`{"price":`+raw+`}`), &row); err != nil || row.Price != want {
				t.Fatalf("%s: price=%s: got %v, err = %v", codec.Name, raw, row.Price, err)
			}
		}
	}
}
