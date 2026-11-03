package contract_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/contract"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "3bba8d2a-2cb6-11e6-8a84-bae50000002c"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ), MediaType: "application/json"}}
}

var (
	full   = contract.CreateRequest{Name: model.Value("12345"), Description: model.Value(""), Code: model.Value("c"), ExternalCode: model.Value("ext"), Archived: model.Value(false), Moment: model.Value(model.DateTime(time.Date(2016, 7, 6, 12, 53, 22, 0, time.UTC))), Sum: model.Value(0.0), ContractType: model.Value("Commission"), RewardType: model.Value("PercentOfSales"), RewardPercent: model.Value(12.5), OwnAgent: ref("entity/organization/fae3561a-2e58-11e6-8a84-bae50000004e", "organization"), Agent: ref("entity/counterparty/12a8b923-2e58-11e6-8a84-bae500000057", "counterparty"), OrganizationAccount: model.Null[model.Reference](), Rate: model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/baac25f0-50ac-11e5-300d-c79b00000055", "currency"))}), Attributes: model.Value([]model.Attribute{}), Shared: model.Value(true), Owner: model.Value(ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee"))}
	update = contract.UpdateRequest{Sum: model.Value(1050.5), State: model.Null[model.Reference](), RewardType: model.Value("None")}
)

func TestContractFixture(t *testing.T) {
	data := read(t, "upstream/contract.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got contract.Contract
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "00000011111" || got.Sum != 200000 || got.ContractType != "Sales" || got.RewardType != "PercentOfSales" || got.RewardPercent != 10 || got.OwnAgent == nil || got.Agent == nil || got.AgentAccount == nil || got.OrganizationAccount == nil || got.State == nil || got.Rate == nil || got.Rate.Currency.Meta.Type != "currency" || got.Moment.IsZero() || got.Printed || got.Published {
				t.Fatalf("contract = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse contract.Contract
			if err := codec.Unmarshal([]byte(`{"id":"x","futureField":{"a":1}}`), &sparse); err != nil || sparse.ID != "x" {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/contract/create-full.golden.json")},
			{"update", update, read(t, "custom/contract/update.golden.json")},
			{"create minimal", contract.CreateRequest{OwnAgent: ref("entity/organization/1", "organization"), Agent: ref("entity/counterparty/2", "counterparty")}, []byte(`{"ownAgent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/1","type":"organization","mediaType":"application/json"}},"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/2","type":"counterparty","mediaType":"application/json"}}}`)},
			{"update sum 0", contract.UpdateRequest{Sum: model.Value(0.0), RewardPercent: model.Value(0.0)}, []byte(`{"sum":0,"rewardPercent":0}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[contract.CreateRequest](), reflect.TypeFor[contract.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Printed", "Published"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *contract.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := contract.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/contract/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*contract.Contract, error)
	}{
		{"get", "GET", "entity/contract/" + testID, nil, func(ctx context.Context, c *transport.Client) (*contract.Contract, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/contract", read(t, "custom/contract/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*contract.Contract, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/contract/" + testID, read(t, "custom/contract/update.golden.json"), func(ctx context.Context, c *transport.Client) (*contract.Contract, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *contract.Contract) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("contract = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no own agent": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, contract.CreateRequest{Agent: ref("entity/counterparty/2", "counterparty")}, o...)
			return err
		},
		"create no agent": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, contract.CreateRequest{OwnAgent: ref("entity/organization/1", "organization")}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, contract.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "contract", read(t, "upstream/contract.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[contract.Contract] {
		return service(c).List
	}, func(t *testing.T, got contract.Contract) {
		if got.ID != testID || got.Code != "contractWithPartner" {
			t.Fatalf("row = %+v", got)
		}
	})
}
