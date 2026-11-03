// Проверка O02: шаблон документа PUT entity/{type}/new у документов из таблицы источника.
package entity_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/entity/facturein"
	"github.com/eteeet/moysklad-sdk-go/entity/factureout"
	"github.com/eteeet/moysklad-sdk-go/entity/internalorder"
	"github.com/eteeet/moysklad-sdk-go/entity/inventory"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashin"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashout"
	"github.com/eteeet/moysklad-sdk-go/entity/retailsalesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestDocumentTemplates(t *testing.T) {
	basis := map[string]any{"customerOrder": model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/customerorder/451cb4c0-3d1d-11e6-8a84-bae500000004", Type: "customerorder"}}}
	for _, tc := range []struct {
		document string
		newFn    func(*transport.Client, map[string]any) (bool, error)
	}{
		{"cashin", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := cashin.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"cashout", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := cashout.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"demand", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := demand.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"enter", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := enter.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"invoicein", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := invoicein.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"invoiceout", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := invoiceout.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"loss", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := loss.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"move", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := move.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"paymentin", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := paymentin.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"paymentout", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := paymentout.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"purchaseorder", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := purchaseorder.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"purchasereturn", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := purchasereturn.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"retaildemand", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := retaildemand.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"retailsalesreturn", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := retailsalesreturn.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"salesreturn", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := salesreturn.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"supply", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := supply.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"retaildrawercashout", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := retaildrawercashout.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"retaildrawercashin", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := retaildrawercashin.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"factureout", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := factureout.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"facturein", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := facturein.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"inventory", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := inventory.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
		{"internalorder", func(c *transport.Client, basis map[string]any) (bool, error) {
			s, _ := internalorder.NewService(c)
			r, err := s.New(t.Context(), basis)
			return r != nil, err
		}},
	} {
		t.Run(tc.document, func(t *testing.T) {
			for _, c := range []struct {
				name, response, body string
				basis                map[string]any
				status               int
				ok                   bool
			}{
				{"empty", `{"name":"00001","applicable":true}`, "", nil, 200, true},
				{"basis", `{"name":"00001"}`, `{"customerOrder":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/customerorder/451cb4c0-3d1d-11e6-8a84-bae500000004","type":"customerorder"}}}`, basis, 200, true},
				{"null", `null`, "", nil, 200, false},
				{"array", `[]`, "", nil, 200, false},
				{"api error", `{"errors":[{"error":"bad","code":3000}]}`, "", nil, 412, false},
				{"unexpected status", `{}`, "", nil, 201, false},
			} {
				t.Run(c.name, func(t *testing.T) {
					var calls atomic.Int32
					server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls.Add(1)
						b, _ := io.ReadAll(r.Body)
						if r.Method != http.MethodPut || r.URL.Path != "/entity/"+tc.document+"/new" || string(b) != c.body {
							t.Errorf("request = %s %s %s", r.Method, r.URL.Path, b)
						}
						w.WriteHeader(c.status)
						io.WriteString(w, c.response)
					}))
					defer server.Close()
					client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
					if err != nil {
						t.Fatal(err)
					}
					got, err := tc.newFn(client, c.basis)
					if got != c.ok || (err == nil) != c.ok || calls.Load() != 1 {
						t.Fatalf("got = %v, err = %v, calls = %d", got, err, calls.Load())
					}
					if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get "+tc.document+" template")) {
						t.Errorf("error = %v", err)
					}
				})
			}
		})
	}
}

// Отмена, сбой авторизации и транспорта — ошибка с причиной, без повторов.
func TestDocumentTemplateLocalErrors(t *testing.T) {
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*supply.Supply, error) {
		s, _ := supply.NewService(c)
		return s.New(ctx, nil)
	})
}
