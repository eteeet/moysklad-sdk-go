package prepaymentreturn_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/prepaymentreturn"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// Денежные поля — копейки float64: 0, 0.0, дробь и большое значение читаются без потерь обоими кодеками.
func TestMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got prepaymentreturn.PrepaymentReturn
			if err := codec.Unmarshal([]byte(`{"sum":`+raw+`}`), &got); err != nil || got.Sum != want {
				t.Fatalf("%s: sum=%s: got %v, err = %v", codec.Name, raw, got.Sum, err)
			}
		}
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got prepaymentreturn.PrepaymentReturn
			if err := codec.Unmarshal([]byte(`{"cashSum":`+raw+`}`), &got); err != nil || got.CashSum != want {
				t.Fatalf("%s: cashSum=%s: got %v, err = %v", codec.Name, raw, got.CashSum, err)
			}
		}
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got prepaymentreturn.PrepaymentReturn
			if err := codec.Unmarshal([]byte(`{"noCashSum":`+raw+`}`), &got); err != nil || got.NoCashSum != want {
				t.Fatalf("%s: noCashSum=%s: got %v, err = %v", codec.Name, raw, got.NoCashSum, err)
			}
		}
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got prepaymentreturn.PrepaymentReturn
			if err := codec.Unmarshal([]byte(`{"qrSum":`+raw+`}`), &got); err != nil || got.QRSum != want {
				t.Fatalf("%s: qrSum=%s: got %v, err = %v", codec.Name, raw, got.QRSum, err)
			}
		}
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got prepaymentreturn.PrepaymentReturn
			if err := codec.Unmarshal([]byte(`{"vatSum":`+raw+`}`), &got); err != nil || got.VATSum != want {
				t.Fatalf("%s: vatSum=%s: got %v, err = %v", codec.Name, raw, got.VATSum, err)
			}
		}
	}
}
