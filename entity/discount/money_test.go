package discount_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

// Денежные Level.Amount и SpecialPrice.Value: 0, 0.0 и дробь читаются обоими кодеками; запись
// сохраняет 0 и дробь.
func TestMoneyForms(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			number string
			want   float64
		}{{"0", 0}, {"0.0", 0}, {"1050.5", 1050.5}} {
			var level discount.Level
			if err := codec.Unmarshal([]byte(`{"amount":`+tc.number+`,"discount":`+tc.number+`}`), &level); err != nil || level.Amount != tc.want || level.Discount != tc.want {
				t.Fatalf("%s level %s: %+v, err = %v", codec.Name, tc.number, level, err)
			}
			var price discount.SpecialPrice
			if err := codec.Unmarshal([]byte(`{"value":`+tc.number+`}`), &price); err != nil || price.Value != tc.want {
				t.Fatalf("%s price %s: %+v, err = %v", codec.Name, tc.number, price, err)
			}
		}
		data, err := codec.Marshal([]any{discount.Level{Amount: 0, Discount: 2.5}, discount.Level{Amount: 1050.5},
			discount.SpecialPriceRequest{Value: model.Value(0.0)}, discount.SpecialPriceRequest{Value: model.Value(99.99)}})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`[{"amount":0,"discount":2.5},{"amount":1050.5,"discount":0},{"value":0},{"value":99.99}]`))
	}
}
