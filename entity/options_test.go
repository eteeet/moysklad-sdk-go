// Проверка T02: каждый service передаёт per-call опции в transport и заполняет ResponseInfo.
package entity_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "77777777-7777-7777-7777-777777777777"

type optionsCase struct {
	name, method, path string
	invoke             func(context.Context, *transport.Client, ...transport.RequestOption) error
}

func getCase[S, T any](name string, newService func(*transport.Client) (*S, error),
	get func(*S, context.Context, string, ...transport.RequestOption) (*T, error)) optionsCase {
	return optionsCase{name + " get", http.MethodGet, "entity/" + name + "/" + testID,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := newService(c)
			if err != nil {
				return err
			}
			_, err = get(s, ctx, testID, opts...)
			return err
		}}
}

func positionsCase[S, P any](name string, newService func(*transport.Client) (*S, error),
	get func(*S, context.Context, string, positions.PageOptions, ...transport.RequestOption) (*model.Collection[P], error)) optionsCase {
	return optionsCase{name + " positions", http.MethodGet, "entity/" + name + "/" + testID + "/positions?limit=1000&offset=0",
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := newService(c)
			if err != nil {
				return err
			}
			_, err = get(s, ctx, testID, positions.PageOptions{}, opts...)
			return err
		}}
}

func createCase[S, I, T any](name string, newService func(*transport.Client) (*S, error),
	create func(*S, context.Context, I, ...transport.RequestOption) (*T, error)) optionsCase {
	return optionsCase{name + " create", http.MethodPost, "entity/" + name,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := newService(c)
			if err != nil {
				return err
			}
			_, err = create(s, ctx, fillRequest[I](), opts...)
			return err
		}}
}

func updateCase[S, I, T any](name string, newService func(*transport.Client) (*S, error),
	update func(*S, context.Context, string, I, ...transport.RequestOption) (*T, error)) optionsCase {
	return optionsCase{name + " update", http.MethodPut, "entity/" + name + "/" + testID,
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := newService(c)
			if err != nil {
				return err
			}
			_, err = update(s, ctx, testID, fillRequest[I](), opts...)
			return err
		}}
}

// fillRequest заполняет обязательные для локальной проверки поля: Name, ссылки и Description.
func fillRequest[I any]() I {
	var input I
	v := reflect.ValueOf(&input).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		switch {
		case v.Type().Field(i).Name == "Name" && field.Kind() == reflect.String:
			field.SetString("x")
		case field.Type() == reflect.TypeFor[model.Reference]():
			field.Set(reflect.ValueOf(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/x"}}))
		case v.Type().Field(i).Name == "Description" && field.Type() == reflect.TypeFor[model.Field[string]]():
			field.Set(reflect.ValueOf(model.Value("x")))
		}
	}
	return input
}

func TestServicesAcceptRequestOptions(t *testing.T) {
	for _, tc := range []optionsCase{
		getCase("cashin", cashin.NewService, (*cashin.Service).Get),
		getCase("cashout", cashout.NewService, (*cashout.Service).Get),
		getCase("counterparty", counterparty.NewService, (*counterparty.Service).Get),
		getCase("country", country.NewService, (*country.Service).Get),
		getCase("demand", demand.NewService, (*demand.Service).Get),
		getCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).Get),
		getCase("invoicein", invoicein.NewService, (*invoicein.Service).Get),
		getCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).Get),
		getCase("move", move.NewService, (*move.Service).Get),
		getCase("organization", organization.NewService, (*organization.Service).Get),
		getCase("paymentin", paymentin.NewService, (*paymentin.Service).Get),
		getCase("paymentout", paymentout.NewService, (*paymentout.Service).Get),
		getCase("product", product.NewService, (*product.Service).Get),
		getCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).Get),
		getCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).Get),
		getCase("supply", supply.NewService, (*supply.Service).Get),
		getCase("uom", uom.NewService, (*uom.Service).Get),
		createCase("cashin", cashin.NewService, (*cashin.Service).Create),
		updateCase("cashin", cashin.NewService, (*cashin.Service).Update),
		createCase("cashout", cashout.NewService, (*cashout.Service).Create),
		updateCase("cashout", cashout.NewService, (*cashout.Service).Update),
		createCase("counterparty", counterparty.NewService, (*counterparty.Service).Create),
		updateCase("counterparty", counterparty.NewService, (*counterparty.Service).Update),
		createCase("country", country.NewService, (*country.Service).Create),
		updateCase("country", country.NewService, (*country.Service).Update),
		createCase("demand", demand.NewService, (*demand.Service).Create),
		updateCase("demand", demand.NewService, (*demand.Service).Update),
		createCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).Create),
		updateCase("expenseitem", expenseitem.NewService, (*expenseitem.Service).Update),
		createCase("invoicein", invoicein.NewService, (*invoicein.Service).Create),
		updateCase("invoicein", invoicein.NewService, (*invoicein.Service).Update),
		createCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).Create),
		updateCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).Update),
		createCase("move", move.NewService, (*move.Service).Create),
		updateCase("move", move.NewService, (*move.Service).Update),
		createCase("organization", organization.NewService, (*organization.Service).Create),
		updateCase("organization", organization.NewService, (*organization.Service).Update),
		createCase("paymentin", paymentin.NewService, (*paymentin.Service).Create),
		updateCase("paymentin", paymentin.NewService, (*paymentin.Service).Update),
		createCase("paymentout", paymentout.NewService, (*paymentout.Service).Create),
		updateCase("paymentout", paymentout.NewService, (*paymentout.Service).Update),
		createCase("product", product.NewService, (*product.Service).Create),
		updateCase("product", product.NewService, (*product.Service).Update),
		createCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).Create),
		updateCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).Update),
		createCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).Create),
		updateCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).Update),
		createCase("supply", supply.NewService, (*supply.Service).Create),
		updateCase("supply", supply.NewService, (*supply.Service).Update),
		createCase("uom", uom.NewService, (*uom.Service).Create),
		updateCase("uom", uom.NewService, (*uom.Service).Update),
		positionsCase("demand", demand.NewService, (*demand.Service).GetPositions),
		positionsCase("invoicein", invoicein.NewService, (*invoicein.Service).GetPositions),
		positionsCase("invoiceout", invoiceout.NewService, (*invoiceout.Service).GetPositions),
		positionsCase("move", move.NewService, (*move.Service).GetPositions),
		positionsCase("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).GetPositions),
		positionsCase("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).GetPositions),
		positionsCase("supply", supply.NewService, (*supply.Service).GetPositions),
	} {
		t.Run(tc.name, func(t *testing.T) {
			testServiceOptions(t, tc, false)
			testServiceOptions(t, tc, true)
		})
	}
}

// Метаданные ответа и APIError доходят до вызывающего кода и при ошибке авторизации.
func testServiceOptions(t *testing.T, tc optionsCase, unauthorized bool) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != tc.method || r.URL.RequestURI() != "/"+tc.path || r.Header.Get("Accept-Language") != "en" ||
			r.Header.Get("X-Lognex-Accept-Region") != "kz" {
			t.Errorf("request = %s %s %v", r.Method, r.URL.RequestURI(), r.Header)
		}
		w.Header().Set("X-Lognex-Content-Region", "kz")
		if unauthorized {
			w.Header().Set("X-Lognex-Auth", "-1")
			w.Header().Set("X-Lognex-Auth-Message", "token expired")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":[{"error":"auth","code":1056}]}`))
			return
		}
		if strings.HasSuffix(r.URL.Path, "/positions") {
			w.Write([]byte(`{"meta":{"href":"https://api.moysklad.ru/x"},"rows":[]}`))
			return
		}
		w.Write([]byte(`{"id":"` + testID + `"}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL,
		Header: http.Header{"Accept-Language": {"ru"}}})
	if err != nil {
		t.Fatal(err)
	}

	var info transport.ResponseInfo
	err = tc.invoke(t.Context(), client, transport.WithAcceptLanguage("en"),
		transport.WithAcceptRegion(transport.RegionKZ), transport.WithResponseInfo(&info))
	want := http.StatusOK
	if unauthorized {
		want = http.StatusUnauthorized
		var apiErr *sdkerrors.APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != want ||
			info.Header.Get(transport.HeaderAuthMessage) != "token expired" || info.Header.Get(transport.HeaderAuth) != "-1" {
			t.Fatalf("err = %v, info = %+v", err, info)
		}
	} else if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 || info.StatusCode != want || info.Header.Get("X-Lognex-Content-Region") != "kz" {
		t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
	}
}

// Регрессия ревью 46e9d06: локальная ошибка (валидация до запроса) не оставляет в ResponseInfo
// метаданных предыдущего успешного вызова.
func TestLocalErrorResetsResponseInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Lognex-Content-Region", "kz")
		w.Write([]byte(`{"id":"` + testID + `"}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	units, err := uom.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	demands, err := demand.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	for name, call := range map[string]func() error{
		"get invalid id": func() error {
			_, err := units.Get(t.Context(), "bad", transport.WithResponseInfo(&info))
			return err
		},
		"list limit": func() error {
			_, err := units.List(t.Context(), query.ListOptions{Limit: 1001}, transport.WithResponseInfo(&info))
			return err
		},
		"positions invalid id": func() error {
			_, err := demands.GetPositions(t.Context(), "bad", positions.PageOptions{}, transport.WithResponseInfo(&info))
			return err
		},
	} {
		if _, err := units.Get(t.Context(), testID, transport.WithResponseInfo(&info)); err != nil || info.StatusCode != 200 {
			t.Fatalf("%s: setup: err = %v, info = %+v", name, err, info)
		}
		if err := call(); err == nil || info.StatusCode != 0 || info.Attempts != 0 || info.Header != nil {
			t.Errorf("%s: err = %v, stale info = %+v", name, err, info)
		}
	}
}
