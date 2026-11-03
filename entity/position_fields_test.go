// Проверка O06: поля строк invoicein, invoiceout и move по таблицам позиций источника
// (pack; у move — things, sourceSlot, targetSlot) уходят в UpdatePosition и читаются из ответа.
package entity_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestDocumentPositionFields(t *testing.T) {
	const (
		docID = "7944ef04-f831-11e5-7a69-971500188b19"
		posID = "7944ef04-f831-11e5-7a69-971500188b21"
		pack  = `"pack":{"id":"7944ef04-f831-11e5-7a69-971500188b22"}`
	)
	slot := func(n string) model.Reference {
		return model.Reference{Meta: model.Meta{Href: linkBase + "store/" + docID + "/slots/" + n, Type: "slot"}}
	}
	withPack := positions.PositionRequest{Pack: model.Value(positions.PackRequest{ID: "7944ef04-f831-11e5-7a69-971500188b22"})}
	moveRow := withPack
	moveRow.Things = model.Value([]string{"S-1", "S-2"})
	moveRow.SourceSlot = model.Value(slot("a"))
	moveRow.TargetSlot = model.Value(slot("b"))
	slotJSON := func(n string) string {
		return `{"meta":{"href":"` + linkBase + "store/" + docID + "/slots/" + n + `","type":"slot"}}`
	}
	for _, tc := range []struct {
		document, body string
		update         func(*transport.Client) (*positions.Position, error)
	}{
		{"invoicein", `{` + pack + `}`, func(c *transport.Client) (*positions.Position, error) {
			s, _ := invoicein.NewService(c)
			return s.UpdatePosition(t.Context(), docID, posID, withPack)
		}},
		{"invoiceout", `{` + pack + `}`, func(c *transport.Client) (*positions.Position, error) {
			s, _ := invoiceout.NewService(c)
			return s.UpdatePosition(t.Context(), docID, posID, withPack)
		}},
		{"move", `{` + pack + `,"things":["S-1","S-2"],"sourceSlot":` + slotJSON("a") + `,"targetSlot":` + slotJSON("b") + `}`,
			func(c *transport.Client) (*positions.Position, error) {
				s, _ := move.NewService(c)
				return s.UpdatePosition(t.Context(), docID, posID, moveRow)
			}},
	} {
		t.Run(tc.document, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodPut || r.URL.Path != "/entity/"+tc.document+"/"+docID+"/positions/"+posID {
					t.Errorf("request = %s %s", r.Method, r.URL.Path)
				}
				documenttest.JSONEqual(t, b, []byte(tc.body))
				io.WriteString(w, `{"id":"`+posID+`",`+tc.body[1:])
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			got, err := tc.update(c)
			if err != nil || got.Pack == nil || got.Pack.ID == "" || (tc.document == "move" && (len(got.Things) != 2 || got.TargetSlot == nil)) {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
		})
	}
}
