package wizard_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
	"github.com/eteeet/moysklad-sdk-go/wizard"
)

const (
	agentHref   = "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/eff93a94-c03a-11ea-c0a8-f00c0000001f"
	productHref = "https://api.moysklad.ru/api/remap/1.2/entity/product/46628fb5-c1c8-11ea-c0a8-f00c0000001a"
	// Ответ примера «Запрос автозаполнения цен» (сокращён) с полем, не влияющим на расчёт.
	response = `{"agent":{"meta":{"href":"` + agentHref + `","type":"counterparty"}},"description":"сохраняется",` +
		`"positions":[{"assortment":{"meta":{"href":"` + productHref + `","type":"product"}},"price":3300.0,"discountedPrice":3300.0,"quantity":12.0,"sum":39600.0},` +
		`{"assortment":{"meta":{"href":"` + productHref + `","type":"product"}},"price":3300.0,"discountedPrice":2640.0,"discount":20.0}]}`
)

var template = wizard.Template{
	Agent: model.Value(model.Reference{Meta: model.Meta{Href: agentHref, Type: "counterparty"}}),
	Positions: []wizard.PositionRequest{
		{Assortment: model.Reference{Meta: model.Meta{Href: productHref, Type: "product"}}, Quantity: model.Value(12.0)},
		{Assortment: model.Reference{Meta: model.Meta{Href: productHref, Type: "product"}}, Discount: model.Value(20.0), VATEnabled: model.Value(false)},
	},
}

const requestBody = `{"agent":{"meta":{"href":"` + agentHref + `","type":"counterparty"}},"positions":[` +
	`{"assortment":{"meta":{"href":"` + productHref + `","type":"product"}},"quantity":12},` +
	`{"assortment":{"meta":{"href":"` + productHref + `","type":"product"}},"discount":20,"vatEnabled":false}]}`

func TestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		data, err := codec.Marshal(template)
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(requestBody))
		var got wizard.Result
		if err := codec.Unmarshal([]byte(response), &got); err != nil {
			t.Fatal(err)
		}
		if got.Agent == nil || len(got.Positions) != 2 || got.Positions[0].Sum != 39600 || got.Positions[1].DiscountedPrice != 2640 ||
			got.Positions[1].Discount != 20 || !strings.Contains(string(got.Raw), `"description":"сохраняется"`) {
			t.Fatalf("%s: result = %+v", codec.Name, got)
		}
		empty, err := codec.Marshal(wizard.Template{Positions: []wizard.PositionRequest{}})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, empty, []byte(`{"positions":[]}`))
	}
}

func TestFill(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		ok         bool
	}{
		{"ok", response, 200, true},
		{"null", `null`, 200, false},
		{"array", `[]`, 200, false},
		{"api error", `{"errors":[{"error":"agent required","code":3000}]}`, 412, false},
		{"unexpected status", response, 201, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				b, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodPost || r.URL.Path != "/wizard/demand" || r.URL.Query().Get("action") != "evaluate_price,evaluate_vat" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				documenttest.JSONEqual(t, b, []byte(requestBody))
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			s, _ := wizard.NewService(c)
			got, err := s.Fill(t.Context(), "demand", []string{wizard.EvaluatePrice, wizard.EvaluateVat}, template)
			if tc.ok != (err == nil) || tc.ok != (got != nil) || calls.Load() != 1 {
				t.Fatalf("got = %+v, err = %v, calls = %d", got, err, calls.Load())
			}
			if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "fill demand template")) {
				t.Errorf("error = %v", err)
			}
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*wizard.Result, error) {
		s, _ := wizard.NewService(c)
		return s.Fill(ctx, "salesreturn", []string{wizard.EvaluateCost}, wizard.Template{})
	})
}

func TestValidation(t *testing.T) {
	c, err := transport.NewClient(transport.Config{BaseURL: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := wizard.NewService(c)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"bad type": func(o ...transport.RequestOption) error {
			_, err := s.Fill(t.Context(), "../x", []string{wizard.EvaluatePrice}, template, o...)
			return err
		},
		"no action": func(o ...transport.RequestOption) error {
			_, err := s.Fill(t.Context(), "demand", nil, template, o...)
			return err
		},
		"empty action": func(o ...transport.RequestOption) error {
			_, err := s.Fill(t.Context(), "demand", []string{""}, template, o...)
			return err
		},
		"comma": func(o ...transport.RequestOption) error {
			_, err := s.Fill(t.Context(), "demand", []string{"a,b"}, template, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
}
