package productiontask_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/productiontask"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// Денежные поля этапа — float64 без пересчёта: 0, 0.0, дробь и большое значение читаются без
// потерь обоими кодеками.
func TestMoneyJSON(t *testing.T) {
	fields := map[string]func(productiontask.Stage) float64{
		"processingUnitCost": func(s productiontask.Stage) float64 { return s.ProcessingUnitCost },
		"labourUnitCost":     func(s productiontask.Stage) float64 { return s.LabourUnitCost },
		"standardHourCost":   func(s productiontask.Stage) float64 { return s.StandardHourCost },
	}
	for _, codec := range documenttest.Codecs {
		for field, get := range fields {
			for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
				var got productiontask.Stage
				if err := codec.Unmarshal([]byte(`{"`+field+`":`+raw+`}`), &got); err != nil || get(got) != want {
					t.Fatalf("%s: %s=%s: got %v, err = %v", codec.Name, field, raw, get(got), err)
				}
			}
		}
	}
}
