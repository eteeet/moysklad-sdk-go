package counterparty_test

import (
	"context"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Пример источника: GET и частичный PUT настроек справочника контрагентов.
func TestSettingsHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/settings","type":"counterpartysettings","mediaType":"application/json"},` +
		`"uniqueCodeRules":{"checkUniqueCode":true,"fillUniqueCode":true},"createShared":true}`
	request := counterparty.SettingsRequest{UniqueCodeRules: model.Value(model.UniqueCodeRulesRequest{CheckUniqueCode: model.Value(true), FillUniqueCode: model.Value(true)}),
		CreateShared: model.Value(false)}
	check := func(got *counterparty.Settings) {
		if got.Meta.Type != "counterpartysettings" || !got.UniqueCodeRules.CheckUniqueCode || !got.UniqueCodeRules.FillUniqueCode || !got.CreateShared {
			t.Fatalf("settings = %+v", got)
		}
	}
	for _, op := range []struct {
		name, method string
		body         []byte
		invoke       func(context.Context, *transport.Client) (*counterparty.Settings, error)
	}{
		{"get", "GET", nil, func(ctx context.Context, c *transport.Client) (*counterparty.Settings, error) {
			return svc(c).GetSettings(ctx)
		}},
		{"update", "PUT", []byte(`{"uniqueCodeRules":{"checkUniqueCode":true,"fillUniqueCode":true},"createShared":false}`),
			func(ctx context.Context, c *transport.Client) (*counterparty.Settings, error) {
				return svc(c).UpdateSettings(ctx, request)
			}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, "entity/counterparty/settings", op.body, op.invoke, check)
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
	if _, err := svc(c).UpdateSettings(t.Context(), counterparty.SettingsRequest{}, transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
		t.Fatalf("empty update: err = %v, info = %+v", err, info)
	}
}
