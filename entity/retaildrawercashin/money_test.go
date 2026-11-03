package retaildrawercashin_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashin"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// Денежные поля — копейки float64: 0, 0.0, дробь и большое значение читаются без потерь обоими кодеками.
func TestMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got retaildrawercashin.RetailDrawerCashIn
			if err := codec.Unmarshal([]byte(`{"sum":`+raw+`}`), &got); err != nil || got.Sum != want {
				t.Fatalf("%s: sum=%s: got %v, err = %v", codec.Name, raw, got.Sum, err)
			}
		}
	}
}
