package payroll_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/payroll"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// Денежные поля — копейки float64: 0, 0.0, дробь и большое значение читаются без потерь обоими кодеками.
func TestMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got payroll.Position
			if err := codec.Unmarshal([]byte(`{"baseSalary":`+raw+`}`), &got); err != nil || got.BaseSalary != want {
				t.Fatalf("%s: baseSalary=%s: got %v, err = %v", codec.Name, raw, got.BaseSalary, err)
			}
		}
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got payroll.Position
			if err := codec.Unmarshal([]byte(`{"pieceworkSalary":`+raw+`}`), &got); err != nil || got.PieceworkSalary != want {
				t.Fatalf("%s: pieceworkSalary=%s: got %v, err = %v", codec.Name, raw, got.PieceworkSalary, err)
			}
		}
	}
}
