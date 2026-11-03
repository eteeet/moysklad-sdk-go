// Проверка T02: каждый service передаёт per-call опции в transport и заполняет ResponseInfo.
package entity_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/assortment"
	"github.com/eteeet/moysklad-sdk-go/entity/bonusprogram"
	"github.com/eteeet/moysklad-sdk-go/entity/bonustransaction"
	"github.com/eteeet/moysklad-sdk-go/entity/bundle"
	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/entity/consignment"
	"github.com/eteeet/moysklad-sdk-go/entity/contract"
	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/entity/currency"
	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/entity/group"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/organization"
	"github.com/eteeet/moysklad-sdk-go/entity/organizationbranch"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/entity/pricetype"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/entity/productfolder"
	"github.com/eteeet/moysklad-sdk-go/entity/project"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/region"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/entity/retailsalesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/retailshift"
	"github.com/eteeet/moysklad-sdk-go/entity/retailstore"
	"github.com/eteeet/moysklad-sdk-go/entity/role"
	"github.com/eteeet/moysklad-sdk-go/entity/saleplatform"
	"github.com/eteeet/moysklad-sdk-go/entity/saleschannel"
	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/service"
	"github.com/eteeet/moysklad-sdk-go/entity/store"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/entity/taxrate"
	"github.com/eteeet/moysklad-sdk-go/entity/thing"
	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/files"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
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

// fillRequest заполняет обязательные для локальной проверки поля: Name/LastName, ссылки, Description,
// Code/ISOCode, rate и непустые срезы (components, characteristics).
func fillRequest[I any]() I {
	var input I
	v := reflect.ValueOf(&input).Elem()
	for i := range v.NumField() {
		field := v.Field(i)
		switch {
		case (v.Type().Field(i).Name == "Name" || v.Type().Field(i).Name == "LastName" || v.Type().Field(i).Name == "Label") && field.Kind() == reflect.String:
			field.SetString("x")
		case v.Type().Field(i).Name == "TransactionType" && field.Kind() == reflect.String:
			field.SetString("EARNING")
		case v.Type().Field(i).Name == "Type" && field.Kind() == reflect.String:
			field.SetString("OTHER")
		case field.Type() == reflect.TypeFor[model.PriceType]():
			field.Set(reflect.ValueOf(model.PriceType{Meta: model.Meta{Href: "https://api.moysklad.ru/x"}}))
		case field.Type() == reflect.TypeFor[model.Reference]():
			field.Set(reflect.ValueOf(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/x"}}))
		case slices.Contains([]string{"Description", "Name", "Code", "ISOCode"}, v.Type().Field(i).Name) && field.Type() == reflect.TypeFor[model.Field[string]]():
			field.Set(reflect.ValueOf(model.Value("x")))
		case field.Type() == reflect.TypeFor[*float64]():
			field.Set(reflect.ValueOf(new(1.0)))
		case v.Type().Field(i).Name == "Rate" && field.Type() == reflect.TypeFor[model.Field[float64]]():
			field.Set(reflect.ValueOf(model.Value(1.0)))
		case field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct:
			field.Set(reflect.MakeSlice(field.Type(), 1, 1))
		}
	}
	return input
}

// catalogCase — GET вне формы entity/{name}/{id}; NewService падает только на nil client.
func catalogCase(name, path string, invoke func(context.Context, *transport.Client, ...transport.RequestOption) error) optionsCase {
	return optionsCase{name, http.MethodGet, path, invoke}
}

func TestServicesAcceptRequestOptions(t *testing.T) {
	for _, tc := range []optionsCase{
		getCase("cashin", cashin.NewService, (*cashin.Service).Get),
		getCase("cashout", cashout.NewService, (*cashout.Service).Get),
		getCase("counterparty", counterparty.NewService, (*counterparty.Service).Get),
		getCase("country", country.NewService, (*country.Service).Get),
		getCase("demand", demand.NewService, (*demand.Service).Get),
		getCase("employee", employee.NewService, (*employee.Service).Get),
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
		getCase("role", role.NewService, (*role.Service).Get),
		getCase("supply", supply.NewService, (*supply.Service).Get),
		getCase("uom", uom.NewService, (*uom.Service).Get),
		getCase("consignment", consignment.NewService, (*consignment.Service).Get),
		getCase("thing", thing.NewService, (*thing.Service).Get),
		getCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).Get),
		getCase("saleplatform", saleplatform.NewService, (*saleplatform.Service).Get),
		getCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).Get),
		getCase("region", region.NewService, (*region.Service).Get),
		getCase("contract", contract.NewService, (*contract.Service).Get),
		getCase("project", project.NewService, (*project.Service).Get),
		getCase("group", group.NewService, (*group.Service).Get),
		getCase("productfolder", productfolder.NewService, (*productfolder.Service).Get),
		getCase("store", store.NewService, (*store.Service).Get),
		getCase("currency", currency.NewService, (*currency.Service).Get),
		getCase("taxrate", taxrate.NewService, (*taxrate.Service).Get),
		getCase("service", service.NewService, (*service.Client).Get),
		getCase("bundle", bundle.NewService, (*bundle.Service).Get),
		getCase("variant", variant.NewService, (*variant.Service).Get),
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
		createCase("employee", employee.NewService, (*employee.Service).Create),
		updateCase("employee", employee.NewService, (*employee.Service).Update),
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
		createCase("role", role.NewService, (*role.Service).Create),
		updateCase("role", role.NewService, (*role.Service).Update),
		createCase("supply", supply.NewService, (*supply.Service).Create),
		updateCase("supply", supply.NewService, (*supply.Service).Update),
		createCase("uom", uom.NewService, (*uom.Service).Create),
		updateCase("uom", uom.NewService, (*uom.Service).Update),
		createCase("consignment", consignment.NewService, (*consignment.Service).Create),
		updateCase("consignment", consignment.NewService, (*consignment.Service).Update),
		createCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).Create),
		updateCase("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).Update),
		createCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).Create),
		updateCase("saleschannel", saleschannel.NewService, (*saleschannel.Service).Update),
		createCase("contract", contract.NewService, (*contract.Service).Create),
		updateCase("contract", contract.NewService, (*contract.Service).Update),
		createCase("project", project.NewService, (*project.Service).Create),
		updateCase("project", project.NewService, (*project.Service).Update),
		createCase("group", group.NewService, (*group.Service).Create),
		updateCase("group", group.NewService, (*group.Service).Update),
		createCase("productfolder", productfolder.NewService, (*productfolder.Service).Create),
		updateCase("productfolder", productfolder.NewService, (*productfolder.Service).Update),
		createCase("store", store.NewService, (*store.Service).Create),
		updateCase("store", store.NewService, (*store.Service).Update),
		createCase("currency", currency.NewService, (*currency.Service).Create),
		updateCase("currency", currency.NewService, (*currency.Service).Update),
		createCase("taxrate", taxrate.NewService, (*taxrate.Service).Create),
		updateCase("taxrate", taxrate.NewService, (*taxrate.Service).Update),
		createCase("service", service.NewService, (*service.Client).Create),
		updateCase("service", service.NewService, (*service.Client).Update),
		createCase("bundle", bundle.NewService, (*bundle.Service).Create),
		updateCase("bundle", bundle.NewService, (*bundle.Service).Update),
		createCase("variant", variant.NewService, (*variant.Service).Create),
		updateCase("variant", variant.NewService, (*variant.Service).Update),
		catalogCase("pricetype get", "context/companysettings/pricetype/"+testID, func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := pricetype.NewService(c)
			_, err := s.Get(ctx, testID, opts...)
			return err
		}),
		catalogCase("variant characteristic get", "entity/variant/metadata/characteristics/"+testID, func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := variant.NewService(c)
			_, err := s.GetCharacteristic(ctx, testID, opts...)
			return err
		}),
		catalogCase("bundle component get", "entity/bundle/"+testID+"/components/"+testID, func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := bundle.NewService(c)
			_, err := s.GetComponent(ctx, testID, testID, opts...)
			return err
		}),
		catalogCase("bundle components", "entity/bundle/"+testID+"/components?limit=1000&offset=0", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := bundle.NewService(c)
			_, err := s.GetComponents(ctx, testID, query.PageOptions{}, opts...)
			return err
		}),
		catalogCase("files list", "entity/product/"+testID+"/files?limit=1000&offset=0", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			_, err := s.List(ctx, "product", testID, query.PageOptions{}, opts...)
			return err
		}),
		catalogCase("turnover all", "report/turnover/all?limit=1000&offset=0", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := turnover.NewService(c)
			_, err := s.All(ctx, turnover.AllOptions{}, opts...)
			return err
		}),
		catalogCase("turnover bystore", "report/turnover/bystore?filter=product%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2Fx", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := turnover.NewService(c)
			_, err := s.ByStore(ctx, turnover.DetailOptions{Assortment: model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/x", Type: "product"}}}, opts...)
			return err
		}),
		catalogCase("turnover byoperations", "report/turnover/byoperations?filter=product%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2Fx", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := turnover.NewService(c)
			_, err := s.ByOperations(ctx, turnover.DetailOptions{Assortment: model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/x", Type: "product"}}}, opts...)
			return err
		}),
		catalogCase("assortment list", "entity/assortment", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := assortment.NewService(c)
			_, err := s.List(ctx, assortment.Options{}, opts...)
			return err
		}),
		catalogCase("assortment legacy", "entity/assortment", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := assortment.NewService(c)
			_, err := s.ListLegacy(ctx, assortment.Options{}, opts...)
			return err
		}),
		getCase("customerorder", customerorder.NewService, (*customerorder.Service).Get),
		createCase("customerorder", customerorder.NewService, (*customerorder.Service).Create),
		updateCase("customerorder", customerorder.NewService, (*customerorder.Service).Update),
		positionsCase("customerorder", customerorder.NewService, (*customerorder.Service).GetPositions),
		getCase("enter", enter.NewService, (*enter.Service).Get),
		createCase("enter", enter.NewService, (*enter.Service).Create),
		updateCase("enter", enter.NewService, (*enter.Service).Update),
		positionsCase("enter", enter.NewService, (*enter.Service).GetPositions),
		getCase("loss", loss.NewService, (*loss.Service).Get),
		createCase("loss", loss.NewService, (*loss.Service).Create),
		updateCase("loss", loss.NewService, (*loss.Service).Update),
		positionsCase("loss", loss.NewService, (*loss.Service).GetPositions),
		getCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).Get),
		createCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).Create),
		updateCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).Update),
		positionsCase("salesreturn", salesreturn.NewService, (*salesreturn.Service).GetPositions),
		getCase("retailstore", retailstore.NewService, (*retailstore.Service).Get),
		createCase("retailstore", retailstore.NewService, (*retailstore.Service).Create),
		updateCase("retailstore", retailstore.NewService, (*retailstore.Service).Update),
		getCase("retailshift", retailshift.NewService, (*retailshift.Service).Get),
		createCase("retailshift", retailshift.NewService, (*retailshift.Service).Create),
		updateCase("retailshift", retailshift.NewService, (*retailshift.Service).Update),
		getCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).Get),
		createCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).Create),
		updateCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).Update),
		getCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).Get),
		createCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).Create),
		updateCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).Update),
		getCase("bonusprogram", bonusprogram.NewService, (*bonusprogram.Service).Get),
		createCase("bonusprogram", bonusprogram.NewService, (*bonusprogram.Service).Create),
		updateCase("bonusprogram", bonusprogram.NewService, (*bonusprogram.Service).Update),
		getCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).Get),
		createCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).Create),
		updateCase("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).Update),
		positionsCase("retaildemand", retaildemand.NewService, (*retaildemand.Service).GetPositions),
		positionsCase("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).GetPositions),
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
		if strings.HasSuffix(r.URL.Path, "/positions") || strings.HasSuffix(r.URL.Path, "/components") || strings.HasSuffix(r.URL.Path, "/assortment") ||
			strings.HasSuffix(r.URL.Path, "/files") || strings.HasPrefix(r.URL.Path, "/report/") {
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

// localCase — вызов, который падает на локальной проверке до запроса.
type localCase struct {
	name   string
	invoke func(context.Context, *transport.Client, ...transport.RequestOption) error
}

func localCases[S, C, U, T any](name string, newService func(*transport.Client) (*S, error),
	get func(*S, context.Context, string, ...transport.RequestOption) (*T, error),
	create func(*S, context.Context, C, ...transport.RequestOption) (*T, error),
	update func(*S, context.Context, string, U, ...transport.RequestOption) (*T, error),
	list func(*S, context.Context, query.ListOptions, ...transport.RequestOption) (*model.Collection[T], error)) []localCase {
	call := func(f func(*S, context.Context, ...transport.RequestOption) error) func(context.Context, *transport.Client, ...transport.RequestOption) error {
		return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := newService(c)
			if err != nil {
				return err
			}
			return f(s, ctx, opts...)
		}
	}
	var zeroCreate C
	var zeroUpdate U
	return []localCase{
		{name + " get invalid id", call(func(s *S, ctx context.Context, opts ...transport.RequestOption) error {
			_, err := get(s, ctx, "bad", opts...)
			return err
		})},
		{name + " create empty", call(func(s *S, ctx context.Context, opts ...transport.RequestOption) error {
			_, err := create(s, ctx, zeroCreate, opts...)
			return err
		})},
		{name + " update no fields", call(func(s *S, ctx context.Context, opts ...transport.RequestOption) error {
			_, err := update(s, ctx, testID, zeroUpdate, opts...)
			return err
		})},
		{name + " list limit", call(func(s *S, ctx context.Context, opts ...transport.RequestOption) error {
			_, err := list(s, ctx, query.ListOptions{Limit: 1001}, opts...)
			return err
		})},
	}
}

func positionsLocal[S, P any](name string, newService func(*transport.Client) (*S, error),
	get func(*S, context.Context, string, positions.PageOptions, ...transport.RequestOption) (*model.Collection[P], error)) localCase {
	return localCase{name + " positions invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
		s, err := newService(c)
		if err != nil {
			return err
		}
		_, err = get(s, ctx, "bad", positions.PageOptions{}, opts...)
		return err
	}}
}

// composed — свежая составная опция поверх WithResponseInfo (сценарий ревью 35078b4).
func composed(info *transport.ResponseInfo) transport.RequestOption {
	return func(o *transport.RequestOptions) {
		transport.WithResponseInfo(info)(o)
	}
}

// Регрессия ревью 46e9d06 и 35078b4: локальная ошибка (валидация до запроса) не оставляет
// в ResponseInfo метаданных предыдущего вызова, и при составной опции тоже.
func TestLocalErrorResetsResponseInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Lognex-Content-Region", "kz")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"errors":[{"error":"x","code":1074}]}`))
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	var cases []localCase
	cases = append(cases, localCases("cashin", cashin.NewService, (*cashin.Service).Get, (*cashin.Service).Create, (*cashin.Service).Update, (*cashin.Service).List)...)
	cases = append(cases, localCases("cashout", cashout.NewService, (*cashout.Service).Get, (*cashout.Service).Create, (*cashout.Service).Update, (*cashout.Service).List)...)
	cases = append(cases, localCases("counterparty", counterparty.NewService, (*counterparty.Service).Get, (*counterparty.Service).Create, (*counterparty.Service).Update, (*counterparty.Service).List)...)
	cases = append(cases, localCases("country", country.NewService, (*country.Service).Get, (*country.Service).Create, (*country.Service).Update, (*country.Service).List)...)
	cases = append(cases, localCases("demand", demand.NewService, (*demand.Service).Get, (*demand.Service).Create, (*demand.Service).Update, (*demand.Service).List)...)
	cases = append(cases, localCases("employee", employee.NewService, (*employee.Service).Get, (*employee.Service).Create, (*employee.Service).Update, (*employee.Service).List)...)
	cases = append(cases, localCases("expenseitem", expenseitem.NewService, (*expenseitem.Service).Get, (*expenseitem.Service).Create, (*expenseitem.Service).Update, (*expenseitem.Service).List)...)
	cases = append(cases, localCases("invoicein", invoicein.NewService, (*invoicein.Service).Get, (*invoicein.Service).Create, (*invoicein.Service).Update, (*invoicein.Service).List)...)
	cases = append(cases, localCases("invoiceout", invoiceout.NewService, (*invoiceout.Service).Get, (*invoiceout.Service).Create, (*invoiceout.Service).Update, (*invoiceout.Service).List)...)
	cases = append(cases, localCases("move", move.NewService, (*move.Service).Get, (*move.Service).Create, (*move.Service).Update, (*move.Service).List)...)
	cases = append(cases, localCases("organization", organization.NewService, (*organization.Service).Get, (*organization.Service).Create, (*organization.Service).Update, (*organization.Service).List)...)
	cases = append(cases, localCases("paymentin", paymentin.NewService, (*paymentin.Service).Get, (*paymentin.Service).Create, (*paymentin.Service).Update, (*paymentin.Service).List)...)
	cases = append(cases, localCases("paymentout", paymentout.NewService, (*paymentout.Service).Get, (*paymentout.Service).Create, (*paymentout.Service).Update, (*paymentout.Service).List)...)
	cases = append(cases, localCases("product", product.NewService, (*product.Service).Get, (*product.Service).Create, (*product.Service).Update, (*product.Service).List)...)
	cases = append(cases, localCases("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).Get, (*purchaseorder.Service).Create, (*purchaseorder.Service).Update, (*purchaseorder.Service).List)...)
	cases = append(cases, localCases("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).Get, (*purchasereturn.Service).Create, (*purchasereturn.Service).Update, (*purchasereturn.Service).List)...)
	cases = append(cases, localCases("role", role.NewService, (*role.Service).Get, (*role.Service).Create, (*role.Service).Update, (*role.Service).List)...)
	cases = append(cases, localCases("supply", supply.NewService, (*supply.Service).Get, (*supply.Service).Create, (*supply.Service).Update, (*supply.Service).List)...)
	cases = append(cases, localCases("uom", uom.NewService, (*uom.Service).Get, (*uom.Service).Create, (*uom.Service).Update, (*uom.Service).List)...)
	cases = append(cases, localCases("consignment", consignment.NewService, (*consignment.Service).Get, (*consignment.Service).Create, (*consignment.Service).Update, (*consignment.Service).List)...)
	cases = append(cases, localCases("organizationbranch", organizationbranch.NewService, (*organizationbranch.Service).Get, (*organizationbranch.Service).Create, (*organizationbranch.Service).Update, (*organizationbranch.Service).List)...)
	cases = append(cases, localCases("saleschannel", saleschannel.NewService, (*saleschannel.Service).Get, (*saleschannel.Service).Create, (*saleschannel.Service).Update, (*saleschannel.Service).List)...)
	cases = append(cases, localCases("contract", contract.NewService, (*contract.Service).Get, (*contract.Service).Create, (*contract.Service).Update, (*contract.Service).List)...)
	cases = append(cases, localCases("project", project.NewService, (*project.Service).Get, (*project.Service).Create, (*project.Service).Update, (*project.Service).List)...)
	cases = append(cases, localCases("group", group.NewService, (*group.Service).Get, (*group.Service).Create, (*group.Service).Update, (*group.Service).List)...)
	cases = append(cases, localCases("productfolder", productfolder.NewService, (*productfolder.Service).Get, (*productfolder.Service).Create, (*productfolder.Service).Update, (*productfolder.Service).List)...)
	cases = append(cases, localCases("store", store.NewService, (*store.Service).Get, (*store.Service).Create, (*store.Service).Update, (*store.Service).List)...)
	cases = append(cases, localCases("currency", currency.NewService, (*currency.Service).Get, (*currency.Service).Create, (*currency.Service).Update, (*currency.Service).List)...)
	cases = append(cases, localCases("taxrate", taxrate.NewService, (*taxrate.Service).Get, (*taxrate.Service).Create, (*taxrate.Service).Update, (*taxrate.Service).List)...)
	cases = append(cases, localCases("service", service.NewService, (*service.Client).Get, (*service.Client).Create, (*service.Client).Update, (*service.Client).List)...)
	cases = append(cases, localCases("bundle", bundle.NewService, (*bundle.Service).Get, (*bundle.Service).Create, (*bundle.Service).Update, (*bundle.Service).List)...)
	cases = append(cases, localCases("variant", variant.NewService, (*variant.Service).Get, (*variant.Service).Create, (*variant.Service).Update, (*variant.Service).List)...)
	cases = append(cases, localCases("customerorder", customerorder.NewService, (*customerorder.Service).Get, (*customerorder.Service).Create, (*customerorder.Service).Update, (*customerorder.Service).List)...)
	cases = append(cases, localCases("enter", enter.NewService, (*enter.Service).Get, (*enter.Service).Create, (*enter.Service).Update, (*enter.Service).List)...)
	cases = append(cases, localCases("loss", loss.NewService, (*loss.Service).Get, (*loss.Service).Create, (*loss.Service).Update, (*loss.Service).List)...)
	cases = append(cases, localCases("salesreturn", salesreturn.NewService, (*salesreturn.Service).Get, (*salesreturn.Service).Create, (*salesreturn.Service).Update, (*salesreturn.Service).List)...)
	cases = append(cases, localCases("retailstore", retailstore.NewService, (*retailstore.Service).Get, (*retailstore.Service).Create, (*retailstore.Service).Update, (*retailstore.Service).List)...)
	cases = append(cases, localCases("retailshift", retailshift.NewService, (*retailshift.Service).Get, (*retailshift.Service).Create, (*retailshift.Service).Update, (*retailshift.Service).List)...)
	cases = append(cases, localCases("retaildemand", retaildemand.NewService, (*retaildemand.Service).Get, (*retaildemand.Service).Create, (*retaildemand.Service).Update, (*retaildemand.Service).List)...)
	cases = append(cases, localCases("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).Get, (*retailsalesreturn.Service).Create, (*retailsalesreturn.Service).Update, (*retailsalesreturn.Service).List)...)
	cases = append(cases, localCases("bonusprogram", bonusprogram.NewService, (*bonusprogram.Service).Get, (*bonusprogram.Service).Create, (*bonusprogram.Service).Update, (*bonusprogram.Service).List)...)
	cases = append(cases, localCases("bonustransaction", bonustransaction.NewService, (*bonustransaction.Service).Get, (*bonustransaction.Service).Create, (*bonustransaction.Service).Update, (*bonustransaction.Service).List)...)
	cases = append(cases,
		positionsLocal("retaildemand", retaildemand.NewService, (*retaildemand.Service).GetPositions),
		positionsLocal("retailsalesreturn", retailsalesreturn.NewService, (*retailsalesreturn.Service).GetPositions),
		positionsLocal("customerorder", customerorder.NewService, (*customerorder.Service).GetPositions),
		positionsLocal("enter", enter.NewService, (*enter.Service).GetPositions),
		positionsLocal("loss", loss.NewService, (*loss.Service).GetPositions),
		positionsLocal("salesreturn", salesreturn.NewService, (*salesreturn.Service).GetPositions),
		positionsLocal("demand", demand.NewService, (*demand.Service).GetPositions),
		positionsLocal("invoicein", invoicein.NewService, (*invoicein.Service).GetPositions),
		positionsLocal("invoiceout", invoiceout.NewService, (*invoiceout.Service).GetPositions),
		positionsLocal("move", move.NewService, (*move.Service).GetPositions),
		positionsLocal("purchaseorder", purchaseorder.NewService, (*purchaseorder.Service).GetPositions),
		positionsLocal("purchasereturn", purchasereturn.NewService, (*purchasereturn.Service).GetPositions),
		positionsLocal("supply", supply.NewService, (*supply.Service).GetPositions),
		localCase{"employee get security invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := employee.NewService(c)
			if err != nil {
				return err
			}
			_, err = s.GetSecurity(ctx, "bad", opts...)
			return err
		}},
		localCase{"employee update security empty", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := employee.NewService(c)
			if err != nil {
				return err
			}
			_, err = s.UpdateSecurity(ctx, testID, employee.UpdateSecurityRequest{}, opts...)
			return err
		}},
		localCase{"pricetype get invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := pricetype.NewService(c)
			_, err := s.Get(ctx, "bad", opts...)
			return err
		}},
		localCase{"pricetype replace empty", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := pricetype.NewService(c)
			_, err := s.ReplaceAll(ctx, nil, opts...)
			return err
		}},
		localCase{"variant characteristic invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := variant.NewService(c)
			_, err := s.GetCharacteristic(ctx, "bad", opts...)
			return err
		}},
		localCase{"variant characteristic empty name", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := variant.NewService(c)
			_, err := s.CreateCharacteristic(ctx, "", opts...)
			return err
		}},
		localCase{"bundle components invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := bundle.NewService(c)
			_, err := s.GetComponents(ctx, "bad", query.PageOptions{}, opts...)
			return err
		}},
		localCase{"bundle component invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := bundle.NewService(c)
			_, err := s.GetComponent(ctx, testID, "bad", opts...)
			return err
		}},
		localCase{"assortment list fields", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := assortment.NewService(c)
			_, err := s.List(ctx, assortment.Options{ListOptions: query.ListOptions{Fields: "stock"}}, opts...)
			return err
		}},
		localCase{"assortment legacy limit", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := assortment.NewService(c)
			_, err := s.ListLegacy(ctx, assortment.Options{ListOptions: query.ListOptions{Limit: 1001}}, opts...)
			return err
		}},
		localCase{"files list invalid id", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			_, err := s.List(ctx, "product", "bad", query.PageOptions{}, opts...)
			return err
		}},
		localCase{"files upload empty", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			_, err := s.Upload(ctx, "product", testID, nil, opts...)
			return err
		}},
		localCase{"files download no href", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := files.NewService(c)
			_, err := s.Download(ctx, files.File{}, io.Discard, opts...)
			return err
		}},
		localCase{"turnover all limit", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := turnover.NewService(c)
			_, err := s.All(ctx, turnover.AllOptions{Limit: -1}, opts...)
			return err
		}},
		localCase{"turnover bystore no assortment", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := turnover.NewService(c)
			_, err := s.ByStore(ctx, turnover.DetailOptions{}, opts...)
			return err
		}},
		localCase{"turnover byoperations no assortment", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, _ := turnover.NewService(c)
			_, err := s.ByOperations(ctx, turnover.DetailOptions{}, opts...)
			return err
		}},
		localCase{"transport download nil writer", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			_, err := c.Download(ctx, "download/1", nil, opts...)
			return err
		}},
		localCase{"role preset unsafe", func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) error {
			s, err := role.NewService(c)
			if err != nil {
				return err
			}
			_, err = s.Preset(ctx, "../x", opts...)
			return err
		}},
	)
	for _, tc := range cases {
		for _, option := range map[string]func(*transport.ResponseInfo) transport.RequestOption{
			"direct": transport.WithResponseInfo, "composed": composed,
		} {
			// Прошлый вызов оставил 503 и три попытки.
			info := transport.ResponseInfo{StatusCode: 503, Attempts: 3, Header: http.Header{"X-Stale": {"1"}}}
			err := tc.invoke(t.Context(), client, option(&info))
			if err == nil || sdkerrors.StatusCode(err) != 0 {
				t.Fatalf("%s: expected local error, got %v", tc.name, err)
			}
			if info.StatusCode != 0 || info.Attempts != 0 || info.Header != nil {
				t.Errorf("%s: stale info = %+v", tc.name, info)
			}
		}
	}
	if len(cases) != 42*4+13+3+8+1+6 {
		t.Fatalf("cases = %d", len(cases))
	}
}
