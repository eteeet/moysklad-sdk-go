package retailstore_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retailstore"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "bc9a4c67-2424-11f1-0a82-04ce00000865"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *retailstore.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := retailstore.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	store        = ref("entity/store/33333333-3333-3333-3333-333333333333", "store")
	priceType    = model.PriceType{Meta: ref("context/companysettings/pricetype/44444444-4444-4444-4444-444444444444", "pricetype").Meta}
	full         = retailstore.CreateRequest{
		Name: "Касса 1", Organization: organization, Store: store, PriceType: priceType,
		Address: model.Value(""), Active: model.Value(false), AllowCreateProducts: model.Value(true), ControlShippingStock: model.Value(false),
		RequiredPhone: model.Value(false), SyncAgents: model.Value(true),
		DefaultTaxSystem: model.Value("GENERAL_TAX_SYSTEM"), FiscalType: model.Value("STANDARD"),
		BankPercent: model.Value(1.5), QRBankPercent: model.Value(0.0), DiscountMaxPercent: model.Value(15.0),
		Acquire:             model.Value(ref("entity/counterparty/55555555-5555-5555-5555-555555555555", "counterparty")),
		AddressFull:         model.Value(model.AddressRequest{City: model.Value("Москва"), Street: model.Value("Тверская")}),
		CreateAgentsTags:    model.Value([]string{"розница"}),
		FilterAgentsTags:    model.Value([]string{}),
		CustomerOrderStates: model.Value([]model.Reference{ref("entity/customerorder/metadata/states/66666666-6666-6666-6666-666666666666", "state")}),
		MasterRetailStores:  model.Value([]model.Reference{ref("entity/retailstore/77777777-7777-7777-7777-777777777777", "retailstore")}),
		ProductFolders:      model.Value([]model.Reference{ref("entity/productfolder/88888888-8888-8888-8888-888888888888", "productfolder")}),
		Description:         model.Value("касса у входа"), DemandPrefix: model.Value("R-"), IDQR: model.Value("qr-1"), QRTerminalID: model.Value("T-1"),
		OrderTaxSystem: model.Value("SIMPLIFIED_TAX_SYSTEM_INCOME"), MarkingSellingMode: model.Value("CORRECT_MARKS_ONLY"),
		MarksCheckMode: model.Value("CORRECT_MARKS_ONLY"), MinionToMasterType: model.Value("ANY"), PriorityOFDSend: model.Value("PHONE"),
		TobaccoMRCControlType: model.Value("USER_PRICE"), QRPayEnabled: model.Value(true), SellReserves: model.Value(false),
		CreateOrderWithState: model.Value(ref("entity/customerorder/metadata/states/99999999-9999-9999-9999-999999999999", "state")),
		ReceiptTemplate:      model.Null[model.Reference](),
	}
	update = retailstore.UpdateRequest{
		Name: model.Value("Касса 2"), Archived: model.Value(true), Acquire: model.Null[model.Reference](),
		ProductFolders: model.Value([]model.Reference{}), MinionToMasterType: model.Value("CHOSEN"), DiscountMaxPercent: model.Value(0.0),
	}
)

func TestRetailStoreFixture(t *testing.T) {
	data := read(t, "upstream/retail_store.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got retailstore.RetailStore
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Store == nil || got.PriceType == nil || got.PriceType.Name == "" ||
				len(got.CreateAgentsTags) != 1 || got.Cashiers == nil || got.Cashiers.Meta.Size != 1 || got.MasterRetailStores == nil ||
				len(got.CustomerOrderStates) != 1 || len(got.LastOperationNames) != 5 || got.AddressFull == nil ||
				len(got.Environment) == 0 || len(got.State) == 0 || got.DefaultTaxSystem == "" || got.DiscountMaxPercent != 100 {
				t.Fatalf("store = %+v", got)
			}
			var env struct {
				Software struct {
					Name string `json:"name"`
				} `json:"software"`
			}
			if err := codec.Unmarshal(got.Environment, &env); err != nil || env.Software.Name != "test" {
				t.Fatalf("environment = %s, err = %v", got.Environment, err)
			}
			var sparse retailstore.RetailStore
			if err := codec.Unmarshal([]byte(`{"id":"x","fiscalType":"FUTURE","future":1}`), &sparse); err != nil ||
				sparse.FiscalType != "FUTURE" || sparse.Environment != nil || sparse.PriceType != nil || sparse.Cashiers != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

// masterRetailStores и productFolders: фикстура — объект коллекции, пример ответа Create — массив.
func TestRetailStoreLinks(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, input := range []string{
			`{"meta":{"href":"h","size":1,"limit":1000,"offset":0}}`,
			`[{"meta":{"href":"h","type":"retailstore"}}]`,
			`[]`,
		} {
			var got retailstore.RetailStore
			data := []byte(`{"id":"x","masterRetailStores":` + input + `}`)
			if err := codec.Unmarshal(data, &got); err != nil || got.MasterRetailStores == nil {
				t.Fatalf("%s %s: got = %+v, err = %v", codec.Name, input, got.MasterRetailStores, err)
			}
			out, err := codec.Marshal(got.MasterRetailStores)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, out, []byte(input))
		}
		// Форма примера: массив из метаданных коллекции; size в ссылке не хранится, ошибки нет.
		var example retailstore.RetailStore
		if err := codec.Unmarshal([]byte(`{"id":"x","productFolders":[{"meta":{"href":"h","size":1}}]}`), &example); err != nil ||
			len(example.ProductFolders.Rows) != 1 {
			t.Fatalf("%s: example = %+v, err = %v", codec.Name, example.ProductFolders, err)
		}
	}
}

func TestRetailStoreRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/retailstore/create-full.golden.json")},
			{"create minimal", retailstore.CreateRequest{Name: "К", Organization: organization, Store: store, PriceType: priceType},
				[]byte(`{"name":"К","organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}},"priceType":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/44444444-4444-4444-4444-444444444444","type":"pricetype"}}}`)},
			{"update", update, read(t, "custom/retailstore/update.golden.json")},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, tc.want)
			})
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[retailstore.CreateRequest](), reflect.TypeFor[retailstore.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "ExternalCode", "AllowCustomPrice", "AuthTokenAttached", "OFDEnabled",
			"Environment", "State", "LastOperationNames", "Cashiers"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or out-of-scope %s", typ.Name(), name)
			}
		}
	}
}

func TestRetailStoreHTTP(t *testing.T) {
	const response = `{"id":"server-id","name":"server-name","active":true}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*retailstore.RetailStore, error)
	}{
		{"get", "GET", "entity/retailstore/" + testID, nil, func(ctx context.Context, c *transport.Client) (*retailstore.RetailStore, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/retailstore", read(t, "custom/retailstore/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*retailstore.RetailStore, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/retailstore/" + testID, read(t, "custom/retailstore/update.golden.json"), func(ctx context.Context, c *transport.Client) (*retailstore.RetailStore, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *retailstore.RetailStore) {
						if got.ID != "server-id" || got.Name != "server-name" || !got.Active {
							t.Fatalf("store = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Пустые обязательные поля отклоняются до запроса; конфликт allowCreateProducts и
// controlShippingStock проверяет сервер — запрос уходит как есть.
func TestRetailStoreValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"id":"server-id"}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	conflict := retailstore.CreateRequest{Name: "К", Organization: organization, Store: store, PriceType: priceType,
		AllowCreateProducts: model.Value(true), ControlShippingStock: model.Value(true)}
	if _, err := s.Create(t.Context(), conflict); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"no name": func() error {
			_, err := s.Create(t.Context(), retailstore.CreateRequest{Organization: organization, Store: store, PriceType: priceType})
			return err
		}(),
		"no organization": func() error {
			_, err := s.Create(t.Context(), retailstore.CreateRequest{Name: "К", Store: store, PriceType: priceType})
			return err
		}(),
		"no store": func() error {
			_, err := s.Create(t.Context(), retailstore.CreateRequest{Name: "К", Organization: organization, PriceType: priceType})
			return err
		}(),
		"no price type": func() error {
			_, err := s.Create(t.Context(), retailstore.CreateRequest{Name: "К", Organization: organization, Store: store})
			return err
		}(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, retailstore.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", update); return err }(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retailstore") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := retailstore.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestRetailStoreList(t *testing.T) {
	listtest.Run(t, "retailstore", read(t, "upstream/retail_store.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[retailstore.RetailStore] {
		return service(c).List
	}, func(t *testing.T, got retailstore.RetailStore) {
		if got.ID != fixtureID || got.PriceType == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
