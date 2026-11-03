// Проверка O04: удаление в корзину у всех документов — POST …/trash без тела, 200 или 204.
package entity_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/entity/commissionreportin"
	"github.com/eteeet/moysklad-sdk-go/entity/commissionreportout"
	"github.com/eteeet/moysklad-sdk-go/entity/counterpartyadjustment"
	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
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
	"github.com/eteeet/moysklad-sdk-go/entity/payroll"
	"github.com/eteeet/moysklad-sdk-go/entity/prepayment"
	"github.com/eteeet/moysklad-sdk-go/entity/prepaymentreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashin"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildrawercashout"
	"github.com/eteeet/moysklad-sdk-go/entity/retailsalesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/retailshift"
	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestDocumentTrash(t *testing.T) {
	for _, tc := range []struct {
		document string
		trash    func(context.Context, *transport.Client, string, ...transport.RequestOption) error
	}{
		{"cashin", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := cashin.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"cashout", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := cashout.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"customerorder", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := customerorder.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"demand", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := demand.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"enter", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := enter.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"invoicein", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := invoicein.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"invoiceout", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := invoiceout.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"loss", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := loss.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"move", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := move.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"paymentin", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := paymentin.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"paymentout", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := paymentout.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"purchaseorder", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := purchaseorder.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"purchasereturn", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := purchasereturn.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"retaildemand", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := retaildemand.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"retailsalesreturn", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := retailsalesreturn.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"retailshift", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := retailshift.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"salesreturn", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := salesreturn.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"supply", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := supply.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"retaildrawercashout", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := retaildrawercashout.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"retaildrawercashin", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := retaildrawercashin.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"prepaymentreturn", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := prepaymentreturn.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"prepayment", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := prepayment.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"payroll", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := payroll.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"commissionreportout", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := commissionreportout.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"commissionreportin", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := commissionreportin.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"counterpartyadjustment", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := counterpartyadjustment.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"factureout", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := factureout.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"facturein", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := facturein.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"inventory", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := inventory.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
		{"internalorder", func(ctx context.Context, c *transport.Client, id string, o ...transport.RequestOption) error {
			s, _ := internalorder.NewService(c)
			return s.Trash(ctx, id, o...)
		}},
	} {
		t.Run(tc.document, func(t *testing.T) {
			for _, c := range []struct {
				status int
				body   string
				ok     bool
			}{{200, "", true}, {204, "", true}, {403, `{"errors":[{"error":"forbidden","code":1016}]}`, false}, {202, "", false}} {
				var calls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != http.MethodPost || r.URL.Path != "/entity/"+tc.document+"/"+testID+"/trash" || len(b) != 0 {
						t.Errorf("request = %s %s %s", r.Method, r.URL.Path, b)
					}
					w.WriteHeader(c.status)
					io.WriteString(w, c.body)
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.trash(t.Context(), client, testID)
				server.Close()
				if (err == nil) != c.ok || calls.Load() != 1 {
					t.Fatalf("status %d: err = %v, calls = %d", c.status, err, calls.Load())
				}
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "trash "+tc.document)) {
					t.Errorf("error = %v", err)
				}
			}
			info := transport.ResponseInfo{StatusCode: 503}
			local, err := transport.NewClient(transport.Config{BaseURL: "http://127.0.0.1:1"})
			if err != nil {
				t.Fatal(err)
			}
			if err := tc.trash(t.Context(), local, testID+"/x", transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
				t.Fatalf("bad id: err = %v, info = %+v", err, info)
			}
		})
	}
}
