package currency_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/currency"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "22222222-2222-2222-2222-222222222222"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

var (
	unit = func(g, s1, s2, s5 string) currency.Unit { return currency.Unit{Gender: g, S1: s1, S2: s2, S5: s5} }
	full = currency.CreateRequest{
		Name: model.Value("доллар"), FullName: model.Value("Доллар США"), Code: model.Value("840"), ISOCode: model.Value("USD"),
		Archived: model.Value(false), Rate: model.Value(63.25), Multiplicity: model.Value(1), Indirect: model.Value(false),
		MajorUnit: model.Value(unit("masculine", "доллар", "доллара", "долларов")),
		MinorUnit: model.Value(unit("masculine", "цент", "цента", "центов")),
	}
	system = currency.CreateRequest{System: model.Value(true), ISOCode: model.Value("EUR"), RateUpdateType: model.Value("auto"), Margin: model.Value(2.5)}
	update = currency.UpdateRequest{Name: model.Value("долл"), Rate: model.Value(66.0), Archived: model.Value(true),
		Multiplicity: model.Value(100), Indirect: model.Value(true), Margin: model.Value(1.5)}
)

func TestCurrencyFixture(t *testing.T) {
	data := read(t, "upstream/currency.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got currency.Currency
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.ISOCode != "RUB" || got.Code != "643" || !got.Default || !got.System || got.Rate != 1 ||
				got.Multiplicity != 1 || got.RateUpdateType != "manual" || got.MajorUnit == nil || got.MinorUnit.S5 != "копеек" {
				t.Fatalf("currency = %+v", got)
			}
		})
	}
}

// Пользовательская валюта и системная с автокурсом — примеры ответов источника (README рядом).
func TestCurrencyExamples(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var user, auto currency.Currency
		if err := codec.Unmarshal(read(t, "custom/currency/user.json"), &user); err != nil || user.System || user.Default ||
			user.ISOCode != "USD" || user.RateUpdateType != "manual" || user.Rate != 63 {
			t.Fatalf("%s: user = %+v, err = %v", codec.Name, user, err)
		}
		if err := codec.Unmarshal(read(t, "custom/currency/system-auto.json"), &auto); err != nil || !auto.System || auto.RateUpdateType != "auto" {
			t.Fatalf("%s: system = %+v, err = %v", codec.Name, auto, err)
		}
		var sparse currency.Currency
		if err := codec.Unmarshal([]byte(`{"id":"x","isoCode":"XXX","futureField":{"a":1}}`), &sparse); err != nil ||
			sparse.MajorUnit != nil || sparse.Rate != 0 || sparse.RateUpdateType != "" {
			t.Fatalf("%s: sparse = %+v, err = %v", codec.Name, sparse, err)
		}
	}
}

// Только для чтения: id, meta, default — ни в одном запросе; system и rateUpdateType — только в Create.
func TestRequestReadOnlyFields(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[currency.CreateRequest](), reflect.TypeFor[currency.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "Default"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"System", "RateUpdateType"} {
		if _, ok := reflect.TypeFor[currency.UpdateRequest]().FieldByName(name); ok {
			t.Errorf("UpdateRequest contains create-only %s", name)
		}
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/currency/create-full.golden.json")},
			{"create system", system, read(t, "custom/currency/create-system.golden.json")},
			{"update", update, read(t, "custom/currency/update.golden.json")},
			{"update archived false", currency.UpdateRequest{Archived: model.Value(false)}, []byte(`{"archived":false}`)},
			{"update nulls", currency.UpdateRequest{FullName: model.Null[string](), MajorUnit: model.Null[currency.Unit](), Margin: model.Null[float64]()},
				[]byte(`{"fullName":null,"majorUnit":null,"margin":null}`)},
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
}

func service(c *transport.Client) *currency.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := currency.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/currency/` + testID + `"},"id":"` + testID + `","isoCode":"USD","rate":63.25}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*currency.Currency, error)
	}{
		{"get", "GET", "entity/currency/" + testID, nil, func(ctx context.Context, c *transport.Client) (*currency.Currency, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/currency", read(t, "custom/currency/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*currency.Currency, error) {
			return service(c).Create(ctx, full)
		}},
		{"create system", "POST", "entity/currency", read(t, "custom/currency/create-system.golden.json"), func(ctx context.Context, c *transport.Client) (*currency.Currency, error) {
			return service(c).Create(ctx, system)
		}},
		{"update", "PUT", "entity/currency/" + testID, read(t, "custom/currency/update.golden.json"), func(ctx context.Context, c *transport.Client) (*currency.Currency, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *currency.Currency) {
						if got.ID != testID || got.Rate != 63.25 {
							t.Fatalf("currency = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Правила источника: обязательные поля обеих форм Create, rate ≠ 0, непустые name/code/isoCode.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(client)
	zeroRate := full
	zeroRate.Rate = model.Value(0.0)
	nullRate := full
	nullRate.Rate = model.Null[float64]()
	for name, input := range map[string]currency.CreateRequest{
		"empty":            {},
		"no isoCode":       {Name: model.Value("x"), Code: model.Value("1")},
		"empty name":       {Name: model.Value(""), Code: model.Value("1"), ISOCode: model.Value("X")},
		"zero rate":        zeroRate,
		"null rate":        nullRate,
		"system no codes":  {System: model.Value(true), Name: model.Value("евро")},
		"system null code": {System: model.Value(true), Code: model.Null[string]()},
	} {
		if _, err := s.Create(t.Context(), input); err == nil {
			t.Errorf("create %s accepted", name)
		}
	}
	for name, input := range map[string]currency.UpdateRequest{
		"empty":      {},
		"zero rate":  {Rate: model.Value(0.0)},
		"empty code": {Code: model.Value("")},
		"null name":  {Name: model.Null[string]()},
	} {
		if _, err := s.Update(t.Context(), testID, input); err == nil {
			t.Errorf("update %s accepted", name)
		}
	}
	if _, err := s.Get(t.Context(), "bad"); err == nil {
		t.Error("bad id accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "currency", read(t, "upstream/currency.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[currency.Currency] {
		return service(c).List
	}, func(t *testing.T, got currency.Currency) {
		if got.ID != testID || got.ISOCode != "RUB" {
			t.Fatalf("row = %+v", got)
		}
	})
}
