package demand_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	carrierHref = "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/cffd2c81-62b4-11e6-8a84-bae500000053"
	branchHref  = "https://api.moysklad.ru/api/remap/1.2/entity/organization/869d7628-6396-11e6-8a84-bae50000000c/branch/1"
	// Транспортные поля примера «Отгрузка с полями ТН» источника (ссылки сокращены до meta).
	transportJSON = `{"cargoName":"Овощи","carrier":{"meta":{"href":"` + carrierHref + `","type":"counterparty"}},` +
		`"consignee":{"meta":{"href":"` + carrierHref + `","type":"counterparty"}},"consigneeBranch":{"meta":{"href":"` + branchHref + `","type":"organizationbranch"}},` +
		`"consignor":{"meta":{"href":"` + carrierHref + `","type":"counterparty"}},"consignorBranch":{"meta":{"href":"` + branchHref + `","type":"organizationbranch"}},` +
		`"goodPackQuantity":10,"shippingInstructions":"Беречь от воздействия солнца","stateContractId":"12345",` +
		`"transportFacility":"Mercedes benz","transportFacilityNumber":"лb777m","shippingDocumentNumber":"76543210",` +
		`"shippingDocumentDate":"2026-04-15 00:00:00.000"}`
)

func TestTransportFields(t *testing.T) {
	carrier := model.Reference{Meta: model.Meta{Href: carrierHref, Type: "counterparty"}}
	branch := model.Reference{Meta: model.Meta{Href: branchHref, Type: "organizationbranch"}}
	date := model.DateTime(time.Date(2026, 4, 14, 21, 0, 0, 0, time.UTC))
	update := demand.UpdateRequest{
		CargoName: model.Value("Овощи"), Carrier: model.Value(carrier), Consignee: model.Value(carrier),
		ConsigneeBranch: model.Value(branch), Consignor: model.Value(carrier), ConsignorBranch: model.Value(branch),
		GoodPackQuantity: model.Value(10), ShippingInstructions: model.Value("Беречь от воздействия солнца"),
		StateContractID: model.Value("12345"), TransportFacility: model.Value("Mercedes benz"),
		TransportFacilityNumber: model.Value("лb777m"), ShippingDocumentNumber: model.Value("76543210"),
		ShippingDocumentDate: model.Value(date),
	}
	for _, codec := range documenttest.Codecs {
		var got demand.Demand
		if err := codec.Unmarshal([]byte(transportJSON), &got); err != nil {
			t.Fatal(err)
		}
		if got.CargoName != "Овощи" || got.Carrier == nil || got.ConsigneeBranch == nil || got.ConsignorBranch.Meta.Type != "organizationbranch" ||
			got.GoodPackQuantity != 10 || got.StateContractID != "12345" || got.TransportFacilityNumber != "лb777m" ||
			!time.Time(got.ShippingDocumentDate).Equal(time.Time(date)) {
			t.Fatalf("%s: demand = %+v", codec.Name, got)
		}
		data, err := codec.Marshal(update)
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(transportJSON))
		// Очистка: пустая строка и null уходят как есть.
		data, err = codec.Marshal(demand.UpdateRequest{CargoName: model.Value(""), Carrier: model.Null[model.Reference](), GoodPackQuantity: model.Value(0)})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`{"cargoName":"","carrier":null,"goodPackQuantity":0}`))
	}

	const id = "0ee24723-f640-11e5-8a84-bae500000065"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		if r.Method != http.MethodPut || r.URL.Path != "/entity/demand/"+id {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		documenttest.JSONEqual(t, b, []byte(transportJSON))
		io.WriteString(w, `{"id":"`+id+`",`+transportJSON[1:])
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := demand.NewService(c)
	got, err := s.Update(t.Context(), id, update)
	if err != nil || got.TransportFacility != "Mercedes benz" || got.Consignor == nil {
		t.Fatalf("got = %+v, err = %v", got, err)
	}
}
