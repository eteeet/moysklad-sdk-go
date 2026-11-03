package assortment_test

import (
	"context"
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/assortment"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestSettingsFixture(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/assortment_settings.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		var got assortment.Settings
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.Meta.Type != "assortmentsettings" || !got.BarcodeRules.FillEAN13Barcode || got.BarcodeRules.WeightBarcode ||
			got.BarcodeRules.WeightBarcodePrefix != 20 || !got.UniqueCodeRules.CheckUniqueCode || got.UniqueCodeRules.FillUniqueCode || !got.CreatedShared {
			t.Fatalf("%s: settings = %+v", codec.Name, got)
		}
	}
}

// Пример изменения источника, но с checkUniqueCode из таблицы вместо checkUniqueCodeBoolean.
func TestSettingsHTTP(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/assortment_settings.json")
	if err != nil {
		t.Fatal(err)
	}
	request := assortment.SettingsRequest{
		UniqueCodeRules: model.Value(model.UniqueCodeRulesRequest{CheckUniqueCode: model.Value(true), FillUniqueCode: model.Value(true)}),
		BarcodeRules:    model.Value(assortment.BarcodeRulesRequest{FillEAN13Barcode: model.Value(true), WeightBarcodePrefix: model.Value(55)}),
		CreatedShared:   model.Value(false),
	}
	newService := func(c *transport.Client) *assortment.Service {
		s, _ := assortment.NewService(c)
		return s
	}
	for _, op := range []struct {
		name, method string
		body         []byte
		invoke       func(context.Context, *transport.Client) (*assortment.Settings, error)
	}{
		{"get", "GET", nil, func(ctx context.Context, c *transport.Client) (*assortment.Settings, error) {
			return newService(c).GetSettings(ctx)
		}},
		{"update", "PUT", []byte(`{"uniqueCodeRules":{"checkUniqueCode":true,"fillUniqueCode":true},"barcodeRules":{"fillEAN13Barcode":true,"weightBarcodePrefix":55},"createdShared":false}`),
			func(ctx context.Context, c *transport.Client) (*assortment.Settings, error) {
				return newService(c).UpdateSettings(ctx, request)
			}},
		{"update zero prefix", "PUT", []byte(`{"barcodeRules":{"weightBarcode":false,"weightBarcodePrefix":0}}`),
			func(ctx context.Context, c *transport.Client) (*assortment.Settings, error) {
				return newService(c).UpdateSettings(ctx, assortment.SettingsRequest{BarcodeRules: model.Value(assortment.BarcodeRulesRequest{WeightBarcode: model.Value(false), WeightBarcodePrefix: model.Value(0)})})
			}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(string(data)) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, "entity/assortment/settings", op.body, op.invoke, func(got *assortment.Settings) {
						if got.BarcodeRules.WeightBarcodePrefix != 20 {
							t.Fatalf("settings = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
	c, err := transport.NewClient(transport.Config{BaseURL: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	info := transport.ResponseInfo{StatusCode: 503}
	if _, err := newService(c).UpdateSettings(t.Context(), assortment.SettingsRequest{}, transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
		t.Fatalf("empty update: err = %v, info = %+v", err, info)
	}
}
