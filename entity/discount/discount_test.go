package discount_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const roundOffID = "8ae26646-b1aa-11ea-ac12-000b00000001"

func read(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/upstream/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *discount.Service {
	s, _ := discount.NewService(c)
	return s
}

// Общий тип читает все виды скидок; поля вида приходят только у своих.
func TestFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var round, acc, personal, special discount.Discount
		for name, dst := range map[string]*discount.Discount{"discount.json": &round, "accumulation_discount.json": &acc,
			"personal_discount.json": &personal, "special_price_discount.json": &special} {
			if err := codec.Unmarshal(read(t, name), dst); err != nil {
				t.Fatal(err)
			}
		}
		if round.ID != roundOffID || round.Meta.Type != "discount" || !round.Active || !round.AllAgents || round.AgentTags == nil {
			t.Fatalf("%s: round = %+v", codec.Name, round)
		}
		if acc.Meta.Type != "accumulationdiscount" || len(acc.Levels) != 2 || acc.Levels[1] != (discount.Level{Amount: 200000, Discount: 15}) ||
			len(acc.Assortment) == 0 || len(acc.ProductFolders) == 0 || acc.AgentTags[0] != "группа агентов" {
			t.Fatalf("%s: accumulation = %+v", codec.Name, acc)
		}
		if personal.Meta.Type != "personaldiscount" || !personal.AllProducts || personal.Levels != nil || personal.SpecialPrice != nil {
			t.Fatalf("%s: personal = %+v", codec.Name, personal)
		}
		if !strings.Contains(string(acc.Raw), `"levels"`) {
			t.Fatalf("%s: raw = %s", codec.Name, acc.Raw)
		}
		var future discount.Discount
		if err := codec.Unmarshal([]byte(`{"meta":{"type":"futurediscount"},"id":"x","productfolders":[{"meta":{"type":"productfolder"}}],"bonus":5}`), &future); err != nil ||
			future.Meta.Type != "futurediscount" || len(future.ProductFolders) != 1 || !strings.Contains(string(future.Raw), `"bonus":5`) {
			t.Fatalf("%s: future = %+v, err = %v", codec.Name, future, err)
		}
		if special.Meta.Type != "specialpricediscount" || !special.UsePriceType || special.SpecialPrice == nil || special.SpecialPrice.Value != 15 ||
			special.SpecialPrice.PriceType == nil || special.SpecialPrice.PriceType.Name != "Цена продажи" {
			t.Fatalf("%s: special = %+v", codec.Name, special)
		}
	}
}

func TestUpdateRoundOff(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*discount.Discount, error) {
		return service(c).UpdateRoundOff(ctx, roundOffID, discount.RoundOffRequest{Name: model.Value("updatedName"), Active: model.Value(true)})
	}
	for _, tc := range documenttest.ResponseCases(string(read(t, "discount.json"))) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "PUT", "entity/discount/"+roundOffID, []byte(`{"name":"updatedName","active":true}`), invoke, func(got *discount.Discount) {
				if got.ID != roundOffID {
					t.Fatalf("discount = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestUpdateMany(t *testing.T) {
	var body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = r.Method + " " + r.URL.Path + " " + string(b)
		w.Write([]byte(`[` + string(read(t, "discount.json")) + `]`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	r, err := service(c).UpdateMany(t.Context(), []batch.Update[discount.RoundOffRequest]{{ID: roundOffID, Request: discount.RoundOffRequest{Active: model.Value(false)}}})
	if err != nil || len(r) != 1 || r[0].Value.ID != roundOffID {
		t.Fatalf("r = %+v, err = %v", r, err)
	}
	if !strings.HasPrefix(body, `POST /entity/discount [{"meta":{"href":"`+server.URL+`/entity/discount/`+roundOffID+`","metadataHref":"`+server.URL+`/entity/discount/metadata","type":"discount"`) ||
		!strings.HasSuffix(body, `"active":false}]`) {
		t.Fatalf("body = %s", body)
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"bad id": func(o ...transport.RequestOption) error {
			_, err := service(c).UpdateRoundOff(t.Context(), "x", discount.RoundOffRequest{Active: model.Value(true)}, o...)
			return err
		},
		"empty": func(o ...transport.RequestOption) error {
			_, err := service(c).UpdateRoundOff(t.Context(), roundOffID, discount.RoundOffRequest{}, o...)
			return err
		},
		"many none": func(o ...transport.RequestOption) error {
			_, err := service(c).UpdateMany(t.Context(), nil, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "discount", read(t, "accumulation_discount.json"), "", "", func(c *transport.Client) listtest.ListFunc[discount.Discount] {
		return service(c).List
	}, func(t *testing.T, got discount.Discount) {
		if got.Meta.Type != "accumulationdiscount" || len(got.Levels) != 2 {
			t.Fatalf("row = %+v", got)
		}
	})
}
