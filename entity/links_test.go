// Проверка O05: связи документов (раздел «Связи с другими документами») пишутся в Create и Update.
package entity_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const linkBase = "https://api.moysklad.ru/api/remap/1.2/entity/"

func ref(typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: linkBase + typ + "/7944ef04-f831-11e5-7a69-971500188b19", Type: model.MetaType(typ)}}
}

func link(typ string) string {
	return `{"meta":{"href":"` + linkBase + typ + `/7944ef04-f831-11e5-7a69-971500188b19","type":"` + typ + `"}}`
}

func TestDocumentLinks(t *testing.T) {
	const docID = "7944ef04-f831-11e5-7a69-971500188b20"
	for _, tc := range []struct {
		document string
		update   any
		body     string
		invoke   func(*transport.Client, any) error
		create   func() any
		key      string
	}{
		{"customerorder", customerorder.UpdateRequest{PurchaseOrders: model.Value([]model.Reference{ref("purchaseorder")}), Demands: model.Value([]model.Reference{ref("demand")}), Payments: model.Value([]model.Reference{ref("paymentin")}), ProductionTasks: model.Value([]model.Reference{ref("productiontask")}), InvoicesOut: model.Value([]model.Reference{ref("invoiceout")}), Moves: model.Value([]model.Reference{ref("move")}), Prepayments: model.Value([]model.Reference{ref("prepayment")})}, `{"purchaseOrders":[` + link("purchaseorder") + `],"demands":[` + link("demand") + `],"payments":[` + link("paymentin") + `],"productionTasks":[` + link("productiontask") + `],"invoicesOut":[` + link("invoiceout") + `],"moves":[` + link("move") + `],"prepayments":[` + link("prepayment") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := customerorder.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(customerorder.UpdateRequest))
				return err
			},
			func() any {
				return customerorder.CreateRequest{PurchaseOrders: model.Value([]model.Reference{ref("purchaseorder")})}
			}, `"purchaseOrders":`},
		{"demand", demand.UpdateRequest{FactureOut: model.Value(ref("factureout")), Returns: model.Value([]model.Reference{ref("salesreturn")}), Payments: model.Value([]model.Reference{ref("paymentin")}), ProductionTasks: model.Value([]model.Reference{ref("productiontask")}), InvoicesOut: model.Value([]model.Reference{ref("invoiceout")})}, `{"factureOut":` + link("factureout") + `,"returns":[` + link("salesreturn") + `],"payments":[` + link("paymentin") + `],"productionTasks":[` + link("productiontask") + `],"invoicesOut":[` + link("invoiceout") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := demand.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(demand.UpdateRequest))
				return err
			},
			func() any { return demand.CreateRequest{FactureOut: model.Value(ref("factureout"))} }, `"factureOut":`},
		{"enter", enter.UpdateRequest{Inventory: model.Value(ref("inventory"))}, `{"inventory":` + link("inventory") + `}`,
			func(c *transport.Client, req any) error {
				s, _ := enter.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(enter.UpdateRequest))
				return err
			},
			func() any { return enter.CreateRequest{Inventory: model.Value(ref("inventory"))} }, `"inventory":`},
		{"invoicein", invoicein.UpdateRequest{Supplies: model.Value([]model.Reference{ref("supply")})}, `{"supplies":[` + link("supply") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := invoicein.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(invoicein.UpdateRequest))
				return err
			},
			func() any { return invoicein.CreateRequest{Supplies: model.Value([]model.Reference{ref("supply")})} }, `"supplies":`},
		{"invoiceout", invoiceout.UpdateRequest{Demands: model.Value([]model.Reference{ref("demand")})}, `{"demands":[` + link("demand") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := invoiceout.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(invoiceout.UpdateRequest))
				return err
			},
			func() any { return invoiceout.CreateRequest{Demands: model.Value([]model.Reference{ref("demand")})} }, `"demands":`},
		{"loss", loss.UpdateRequest{Inventory: model.Value(ref("inventory"))}, `{"inventory":` + link("inventory") + `}`,
			func(c *transport.Client, req any) error {
				s, _ := loss.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(loss.UpdateRequest))
				return err
			},
			func() any { return loss.CreateRequest{Inventory: model.Value(ref("inventory"))} }, `"inventory":`},
		{"move", move.UpdateRequest{ProductionTasks: model.Value([]model.Reference{ref("productiontask")})}, `{"productionTasks":[` + link("productiontask") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := move.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(move.UpdateRequest))
				return err
			},
			func() any {
				return move.CreateRequest{ProductionTasks: model.Value([]model.Reference{ref("productiontask")})}
			}, `"productionTasks":`},
		{"purchaseorder", purchaseorder.UpdateRequest{CustomerOrders: model.Value([]model.Reference{ref("customerorder")}), InvoicesIn: model.Value([]model.Reference{ref("invoicein")}), ProductionTasks: model.Value([]model.Reference{ref("productiontask")}), Payments: model.Value([]model.Reference{ref("paymentin")}), Supplies: model.Value([]model.Reference{ref("supply")})}, `{"customerOrders":[` + link("customerorder") + `],"invoicesIn":[` + link("invoicein") + `],"productionTasks":[` + link("productiontask") + `],"payments":[` + link("paymentin") + `],"supplies":[` + link("supply") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := purchaseorder.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(purchaseorder.UpdateRequest))
				return err
			},
			func() any {
				return purchaseorder.CreateRequest{CustomerOrders: model.Value([]model.Reference{ref("customerorder")})}
			}, `"customerOrders":`},
		{"purchasereturn", purchasereturn.UpdateRequest{FactureOut: model.Value(ref("factureout")), FactureIn: model.Value(ref("facturein")), Payments: model.Value([]model.Reference{ref("paymentin")})}, `{"factureOut":` + link("factureout") + `,"factureIn":` + link("facturein") + `,"payments":[` + link("paymentin") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := purchasereturn.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(purchasereturn.UpdateRequest))
				return err
			},
			func() any { return purchasereturn.CreateRequest{FactureOut: model.Value(ref("factureout"))} }, `"factureOut":`},
		{"salesreturn", salesreturn.UpdateRequest{Losses: model.Value([]model.Reference{ref("loss")}), Payments: model.Value([]model.Reference{ref("paymentin")}), FactureOut: model.Value(ref("factureout"))}, `{"losses":[` + link("loss") + `],"payments":[` + link("paymentin") + `],"factureOut":` + link("factureout") + `}`,
			func(c *transport.Client, req any) error {
				s, _ := salesreturn.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(salesreturn.UpdateRequest))
				return err
			},
			func() any { return salesreturn.CreateRequest{Losses: model.Value([]model.Reference{ref("loss")})} }, `"losses":`},
		{"supply", supply.UpdateRequest{FactureIn: model.Value(ref("facturein")), InvoicesIn: model.Value([]model.Reference{ref("invoicein")}), Payments: model.Value([]model.Reference{ref("paymentin")}), ProductionTask: model.Value(ref("productiontask")), Returns: model.Value([]model.Reference{ref("salesreturn")})}, `{"factureIn":` + link("facturein") + `,"invoicesIn":[` + link("invoicein") + `],"payments":[` + link("paymentin") + `],"productionTask":` + link("productiontask") + `,"returns":[` + link("salesreturn") + `]}`,
			func(c *transport.Client, req any) error {
				s, _ := supply.NewService(c)
				_, err := s.Update(t.Context(), docID, req.(supply.UpdateRequest))
				return err
			},
			func() any { return supply.CreateRequest{FactureIn: model.Value(ref("facturein"))} }, `"factureIn":`},
	} {
		t.Run(tc.document, func(t *testing.T) {
			for _, codec := range documenttest.Codecs {
				data, err := codec.Marshal(tc.update)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, data, []byte(tc.body))
				// Create пишет то же значение связи, что Update: сравниваем поле целиком.
				if data, err = codec.Marshal(tc.create()); err != nil {
					t.Fatal(err)
				}
				var created, updated map[string]jsontext.Value
				name := strings.Trim(tc.key, `":`)
				if err := json.Unmarshal(data, &created); err != nil || json.Unmarshal([]byte(tc.body), &updated) != nil || created[name] == nil {
					t.Fatalf("%s: create = %s, err = %v", codec.Name, data, err)
				}
				documenttest.JSONEqual(t, created[name], updated[name])
			}
			for _, c := range []struct {
				name, response string
				status         int
				ok             bool
			}{
				{"ok", `{"id":"` + docID + `",` + tc.body[1:], 200, true},
				{"api error", `{"errors":[{"error":"bad link","code":3000}]}`, 412, false},
			} {
				var calls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != http.MethodPut || r.URL.Path != "/entity/"+tc.document+"/"+docID {
						t.Errorf("request = %s %s", r.Method, r.URL.Path)
					}
					documenttest.JSONEqual(t, b, []byte(tc.body))
					w.WriteHeader(c.status)
					io.WriteString(w, c.response)
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.invoke(client, tc.update)
				server.Close()
				if (err == nil) != c.ok || calls.Load() != 1 || (err != nil && strings.Count(err.Error(), "moysklad:") != 1) {
					t.Fatalf("%s: err = %v, calls = %d", c.name, err, calls.Load())
				}
			}
		})
	}
	// Пустой массив отправляется как [] (что делает сервер, источник не описывает), null — как null.
	for _, codec := range documenttest.Codecs {
		data, err := codec.Marshal(demand.UpdateRequest{Payments: model.Value([]model.Reference{}), FactureOut: model.Null[model.Reference]()})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`{"payments":[],"factureOut":null}`))
	}
}
